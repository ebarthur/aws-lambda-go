package main

import (
	"fmt"
	"lambda-func/app"

	"github.com/aws/aws-lambda-go/lambda"
)


type MyEvent struct {
	Username string `json:"username"`
}

// Take in a payload and do something with it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("Successfully called by - %s", event.Username), nil
}

func main() {
	myapp := app.NewApp()
	lambda.Start(myapp.ApiHandler.RegisterUserHandler)
}