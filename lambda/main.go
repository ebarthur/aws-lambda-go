package main

import (
	"fmt"
	"lambda-func/app"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
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
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myapp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myapp.ApiHandler.LoginUserHandler(request)
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       "Not Found",
			}, nil
		}
	})
}
