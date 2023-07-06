package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	UUID        string `json:"uuid" validate:"omitempty,uuid4_rfc4122"`
	Name        string `json:"name" validate:"omitempty"`
	Email       string `json:"email" validate:"omitempty,email"`
	Pass        string `json:"password" validate:"omitempty"`
	LcpHintMsg  string `json:"lcpHintMsg"`
	LcpPassHash string `json:"lcpPassHash"`
	SessionId   string `json:"-" validate:"-"`
}

func ConvertUserFromUserStor(u stor.User) *User {
	return &User{
		UUID:        u.UUID,
		Name:        u.Name,
		Email:       u.Email,
		Pass:        u.Pass,
		LcpHintMsg:  u.LcpHintMsg,
		LcpPassHash: u.LcpPassHash,
		SessionId:   u.SessionId,
	}
}

func ConvertUserToUserStor(u User, originalUser *stor.User) *stor.User {
	if originalUser == nil {
		originalUser = &stor.User{}
	}

	if u.UUID != "" {
		originalUser.UUID = u.UUID
	}
	if u.Name != "" {
		originalUser.Name = u.Name
	}
	if u.Email != "" {
		originalUser.Email = u.Email
	}
	if u.Pass != "" {
		originalUser.Pass = u.Pass
	}
	if u.LcpHintMsg != "" {
		originalUser.LcpHintMsg = u.LcpHintMsg
	}
	if u.LcpPassHash != "" {
		originalUser.LcpPassHash = u.LcpPassHash
	}
	if u.SessionId != "" {
		originalUser.SessionId = u.SessionId
	}

	return originalUser
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

// @Summary Create a new user
// @Description Create a new user with the provided payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object"
// @Success 201 {object} User "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or validation errors"
// @Failure 500 {object} ErrorResponse "Failed to create user"
// @Router /user [post]
func (api *Api) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var viewUser User
	// Parse the JSON request body into the view model user
	err := json.NewDecoder(r.Body).Decode(&viewUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	viewUser.SessionId = ""
	viewUser.LcpPassHash = lcp.CreateLcpPassHash(viewUser.LcpPassHash)

	// Generate UUID for the user
	if len(viewUser.UUID) == 0 {
		viewUser.UUID = uuid.New().String()
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
	storUser := ConvertUserToUserStor(viewUser, &stor.User{})

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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(viewUser)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// @Summary Get a user by ID
// @Description Retrieve a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User "OK"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/{id} [get]
func (api *Api) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storUser, ok := ctx.Value("user").(*stor.User)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Convert the storage user to the view model user
	viewUser := ConvertUserFromUserStor(*storUser)

	// Set the content type header and write the response
	w.Header().Set("Content-Type", "application/json")

	// Encode the user as JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(viewUser)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// @Summary Update a user by ID
// @Description Update a user with the provided payload
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body User true "User object"
// @Success 200 {object} User "User updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload or validation errors"
// @Failure 500 {object} ErrorResponse "Failed to update user"
// @Router /user/{id} [put]
func (api *Api) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storUser, ok := ctx.Value("user").(*stor.User)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Parse and validate the request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the user struct using the validator
	err = validate.Struct(user)
	if err != nil {
		// If validation fails, return the validation errors
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(user)

	storUserConverted := ConvertUserToUserStor(user, storUser)

	fmt.Println(storUserConverted)

	err = api.stor.UpdateUser(storUserConverted)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	user = *ConvertUserFromUserStor(*storUserConverted)

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the user as JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// @Summary Delete a user by ID
// @Description Delete a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 "User deleted successfully"
// @Failure 500 {object} ErrorResponse "Failed to delete user"
// @Router /user/{id} [delete]
func (api *Api) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	storUser, ok := ctx.Value("user").(*stor.User)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err := api.stor.DeleteUser(storUser)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
