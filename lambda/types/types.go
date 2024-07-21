package types

import "golang.org/x/crypto/bcrypt"

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	HashPassword string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:     registerUser.Username,
		HashPassword: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashPassword, plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainTextPassword))
	return err == nil
}
