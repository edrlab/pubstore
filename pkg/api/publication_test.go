package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestPublicationHandler(t *testing.T) {
	// Initialize router
	r := chi.NewRouter()
	r.Group(api.Rooter)

	jsonData := `
	{
	    "title": "Test Publication",
	    "datePublication": "2023-06-16T12:00:00Z",
	    "description": "Test description",
	    "coverUrl": "http://example.com/cover.jpg",
	    "language": [
	        {
	            "code": "en"
	        },
	        {
	            "code": "fr"
	        }
	    ],
	    "publisher": [
	        {
	            "name": "Test Publisher A"
	        },
	        {
	            "name": "Test Publisher B"
	        }
	    ],
	    "author": [
	        {
	            "name": "Test Author A"
	        },
	        {
	            "name": "Test Author B"
	        }
	    ],
	    "category": [
	        {
	            "name": "Test Category A"
	        },
	        {
	            "name": "Test Category B"
	        }
	    ]
	}
	`

	sessionId := uuid.New().String()
	// Create a new user for testing
	createdUser := &stor.User{
		UUID:        gofakeit.UUID(),
		Name:        "Pierre ler",
		Email:       gofakeit.Email(),
		Pass:        "password123",
		LcpHintMsg:  "Hint",
		LcpPassHash: "hash123",
		SessionId:   sessionId,
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
	fmt.Println(tokenRecorder.Body.String())

	// Retrieve the access token from the response
	var tokenResp struct {
		Token string `json:"access_token"`
	}
	err = json.Unmarshal(tokenRecorder.Body.Bytes(), &tokenResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenResp.Token)
	fmt.Println(tokenResp.Token)

	req, err := http.NewRequest("POST", "/api/v1/publication", bytes.NewBuffer([]byte(jsonData)))
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	assert.NoError(t, err)
	req, err = http.NewRequest("POST", "/api/v1/publication", bytes.NewBuffer([]byte(jsonData)))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	// Unmarshal the response to get the created user
	var createPublicationStor stor.Publication
	err = json.Unmarshal(recorder.Body.Bytes(), &createPublicationStor)
	assert.NoError(t, err)
	assert.NotEmpty(t, createPublicationStor.UUID)
	assert.Equal(t, "Test Publication", createPublicationStor.Title)
	assert.Equal(t, "Test description", createPublicationStor.Description)

	// Test GET /api/v1/publication/{id}
	getUserURL := "/api/v1/publication/" + createPublicationStor.UUID
	req, err = http.NewRequest("GET", getUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var retrievedPub stor.Publication
	err = json.Unmarshal(recorder.Body.Bytes(), &retrievedPub)
	assert.NoError(t, err)
	// Check the retrieved user details

	// Test PUT /api/v1/user/{id}
	updateUserURL := "/api/v1/publication/" + createPublicationStor.UUID
	updateUserData := map[string]interface{}{
		"title": "Jane Doe",
	}
	updateUserDataBytes, err := json.Marshal(updateUserData)
	assert.NoError(t, err)
	req, err = http.NewRequest("PUT", updateUserURL, bytes.NewBuffer(updateUserDataBytes))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Retrieve publication by ID and validate updated name
	pubGetFromStor, err := api.stor.GetPublicationByUUID(createPublicationStor.UUID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", pubGetFromStor.Title)

	// Test DELETE /api/v1/publication/{id}
	deleteUserURL := "/api/v1/publication/" + createPublicationStor.UUID
	req, err = http.NewRequest("DELETE", deleteUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Retrieve user by ID and ensure it's not found
	pubDeleteFromStor, err := api.stor.GetPublicationByUUID(createdUser.UUID)
	assert.Error(t, err)
	assert.Nil(t, pubDeleteFromStor)
}
