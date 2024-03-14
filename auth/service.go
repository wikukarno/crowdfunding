package auth

import (
	"errors"
	"os"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
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

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
		SECRET_KEY := os.Getenv("SECRET_KEY")
		return []byte(SECRET_KEY), nil
	
	})

	if err != nil {
		return token, err
	}

	return token, nil
}