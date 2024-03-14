package auth

import (
	"errors"
	"os"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Service interface {
	GenerateToken(userID int) (string, error)

}

type jwtService struct {

}


func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {

	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")
	if SECRET_KEY == ""{
		return "", errors.New("SECRET_KEY is not defined")
	}

	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}