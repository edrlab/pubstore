// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package auth

import (
	"errors"
	"net/http"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/oauth"
	"golang.org/x/crypto/bcrypt"
)

type UserVerifier struct {
	*stor.Store
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (v *UserVerifier) ValidateUser(username, password, scope string, r *http.Request) error {
	var err error
	var user *stor.User
	user, err = v.Store.GetUserByEmail(username)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (*UserVerifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	if clientID == "lcp-server" && clientSecret == "secret-123" {
		return nil
	}
	return errors.New("wrong client")
}

// ValidateCode validates token ID
func (*UserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token
func (*UserVerifier) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	return claims, nil
}

// AddProperties provides additional information to the token response
func (*UserVerifier) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	return props, nil
}

// ValidateTokenID validates token ID
func (*UserVerifier) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*UserVerifier) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
