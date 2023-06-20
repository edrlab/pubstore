package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	// Initialize router
	r := chi.NewRouter()
	r.Post("/api/v1/user", api.createUserHandler)
	r.Get("/api/v1/user/{id}", api.getUserHandler)

	sessionId := uuid.New().String()
	// Create a new user for testing
	newUser := &User{
		Name:        "John Doe",
		Email:       gofakeit.Email(),
		Pass:        "password123",
		LcpHintMsg:  "Hint",
		LcpPassHash: "hash123",
		SessionId:   sessionId,
	}

	// Encode the user as JSON
	userJSON, err := json.Marshal(newUser)
	assert.NoError(t, err)

	// Perform a POST request to create the user
	createReq := httptest.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(userJSON))
	createRec := httptest.NewRecorder()
	r.ServeHTTP(createRec, createReq)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, createRec.Code)

	// Perform a GET request to retrieve the user
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/user/%s", sessionId), nil)
	getRec := httptest.NewRecorder()
	r.ServeHTTP(getRec, getReq)

	// Check the response status code
	assert.Equal(t, http.StatusOK, getRec.Code)

	// Decode the response body
	var retrievedUser User
	err = json.Unmarshal(getRec.Body.Bytes(), &retrievedUser)
	assert.NoError(t, err)

	// Check the retrieved user details
	assert.Equal(t, newUser.Name, retrievedUser.Name)
	assert.Equal(t, newUser.Email, retrievedUser.Email)
	assert.Equal(t, newUser.Pass, retrievedUser.Pass)
	assert.Equal(t, newUser.LcpHintMsg, retrievedUser.LcpHintMsg)
	assert.Equal(t, newUser.LcpPassHash, retrievedUser.LcpPassHash)
	assert.Equal(t, newUser.SessionId, retrievedUser.SessionId)
}
