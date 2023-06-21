package api

import (
	"encoding/json"
	"net/http"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Name        string `json:"name" validate:"omitempty"`
	Email       string `json:"email" validate:"omitempty,email"`
	Pass        string `json:"password" validate:"omitempty"`
	LcpHintMsg  string `json:"lcpHintMsg"`
	LcpPassHash string `json:"lcpPassHash"`
	SessionId   string `json:"sessionId" validate:"-"`
}

func ConvertUserFromUserStor(u stor.User) *User {
	return &User{
		Name:        u.Name,
		Email:       u.Email,
		Pass:        u.Pass,
		LcpHintMsg:  u.LcpHintMsg,
		LcpPassHash: u.LcpPassHash,
		SessionId:   u.SessionId,
	}
}

func ConvertUserToUserStor(u User) *stor.User {
	return &stor.User{
		Name:        u.Name,
		Email:       u.Email,
		Pass:        u.Pass,
		LcpHintMsg:  u.LcpHintMsg,
		LcpPassHash: u.LcpPassHash,
		SessionId:   u.SessionId,
	}
}

func (api *Api) getUserHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")

	// Call your storage function to get the user by session ID
	storUser, err := api.stor.GetUserBySessionId(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the storage user to the view model user
	viewUser := ConvertUserFromUserStor(*storUser)

	// Set the content type header and write the response
	w.Header().Set("Content-Type", "application/json")

	// Encode the user as JSON and write it to the response
	err = json.NewEncoder(w).Encode(viewUser)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

/*
*

	curl -X POST -H "Content-Type: application/json" -d '{
	  "name": "John Doe",
	  "email": "johndoe@example.com",
	  "password": "password123",
	  "lcpHintMsg": "Hint",
	  "lcpPassHash": "hash123"
	}' http://localhost:8080/api/v1/user
*/
func (api *Api) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var viewUser User
	// Parse the JSON request body into the view model user
	err := json.NewDecoder(r.Body).Decode(&viewUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the publication struct using the validator
	err = validate.Struct(viewUser)
	if err != nil {
		// If validation fails, return the validation errors
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	// Convert the view model user to the storage user
	storUser := ConvertUserToUserStor(viewUser)

	// Call your storage function to create the user
	err = api.stor.CreateUser(storUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return success response
	w.WriteHeader(http.StatusCreated)
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the publication as JSON and write it to the response
	err = json.NewEncoder(w).Encode(viewUser)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
