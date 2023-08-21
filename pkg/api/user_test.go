package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	// Initialize router
	r := chi.NewRouter()
	r.Group(api.Router)

	// Create a new user for testing
	createdUser := &stor.User{
		UUID:        gofakeit.UUID(),
		Name:        "Pierre ler",
		Email:       gofakeit.Email(),
		Pass:        "password123",
		LcpHintMsg:  "Hint",
		LcpPassHash: "hash123",
		SessionId:   uuid.New().String(),
	}

	// Create the user in the storage
	err := api.stor.CreateUser(createdUser)
	assert.NoError(t, err)

	// Generate the bearer token by making a POST request to /api/v1/token
	tokenURL := "/api/v1/token"
	tokenData := url.Values{
		"grant_type": {"password"},
		"username":   {createdUser.Email},
		"password":   {"password123"},
	}
	tokenReq, err := http.NewRequest("POST", tokenURL, strings.NewReader(tokenData.Encode()))
	assert.NoError(t, err)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenRecorder := httptest.NewRecorder()
	r.ServeHTTP(tokenRecorder, tokenReq)
	assert.Equal(t, http.StatusOK, tokenRecorder.Code)

	// Retrieve the access token from the response
	var tokenResp struct {
		Token string `json:"access_token"`
	}
	err = json.Unmarshal(tokenRecorder.Body.Bytes(), &tokenResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenResp.Token)

	// Try creating a user without any token
	newUser := &stor.User{
		Name:        "Jean Mène",
		Email:       gofakeit.Email(),
		Pass:        "password123",
		LcpHintMsg:  "Hint",
		LcpPassHash: "hash123",
		SessionId:   uuid.New().String(),
	}
	newUserBytes, err := json.Marshal(newUser)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(newUserBytes))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.NoError(t, err)

	// Try creating a user with a token
	// Test POST /api/v1/users
	req, err = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(newUserBytes))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	// Unmarshal the response to get the created user
	var createdUserFromPostRequest stor.User
	err = json.Unmarshal(recorder.Body.Bytes(), &createdUserFromPostRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdUserFromPostRequest.UUID)

	// Get the user previously created by its id
	// Test GET /api/v1/users/{id}
	getUserURL := "/api/v1/users/" + createdUser.UUID
	req, err = http.NewRequest("GET", getUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	var retrievedUser stor.User
	err = json.Unmarshal(recorder.Body.Bytes(), &retrievedUser)
	assert.NoError(t, err)

	// Check the retrieved user details
	assert.Equal(t, createdUser.Name, retrievedUser.Name)
	assert.Equal(t, createdUser.Email, retrievedUser.Email)
	assert.Equal(t, "", retrievedUser.Pass)
	assert.Equal(t, createdUser.LcpHintMsg, retrievedUser.LcpHintMsg)
	assert.Equal(t, createdUser.LcpPassHash, retrievedUser.LcpPassHash)
	assert.Equal(t, "", retrievedUser.SessionId)

	// Update the user
	// Test PUT /api/v1/users/{id}
	updateUserURL := "/api/v1/users/" + createdUser.UUID
	updateUserData := map[string]interface{}{
		"name": "Jane Doe",
	}
	updateUserDataBytes, err := json.Marshal(updateUserData)
	assert.NoError(t, err)
	req, err = http.NewRequest("PUT", updateUserURL, bytes.NewBuffer(updateUserDataBytes))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Retrieve user by ID and validate updated name
	userGetFromStor, err := api.stor.GetUserByUUID(createdUser.UUID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", userGetFromStor.Name)

	// Delete the user
	// Test DELETE /api/v1/users/{id}
	deleteUserURL := "/api/v1/users/" + createdUser.UUID
	req, err = http.NewRequest("DELETE", deleteUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Retrieve user by ID and ensure it's not found
	userDeleteFromStor, err := api.stor.GetUserByUUID(createdUser.UUID)
	assert.Error(t, err)
	assert.Nil(t, userDeleteFromStor)
}
