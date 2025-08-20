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
	"github.com/stretchr/testify/assert"
)

func TestClientHandler(t *testing.T) {
	// Initialize router
	r := chi.NewRouter()
	r.Group(testapi.Router)

	// Generate the bearer token for the client
	tokenURL := "/api/auth"
	tokenData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {"lcp-server"},
		"client_secret": {"secret-123"},
	}
	tokenReq, err := http.NewRequest("POST", tokenURL, strings.NewReader(tokenData.Encode()))
	assert.NoError(t, err)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenRecorder := httptest.NewRecorder()
	r.ServeHTTP(tokenRecorder, tokenReq)
	if !assert.Equal(t, http.StatusOK, tokenRecorder.Code) {
		t.FailNow()
	}
	// Retrieve the access token from the response
	var tokenResp struct {
		Token string `json:"access_token"`
	}
	err = json.Unmarshal(tokenRecorder.Body.Bytes(), &tokenResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenResp.Token)

	// init a test user
	newUser := &stor.User{
		UUID:       gofakeit.UUID(),
		Name:       "Pierre ler",
		Email:      gofakeit.Email(),
		Password:   "password",
		TextHint:   "hint",
		Passphrase: "passphrase",
		SessionId:  gofakeit.UUID(),
	}

	newUserBytes, err := json.Marshal(newUser)
	assert.NoError(t, err)

	// create a user with a client token
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(newUserBytes))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	// count users
	userCount, err := testapi.Store.CountUsers()
	assert.NoError(t, err)
	assert.Equal(t, 1, int(userCount))

	// delete the user
	req = httptest.NewRequest("DELETE", "/api/users/"+newUser.UUID, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
