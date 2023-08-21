package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/config"
	_ "github.com/edrlab/pubstore/pkg/docs"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

type Api struct {
	stor *stor.Stor
}

func Init(s *stor.Stor) *Api {
	return &Api{stor: s}
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

// @title Pubstore API
// @version 1.0
// @description Pubstore API.

// @contact.name edrlab
// @contact.url https://edrlab.org
// @contact.email support@edrlab.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host pubstore.edrlab.org
// @BasePath /api/v1

//	@securitydefinitions.oauth2.password	OAuth2Password
//	@tokenUrl								https://pubstore.edrlab.org/api/v1/token
//	@scope.read								Grants read access
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information:w

func (api *Api) Router(r chi.Router) {

	validate = validator.New()

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
		time.Second*120,
		&UserVerifier{stor: api.stor},
		nil)

	r.Get("/api/v1/swagger/*", httpSwagger.WrapHandler)

	credentials := make(map[string]string)
	credentials[config.PUBSTORE_USERNAME] = config.PUBSTORE_PASSWORD

	// Create a publication using basic auth (used by the LCP encryption tool)
	r.Route("/api/v1/notify", func(r chi.Router) {
		r.Use(middleware.BasicAuth("restricted", credentials))
		r.Post("/", api.createPublicationHandler)
	})

	/*
		 Generate Token using username & password
			    	POST http://localhost:8080/api/v1/token
					Content-Type: application/x-www-form-urlencoded
					grant_type=password&username=user01&password=12345
	*/
	/*
		RefreshTokenGrant Token
			POST http://localhost:8080/api/v1/token
			Content-Type: application/x-www-form-urlencoded
			grant_type=refresh_token&refresh_token={the refresh_token obtained in the previous response}
	*/
	r.Post("/api/v1/token", s.UserCredentials)
	r.Post("/api/v1/auth", s.ClientCredentials)

	r.Route("/api/v1/publications", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(oauth.Authorize(config.OauthSeed, nil))
			r.Post("/", api.createPublicationHandler)
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Use(api.publicationCtx)
			r.Get("/", api.getPublicationHandler)
			r.Group(func(r chi.Router) {
				r.Use(oauth.Authorize(config.OauthSeed, nil))
				r.Put("/", api.updatePublicationHandler)
				r.Delete("/", api.deletePublicationHandler)
			})
		})
	})
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(oauth.Authorize(config.OauthSeed, nil))
			r.Post("/", api.createUserHandler)
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Use(api.userCtx)
			r.Get("/", api.getUserHandler)
			r.Group(func(r chi.Router) {
				r.Use(oauth.Authorize(config.OauthSeed, nil))
				r.Put("/", api.updateUserHandler)
				r.Delete("/", api.deleteUserHandler)
			})
		})
	})
}

func (api *Api) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		user, err := api.stor.GetUserByUUID(userID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *Api) publicationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pubID := chi.URLParam(r, "id")
		pub, err := api.stor.GetPublicationByUUID(pubID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "publication", pub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
