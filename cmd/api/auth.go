package main

import (
	"encoding/json"
	"time"

	"net/http"

	"github.com/devaartana/e01-oprec-rpl/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if length := len(payload.Password); length < 8 {
		http.Error(w, "Password is not long enough", http.StatusBadRequest)
		return
	}

	user := &store.User{
		Username:   payload.Username,
		Email: 	payload.Email,
		Created_at: time.Now(),
		Links: 	[]store.Link{},
	}

	user.SetPassword(payload.Password)

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User is registered"))
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if length := len(payload.Password); length < 8 {
		http.Error(w, "Password is not long enough", http.StatusBadRequest)
		return
	}

	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		http.Error(w, "Email is not exist", http.StatusInternalServerError)
		return
	}

	if user.Compare(payload.Password) != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"email": payload.Email,
		"exp":   time.Now().Add(app.config.auth.exp).Unix(),
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
		"iss":   app.config.auth.iss,
		"aud":   app.config.auth.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}


func (app *application) UserHandler(w http.ResponseWriter, r *http.Request) {
	user, err:= app.store.Users.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}