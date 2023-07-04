package opds

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/edrlab/pubstore/pkg/config"
	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	"golang.org/x/crypto/bcrypt"
)

type Opds struct {
	stor *stor.Stor
}

func Init(s *stor.Stor) *Opds {
	return &Opds{stor: s}
}

// TestUserVerifier provides user credentials verifier for testing.
type UserVerifier struct {
	stor *stor.Stor
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (u *UserVerifier) ValidateUser(username, password, scope string, r *http.Request) error {
	user, err := u.stor.GetUserByEmail(username)
	if err == nil && bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(password)) == nil {
		return nil
	}

	return errors.New("wrong user")
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (*UserVerifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
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

func (opds *Opds) catalogHandler(w http.ResponseWriter, r *http.Request) {

	opdsFeed, err := opds.GenerateOpdsFeed(1, 50)
	if err != nil {
		fmt.Fprintf(w, "opds feed : %v!", err)
	}
	// Encode the publication as JSON and write it to the response
	w.Header().Set("Content-Type", "application/opds+json")
	err = json.NewEncoder(w).Encode(opdsFeed)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (opds *Opds) publicationHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	credential, ok := ctx.Value(CredentialContext).(string)
	authentified := false
	if ok {
		authentified = true
	}

	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	pub, err := convertToPublication(storPublication)
	if err != nil {
		http.Error(w, "Failed to convert publication", http.StatusInternalServerError)
		return
	}

	defer func() {
		// Encode the publication as JSON and write it to the response
		w.Header().Set("Content-Type", "application/opds+json")
		err = json.NewEncoder(w).Encode(pub)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}()

	if authentified {
		// add authentified acquisition link

		user, err := opds.stor.GetUserByEmail(credential)
		if err != nil {
			// http.Error(w, "Failed to get user", http.StatusInternalServerError)
			fmt.Println("Failed to get user : " + credential)
			pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("notAuthentified", storPublication.UUID, "", "", time.Time{}, time.Time{}))
			return
		}

		transactions, err := opds.getTransactionFromUserAndPubUUID(user, storPublication.UUID)
		if err != nil {
			pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("authentified", storPublication.UUID, "", "", time.Time{}, time.Time{}))
			return
		}

		lsdStatus, err := lcp.GetLsdStatus(transactions.LicenceId, user.Email, user.LcpHintMsg, user.LcpPassHash)
		if err != nil {
			lsdStatus = &lcp.LsdStatus{}
		}
		pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("authentifiedAndBorrowed", storPublication.UUID, lsdStatus.StatusCode, user.LcpPassHash, lsdStatus.StartDate, lsdStatus.EndDate))

	} else {
		// add borrow link
		pub.Links = append(pub.Links, publicationAcquisitionLinkChoice("notAuthentified", storPublication.UUID, "", "", time.Time{}, time.Time{}))
	}
}

func (opds *Opds) publicationBorrowHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(CredentialContext).(string)
	if !ok {
		opdsAuthenticationHandler(w, r)
		return
	}
	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Redirect(w, r, "/opds/v1/publication/"+storPublication.UUID, http.StatusFound)
}

func (opds *Opds) publicationLoanHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	credential, ok := ctx.Value(CredentialContext).(string)
	if !ok {
		opdsAuthenticationHandler(w, r)
		return
	}
	user, err := opds.stor.GetUserByEmail(credential)
	if err != nil {
		http.Error(w, "Failed to get transaction", http.StatusInternalServerError)
		return
	}

	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	licenceBytes, err := lcp.LicenceLoan(storPublication.UUID, user.UUID, user.Email, user.LcpHintMsg, user.LcpPassHash, 100, 2000, time.Now(), time.Now().AddDate(0, 0, 7))
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
	io.Copy(w, bytes.NewReader(licenceBytes))
}

func (opds *Opds) publicationFreshLicenseHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	credential, ok := ctx.Value(CredentialContext).(string)
	if !ok {
		opdsAuthenticationHandler(w, r)
		return
	}
	user, err := opds.stor.GetUserByEmail(credential)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	storPublication, ok := ctx.Value("publication").(*stor.Publication)
	if !ok {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	transaction, err := opds.stor.GetTransactionByUserAndPublication(user.ID, storPublication.ID)
	if err != nil {
		http.Error(w, "Failed to get transaction", http.StatusInternalServerError)
		return
	}

	licenceBytes, err := lcp.GenerateFreshLicenceFromLcpServer(transaction.LicenceId, user.Email, user.LcpHintMsg, user.LcpPassHash)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
	io.Copy(w, bytes.NewReader(licenceBytes))
}

func (opds *Opds) bookshelfHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	credential, ok := ctx.Value(CredentialContext).(string)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	opdsFeed, err := opds.GenerateBookshelfFeed(credential)
	if err != nil {
		fmt.Println("Bookshelf : " + err.Error())
	}

	// Encode the publication as JSON and write it to the response
	w.Header().Set("Content-Type", "application/opds+json")
	err = json.NewEncoder(w).Encode(opdsFeed)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (opds *Opds) Router(r chi.Router) {

	// https://github.com/go-chi/oauth/blob/master/example/authserver/main.go
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	s := oauth.NewBearerServer(
		config.OauthSeed,
		time.Second*3600,
		&UserVerifier{stor: opds.stor},
		nil)

	/*
		 Generate Token using username & password
			    	POST http://localhost:8080/token
					Content-Type: application/x-www-form-urlencoded
					grant_type=password&username=user01&password=12345
	*/
	/*
		RefreshTokenGrant Token
			POST http://localhost:8080/token
			Content-Type: application/x-www-form-urlencoded
			grant_type=refresh_token&refresh_token={the refresh_token obtained in the previous response}
	*/
	r.Post("/opds/v1/token", s.UserCredentials)
	r.Get("/401", opdsAuthenticationHandler)
	r.Route("/opds/v1", func(opdsRouter chi.Router) {
		opdsRouter.Get("/catalog", opds.catalogHandler)
		opdsRouter.Route("/publication/{id}", func(idRouter chi.Router) {
			idRouter.Use(authorizePassthrough(config.OauthSeed, nil))
			idRouter.Use(opds.publicationCtx)
			idRouter.Get("/", opds.publicationHandler)
			idRouter.Get("/loan", opds.publicationLoanHandler)
			idRouter.Get("/borrow", opds.publicationBorrowHandler)
			idRouter.Get("/license", opds.publicationFreshLicenseHandler)
		})
		opdsRouter.Group(func(opdsRouterAuth chi.Router) {
			opdsRouterAuth.Use(authorize(config.OauthSeed, nil))
			opdsRouterAuth.Get("/bookshelf", opds.bookshelfHandler)
		})
	})

}

func (opds *Opds) publicationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pubID := chi.URLParam(r, "id")
		pub, err := opds.stor.GetPublicationByUUID(pubID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "publication", pub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func opdsAuthenticationHandler(w http.ResponseWriter, _ *http.Request) {
	jsonData := `
	{
		"id": "org.edrlab.pubstore",
		"title": "LOGIN",
		"description": "PUBSTORE LOGIN",
		"links": [
		  {"rel": "logo", "href": "http://localhost:8080/static/images/edrlab-logo.jpeg", "type": "image/jpeg", "width": 90, "height": 90}
		],
		"authentication": [
		  {
			"type": "http://opds-spec.org/auth/oauth/password",
			"links": [
			  {"rel": "authenticate", "href": "http://localhost:8080/opds/v1/token", "type": "application/json"}
			]
		  }
		]
	  }`

	w.Header().Set("Content-Type", "application/opds-authentication+json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(jsonData))
}

type contextKey string

const (
	CredentialContext  contextKey = "oauth.credential"
	ClaimsContext      contextKey = "oauth.claims"
	ScopeContext       contextKey = "oauth.scope"
	TokenTypeContext   contextKey = "oauth.tokentype"
	AccessTokenContext contextKey = "oauth.accesstoken"
)

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

// Authorize is the OAuth 2.0 middleware for go-chi resource server.
// Authorize creates a BearerAuthentication middleware and return the Authorize method.
func authorize(secretKey string, formatter oauth.TokenSecureFormatter) func(next http.Handler) http.Handler {
	return newBearerAuthentication(secretKey, formatter).Authorize
}

// Authorize verifies the bearer token authorizing or not the request.
// Token is retrieved from the Authorization HTTP header that respects the format
// Authorization: Bearer {access_token}
func (ba *BearerAuthentication) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token, err := ba.checkAuthorizationHeader(auth)
		if err != nil {
			opdsAuthenticationHandler(w, r)
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
		return nil, errors.New("Invalid bearer authorization header")
	}
	authType := strings.ToLower(auth[:6])
	if authType != "bearer" {
		return nil, errors.New("Invalid bearer authorization header")
	}
	token, err := ba.provider.DecryptToken(auth[7:])
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	if time.Now().UTC().After(token.CreationDate.Add(token.ExpiresIn)) {
		return nil, errors.New("Token expired")
	}
	return token, nil
}

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
		return nil, errors.New("Invalid bearer authorization header")
	}
	authType := strings.ToLower(auth[:6])
	if authType != "bearer" {
		return nil, errors.New("Invalid bearer authorization header")
	}
	token, err := ba.provider.DecryptToken(auth[7:])
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	if time.Now().UTC().After(token.CreationDate.Add(token.ExpiresIn)) {
		return nil, errors.New("Token expired")
	}
	return token, nil
}
