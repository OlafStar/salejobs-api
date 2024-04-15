package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/olafstar/salejobs-api/internal/passwords"
	"github.com/olafstar/salejobs-api/internal/types"
)

var secretKey = []byte("secret-key")

type User = types.User
type DecodedToken = types.DecodedToken
type LoginResponse struct {
	Token string `json:"token"`
}

func (s *APIServer) JWTLoginHandler(w http.ResponseWriter, r *http.Request) error {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			return &HTTPError{StatusCode: http.StatusBadRequest, Message: "error decoding user data"}
	}

	storedHash, err := s.store.GetJWTUser(u.Username)
	if err != nil {
			return &HTTPError{StatusCode: http.StatusUnauthorized, Message: "Authorization failed"}
	}

	if !passwords.CheckPasswordHash(u.Password, storedHash) {
			return &HTTPError{StatusCode: http.StatusUnauthorized, Message: "invalid credentials"}
	}

	tokenString, err := createToken(u.Username)
	if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "error creating token"}
	}

	return WriteJSON(w, http.StatusOK, LoginResponse{
		Token: tokenString,
	})
}

func ProtectedRequest(callback apiFunc) apiFunc{
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			return &HTTPError{StatusCode: http.StatusUnauthorized, Message: "Missing auth header"}
		}

		tokenString = tokenString[len("Bearer "):]

		err := verifyToken(tokenString)

		if err != nil {
			return &HTTPError{StatusCode: http.StatusUnauthorized, Message: "Missing auth header"}
		}

		return callback(w, r)
	}
}

func DecodeToken(tokenString string) (DecodedToken, error) {
	token, err := jwt.ParseWithClaims(tokenString[len("Bearer "):], &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
	})

	if err != nil {
			return DecodedToken{}, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
			return DecodedToken{}, fmt.Errorf("invalid token")
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
			return DecodedToken{}, fmt.Errorf("unable to extract user from token")
	}

	exp, ok := (*claims)["exp"].(float64) 
	if !ok {
			return DecodedToken{}, fmt.Errorf("unable to extract expiration from token")
	}

	return DecodedToken{
			Username:   username,
			Exp: int64(exp),
	}, nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func verifyToken(tokenString string) error{
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
 })

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// func refreshToken(decodedToken DecodedToken) (string, error) {
// 	expirationTime := time.Now().Add(30 * time.Minute)
// 	claims := &jwt.MapClaims{
// 		"username": decodedToken.Username,
// 		"exp":      expirationTime.Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(secretKey)
// }
