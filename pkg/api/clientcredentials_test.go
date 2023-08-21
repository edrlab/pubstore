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

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestClientHandler(t *testing.T) {
	// Initialize router
	r := chi.NewRouter()
	r.Group(api.Router)

	// Generate the bearer token by making a POST request to /api/v1/auth
	tokenURL := "/api/v1/auth"
	tokenData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {"lcp-server"},
		"client_secret": {"secret-123"},
	}
	fmt.Println(tokenData.Encode())
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

	// Try creating a publication with a token
	// Test POST /api/v1/publicationss
	req, err := http.NewRequest("POST", "/api/v1/publications", bytes.NewBuffer([]byte(jsonData)))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	// Unmarshal the response to get the created publication
	var createdPublication stor.Publication
	err = json.Unmarshal(recorder.Body.Bytes(), &createdPublication)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdPublication.UUID)

}
