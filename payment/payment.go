package payment

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type BankingService interface {
	Withdraw(accountNumber string, amount int, referenceID string) (string, error)
	Deposit(accountNumber string, amount int, referenceID string) (string, error)
}

type Activities struct {
	bank BankingService
}

type PaymentDetails struct {
	ReferenceID   string
	SourceAccount string
	TargetAccount string
	Amount        int
}

func MoneyTransferWorkflow(ctx workflow.Context, input PaymentDetails) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        100 * time.Second,
			NonRetryableErrorTypes: []string{"ErrInvalidAccount", "ErrInsufficientFunds"},
		},
	})

	var withdrawOutput string
	err := workflow.ExecuteActivity(ctx, "Withdraw", input).Get(ctx, &withdrawOutput)
	if err != nil {
		return "", err
	}

	var depositOutput string
	depositErr := workflow.ExecuteActivity(ctx, "Deposit", input).Get(ctx, &depositOutput)
	if depositErr != nil {
		var result string
		refundErr := workflow.ExecuteActivity(ctx, "Refund", input).Get(ctx, &result)

		if refundErr != nil {
			return "",
				fmt.Errorf("deposit: failed to deposit money into %v: %v. Money could not be returned to %v: %w",
					input.TargetAccount, depositErr, input.SourceAccount, refundErr)
		}

		return "", fmt.Errorf("deposit: failed to deposit money into %v: Money returned to %v: %w",
			input.TargetAccount, input.SourceAccount, depositErr)
	}

	return fmt.Sprintf("Transfer complete (transaction IDs: %s, %s)", withdrawOutput, depositOutput), nil
}

func (a *Activities) Withdraw(ctx context.Context, data PaymentDetails) (string, error) {
	referenceID := fmt.Sprintf("%s-withdrawal", data.ReferenceID)
	confirmation, err := a.bank.Withdraw(data.SourceAccount, data.Amount, referenceID)
	return confirmation, err
}

func (a *Activities) Deposit(ctx context.Context, data PaymentDetails) (string, error) {
	referenceID := fmt.Sprintf("%s-deposit", data.ReferenceID)
	confirmation, err := a.bank.Deposit(data.TargetAccount, data.Amount, referenceID)
	return confirmation, err
}

func (a *Activities) Refund(ctx context.Context, data PaymentDetails) (string, error) {
	referenceID := fmt.Sprintf("%s-refund", data.ReferenceID)
	confirmation, err := a.bank.Deposit(data.SourceAccount, data.Amount, referenceID)
	return confirmation, err
}
