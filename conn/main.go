package main

import (
	"go.temporal.io/sdk/client"
)

func main() {
	_, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		panic(err)
	}
}
