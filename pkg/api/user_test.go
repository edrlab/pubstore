// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

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

func TestUserHandler(t *testing.T) {
	// Initialize the router
	r := chi.NewRouter()
	r.Group(testapi.Router)

	// init an admin user, who will be able to create other users
	// this user has no initial UUID and no session id
	adminUser := &stor.User{
		Name:       "Admin",
		Email:      gofakeit.Email(),
		Password:   "password",
		TextHint:   "hint",
		Passphrase: "passphrase",
	}

	// create the user in the database, so that a token can be acquired
	// note: the creation replace the clear password by its hash
	err := testapi.Store.CreateUser(adminUser)
	assert.NoError(t, err)

	// generate a bearer token
	tokenURL := "/api/v1/token"
	tokenData := url.Values{
		"grant_type": {"password"},
		"username":   {adminUser.Email},
		"password":   {adminUser.Password},
	}
	tokenReq := httptest.NewRequest("POST", tokenURL, strings.NewReader(tokenData.Encode()))
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenRecorder := httptest.NewRecorder()
	r.ServeHTTP(tokenRecorder, tokenReq)
	if !assert.Equal(t, http.StatusOK, tokenRecorder.Code) {
		t.FailNow()
	}

	// retrieve the access token from the response
	var tokenResp struct {
		Token string `json:"access_token"`
	}
	err = json.Unmarshal(tokenRecorder.Body.Bytes(), &tokenResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenResp.Token)

	// init a test user
	// this user has an initial UUID
	newUser := &stor.User{
		UUID:       gofakeit.UUID(),
		Name:       gofakeit.Name(),
		Email:      gofakeit.Email(),
		Password:   "password",
		TextHint:   "hint",
		Passphrase: "passphrase",
	}
	userBytes, err := json.Marshal(newUser)
	assert.NoError(t, err)

	// try creating the user with no token
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(userBytes))
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.NoError(t, err)

	// create the user with a token
	req = httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(userBytes))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	if !assert.Equal(t, http.StatusCreated, recorder.Code) {
		t.FailNow()
	}

	// unmarshal the response to get the created user
	var createdUser stor.User
	err = json.Unmarshal(recorder.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdUser.UUID)

	// get the user previously created by its id
	getUserURL := "/api/v1/users/" + newUser.UUID
	req = httptest.NewRequest("GET", getUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var retrievedUser stor.User
	err = json.Unmarshal(recorder.Body.Bytes(), &retrievedUser)
	assert.NoError(t, err)

	// check the retrieved user details
	assert.NotEqual(t, "", retrievedUser.UUID)
	assert.Equal(t, newUser.Name, retrievedUser.Name)
	assert.Equal(t, newUser.Email, retrievedUser.Email)
	assert.Equal(t, newUser.TextHint, retrievedUser.TextHint)
	// these are not returned (filtered on rendering)
	assert.Equal(t, "", retrievedUser.Password)
	assert.Equal(t, "", retrievedUser.HPassword)
	assert.Equal(t, "", retrievedUser.Passphrase)
	assert.Equal(t, "", retrievedUser.HPassphrase)
	assert.Equal(t, "", retrievedUser.SessionId)

	// update the user
	updateUserURL := "/api/v1/users/" + newUser.UUID
	newUser.Name = "Jane Doe"
	updateUserBytes, err := json.Marshal(newUser)
	assert.NoError(t, err)
	req = httptest.NewRequest("PUT", updateUserURL, bytes.NewBuffer(updateUserBytes))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// retrieve the user by ID and validate the updated name
	userFromStor, err := testapi.Store.GetUser(newUser.UUID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", userFromStor.Name)

	// list users
	listUserURL := "/api/v1/users"
	req = httptest.NewRequest("GET", listUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var retrievedUsers []stor.User
	err = json.Unmarshal(recorder.Body.Bytes(), &retrievedUsers)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedUsers))
	assert.Equal(t, "Admin", retrievedUsers[0].Name)
	assert.Equal(t, "Jane Doe", retrievedUsers[1].Name)

	// delete the user
	deleteUserURL := "/api/v1/users/" + newUser.UUID
	req = httptest.NewRequest("DELETE", deleteUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// retrieve user by ID and ensure it's not found
	_, err = testapi.Store.GetUser(newUser.UUID)
	assert.Error(t, err)
}
