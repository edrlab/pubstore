// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

// TODO: check why this code is needed. It is a copy of https://github.com/go-chi/oauth/blob/master/middleware.go
// with a copy-paste for a "passthrough" variant.

package opds

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/oauth"
)

type contextKey string

const (
	CredentialContext  contextKey = "oauth.credential"
	ClaimsContext      contextKey = "oauth.claims"
	ScopeContext       contextKey = "oauth.scope"
	TokenTypeContext   contextKey = "oauth.tokentype"
	AccessTokenContext contextKey = "oauth.accesstoken"
)

// Authorize is the OAuth 2.0 middleware for go-chi resource server.
// Authorize creates a BearerAuthentication middleware and return the Authorize method.
func authorize(secretKey string, formatter oauth.TokenSecureFormatter) func(next http.Handler) http.Handler {
	return newBearerAuthentication(secretKey, formatter).Authorize
}

// BearerAuthentication middleware for go-chi
type BearerAuthentication struct {
	secretKey string
	provider  *oauth.TokenProvider
}

// NewBearerAuthentication create a BearerAuthentication middleware
func newBearerAuthentication(secretKey string, formatter oauth.TokenSecureFormatter) *BearerAuthentication {
	ba := &BearerAuthentication{secretKey: secretKey}
	if formatter == nil {
		formatter = oauth.NewSHA256RC4TokenSecurityProvider([]byte(secretKey))
	}
	ba.provider = oauth.NewTokenProvider(formatter)
	return ba
}

// Authorize verifies the bearer token authorizing or not the request.
// Token is retrieved from the Authorization HTTP header that respects the format
// Authorization: Bearer {access_token}
func (ba *BearerAuthentication) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token, err := ba.checkAuthorizationHeader(auth)
		if err != nil {
			GetAuthenticationDoc(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CredentialContext, token.Credential)
		ctx = context.WithValue(ctx, ClaimsContext, token.Claims)
		ctx = context.WithValue(ctx, ScopeContext, token.Scope)
		ctx = context.WithValue(ctx, TokenTypeContext, token.TokenType)
		ctx = context.WithValue(ctx, AccessTokenContext, auth[7:])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Check header and token.
func (ba *BearerAuthentication) checkAuthorizationHeader(auth string) (t *oauth.Token, err error) {
	if len(auth) < 7 {
		return nil, errors.New("invalid bearer authorization header")
	}
	authType := strings.ToLower(auth[:6])
	if authType != "bearer" {
		return nil, errors.New("invalid bearer authorization header")
	}
	token, err := ba.provider.DecryptToken(auth[7:])
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if time.Now().UTC().After(token.CreationDate.Add(token.ExpiresIn)) {
		return nil, errors.New("token expired")
	}
	return token, nil
}

// -----------------------------------------------------------------

// BearerAuthentication middleware for go-chi
type BearerAuthenticationPassthrough struct {
	secretKey string
	provider  *oauth.TokenProvider
}

// NewBearerAuthentication create a BearerAuthentication middleware
func newBearerAuthenticationPassthrough(secretKey string, formatter oauth.TokenSecureFormatter) *BearerAuthenticationPassthrough {
	ba := &BearerAuthenticationPassthrough{secretKey: secretKey}
	if formatter == nil {
		formatter = oauth.NewSHA256RC4TokenSecurityProvider([]byte(secretKey))
	}
	ba.provider = oauth.NewTokenProvider(formatter)
	return ba
}

// Authorize is the OAuth 2.0 middleware for go-chi resource server.
// Authorize creates a BearerAuthentication middleware and return the Authorize method.
func authorizePassthrough(secretKey string, formatter oauth.TokenSecureFormatter) func(next http.Handler) http.Handler {
	return newBearerAuthenticationPassthrough(secretKey, formatter).AuthorizePassthrough
}

// Authorize verifies the bearer token authorizing or not the request.
// Token is retrieved from the Authorization HTTP header that respects the format
// Authorization: Bearer {access_token}
func (ba *BearerAuthenticationPassthrough) AuthorizePassthrough(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token, err := ba.checkAuthorizationHeaderPassthrough(auth)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CredentialContext, token.Credential)
		ctx = context.WithValue(ctx, ClaimsContext, token.Claims)
		ctx = context.WithValue(ctx, ScopeContext, token.Scope)
		ctx = context.WithValue(ctx, TokenTypeContext, token.TokenType)
		ctx = context.WithValue(ctx, AccessTokenContext, auth[7:])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Check header and token.
func (ba *BearerAuthenticationPassthrough) checkAuthorizationHeaderPassthrough(auth string) (t *oauth.Token, err error) {
	if len(auth) < 7 {
		return nil, errors.New("invalid bearer authorization header")
	}
	authType := strings.ToLower(auth[:6])
	if authType != "bearer" {
		return nil, errors.New("invalid bearer authorization header")
	}
	token, err := ba.provider.DecryptToken(auth[7:])
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if time.Now().UTC().After(token.CreationDate.Add(token.ExpiresIn)) {
		return nil, errors.New("token expired")
	}
	return token, nil
}
