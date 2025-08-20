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

func TestPublicationHandler(t *testing.T) {
	// Initialize the router
	r := chi.NewRouter()
	r.Group(testapi.Router)

	// init a publication
	newPublication := &stor.Publication{
		UUID:          gofakeit.UUID(),
		Title:         "Test Publication",
		ContentType:   "application/epub+zip",
		DatePublished: "2022-12-31",
		Description:   "Test description",
		CoverUrl:      "http://example.com/cover.jpg",
		Language: []stor.Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []stor.Publisher{
			{Name: "Publisher A"},
			{Name: "Publisher B"},
		},
		Author: []stor.Author{
			{Name: "Author A"},
			{Name: "Author B"},
		},
		Category: []stor.Category{
			{Name: "Category A"},
			{Name: "Category B"},
		},
	}

	pubBytes, err := json.Marshal(newPublication)
	assert.NoError(t, err)

	// try creating a publication with no token
	req := httptest.NewRequest("POST", "/api/publications", bytes.NewBuffer([]byte(pubBytes)))
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.NoError(t, err)

	// init a new user for testing
	newUser := &stor.User{
		Name:       "Albert ler",
		Email:      gofakeit.Email(),
		Password:   "password",
		TextHint:   "hint",
		Passphrase: "passphrase",
	}
	// create the user in the database
	err = testapi.CreateUser(newUser)
	assert.NoError(t, err)

	// generate the bearer token
	tokenURL := "/api/token"
	tokenData := url.Values{
		"grant_type": {"password"},
		"username":   {newUser.Email},
		"password":   {newUser.Password},
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

	// create a publication
	req = httptest.NewRequest("POST", "/api/publications", bytes.NewBuffer([]byte(pubBytes)))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	if !assert.Equal(t, http.StatusCreated, recorder.Code) {
		t.FailNow()
	}

	// unmarshal the response to get the newly created publication
	var createdPublication stor.Publication
	err = json.Unmarshal(recorder.Body.Bytes(), &createdPublication)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdPublication.UUID)
	assert.Equal(t, "Test Publication", createdPublication.Title)
	assert.Equal(t, "Test description", createdPublication.Description)

	// get the publication previously created by its id
	getPubURL := "/api/publications/" + newPublication.UUID
	req = httptest.NewRequest("GET", getPubURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var retrievedPub stor.Publication
	err = json.Unmarshal(recorder.Body.Bytes(), &retrievedPub)
	assert.NoError(t, err)

	// check the retrieved pub details
	assert.Equal(t, newPublication.UUID, retrievedPub.UUID)
	assert.Equal(t, newPublication.Title, retrievedPub.Title)
	assert.Equal(t, newPublication.ContentType, retrievedPub.ContentType)
	assert.Equal(t, newPublication.DatePublished, retrievedPub.DatePublished)
	assert.Equal(t, newPublication.Description, retrievedPub.Description)
	assert.Equal(t, newPublication.CoverUrl, retrievedPub.CoverUrl)
	assert.Equal(t, newPublication.Language[0].Code, retrievedPub.Language[0].Code)
	assert.Equal(t, newPublication.Publisher[0].Name, retrievedPub.Publisher[0].Name)
	assert.Equal(t, newPublication.Author[0].Name, retrievedPub.Author[0].Name)
	assert.Equal(t, newPublication.Category[0].Name, retrievedPub.Category[0].Name)

	// update the publication
	updatePubURL := "/api/publications/" + newPublication.UUID
	newPublication.Title = "Update Test Publication"
	updatePubBytes, err := json.Marshal(newPublication)
	assert.NoError(t, err)
	req = httptest.NewRequest("PUT", updatePubURL, bytes.NewBuffer(updatePubBytes))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// retrieve the publication by ID and validate updated title
	pubFromStor, err := testapi.Store.GetPublication(newPublication.UUID)
	assert.NoError(t, err)
	assert.Equal(t, "Update Test Publication", pubFromStor.Title)

	// create a second publication
	newPublication.UUID = gofakeit.UUID()
	newPublication.Title = "Test Publication 2"
	newPublication.ContentType = "application/pdf"
	pubBytes, err = json.Marshal(newPublication)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/api/publications", bytes.NewBuffer([]byte(pubBytes)))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	if !assert.Equal(t, http.StatusCreated, recorder.Code) {
		t.FailNow()
	}

	// list publications
	req = httptest.NewRequest("GET", "/api/publications?page=1&pageSize=5", bytes.NewBuffer([]byte(pubBytes)))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	if !assert.Equal(t, http.StatusOK, recorder.Code) {
		t.FailNow()
	}

	var retrievedPubs []stor.Publication
	err = json.Unmarshal(recorder.Body.Bytes(), &retrievedPubs)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(retrievedPubs))

	// search publications
	req = httptest.NewRequest("GET", "/api/publications/search?format=epub&page=1&pageSize=5", bytes.NewBuffer([]byte(pubBytes)))
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	if !assert.Equal(t, http.StatusOK, recorder.Code) {
		t.FailNow()
	}

	err = json.Unmarshal(recorder.Body.Bytes(), &retrievedPubs)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(retrievedPubs))

	// delete the publication
	deleteUserURL := "/api/publications/" + newPublication.UUID
	req = httptest.NewRequest("DELETE", deleteUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// retrieve the publication by ID and ensure it's not found
	_, err = testapi.Store.GetPublication(newPublication.UUID)
	assert.Error(t, err)

	// delete the test user
	err = testapi.Store.DeleteUser(newUser)
	assert.NoError(t, err)
}
