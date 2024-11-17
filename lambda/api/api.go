package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) *ApiHandler {
	return &ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return types.ErrorResponse(http.StatusBadRequest, "Invalid request payload"), err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return types.ErrorResponse(http.StatusBadRequest, "Invalid request - fields empty"), err
	}

	userExists, err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return types.ErrorResponse(http.StatusInternalServerError, "Internal server error"), err

	}

	if userExists {
		return types.ErrorResponse(http.StatusConflict, "user already exists"), err
	}

	user, err := types.NewUser(registerUser)
	if err != nil {
		return types.ErrorResponse(http.StatusInternalServerError, "Internal server error"), fmt.Errorf("could not create user: %w", err)
	}

	err = api.dbStore.InsertUser(user)
	if err != nil {
		return types.ErrorResponse(http.StatusInternalServerError, "Internal server error"), err
	}

	return types.ErrorResponse(http.StatusCreated, "Successfully registered user"), nil
}

func (api ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return types.ErrorResponse(http.StatusBadRequest, "Invalid request payload"), err
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		return types.ErrorResponse(http.StatusBadRequest, "Invalid request - fields empty"), err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return types.ErrorResponse(http.StatusNotFound, "user not found"), err
	}

	if !types.ValidatePassword(user.HashPassword, loginRequest.Password) {
		return types.ErrorResponse(http.StatusBadRequest, "Invalid user credentials"), err
	}

	accessToken := types.CreateToken(user)

	messageBytes, err := json.Marshal(map[string]string{
		"token":   accessToken,
		"message": "Successfully logged in",
	})
	if err != nil {
		return types.ErrorResponse(http.StatusInternalServerError, "Internal server error"), err
	}
	successMessage := string(messageBytes)

	return types.ErrorResponse(http.StatusOK, successMessage), nil
}
