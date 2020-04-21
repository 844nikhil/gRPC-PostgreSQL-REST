package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"gRPC-PostgreSQL-REST/server/interfaces"
	"gRPC-PostgreSQL-REST/server/model"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func Encrypt(pwd string) string {
	h := sha1.New()
	h.Write([]byte(pwd))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

//Create token with user details and expiry token
func CreateToken(user model.User) (string, error) {

	secretKey, err := getSecretKey()
	if err != nil {
		log.Println(err)
		return "", errors.New("Couldn't find secret key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.UserID,
		"email":    user.Email,
		"password": user.Password,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}

//Retreive secret key from env
func getSecretKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	return os.Getenv("TOKENSECRET"), nil
}

func ValidateToken(token string) (map[string]interface{}, error) {
	//decrypt the token and get the jwt claims
	mapClaims, tokenErr := DecryptToken(token)
	if tokenErr != nil {
		return nil, tokenErr
	}

	//Authenticate the token claims
	userExists, _ := interfaces.DBEngine.Authenticate(fmt.Sprintf("%v", mapClaims["email"]), fmt.Sprintf("%v", mapClaims["password"]))
	if userExists != nil {
		return nil, userExists
	}
	return mapClaims, nil
}

func DecryptToken(tokenString string) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, errors.New("Token not provided")
	}
	secretKey, err := getSecretKey()
	if err != nil {
		log.Println(err)
		return nil, errors.New("Couldn't find secret key")
	}
	token, _ := jwt.Parse(string(tokenString), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("There was an error")
		}
		return []byte(secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		mapClaims := token.Claims.(jwt.MapClaims)
		return mapClaims, nil
	}
	return nil, errors.New("Token error")
}
