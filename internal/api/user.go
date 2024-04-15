package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/olafstar/salejobs-api/internal/passwords"
)

func (s *APIServer) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.JWTLoginHandler(w, r)
	}
	
	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

func (s *APIServer) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.registerUser(w, r)
	}

	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

type RegisterBody struct {
	Username string
	Password string
}

// func validateInput(username, password string) bool {
// 	return true
// }

func (s *APIServer) registerUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var body RegisterBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Invalid request body"}
	}
	defer r.Body.Close()

	// Uncomment and modify the validation part if needed
	// if !validateInput(body.Username, body.Password) {
	//    return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Username or password does not meet the criteria"}
	// }

	fmt.Printf("The user request value %v\n", body)

	hashPass, hashErr := passwords.HashPassword(body.Password)
	if hashErr != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to hash password"}
	}

	var exists int
	if err := s.store.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", body.Username).Scan(&exists); err != nil || exists > 0 {
			if exists > 0 {
					return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Username already exists"}
			}
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Error checking username existence"}
	}

	_, err := s.store.db.Exec(
			"INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)",
			body.Username, hashPass, time.Now(),
	)
	if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Internal server error"}
	}

	fmt.Fprint(w, "User registered successfully")
	return nil
}

func (s *APIServer) Iam( w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		tokenString := r.Header.Get("Authorization")
		user, err := DecodeToken(tokenString) 

		if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Internal server error"}
		}

		return WriteJSON(w, http.StatusOK, user)
	}

	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}
