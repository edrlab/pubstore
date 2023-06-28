package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	"github.com/go-playground/validator/v10"
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

func (api *Api) Rooter(r chi.Router) {

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
		"Edrlab-Rocks",
		time.Second*120,
		&UserVerifier{stor: api.stor},
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
	r.Post("/api/v1/token", s.UserCredentials)

	r.Route("/api/v1/publication", func(publicationRouter chi.Router) {
		publicationRouter.Group(func(postRouter chi.Router) {
			postRouter.Use(oauth.Authorize("Edrlab-Rocks", nil))
			postRouter.Post("/", api.createPublicationHandler)
		})
		publicationRouter.Route("/{id}", func(idRouter chi.Router) {
			idRouter.Use(api.publicationCtx)
			idRouter.Get("/", api.getPublicationHandler)
			idRouter.Group(func(idRouterGroup chi.Router) {
				idRouterGroup.Use(oauth.Authorize("Edrlab-Rocks", nil))
				idRouterGroup.Put("/", api.updatePublicationHandler)
				idRouterGroup.Delete("/", api.deletePublicationHandler)
			})
		})
	})
	r.Route("/api/v1/user", func(userRouter chi.Router) {
		userRouter.Group(func(postRouter chi.Router) {
			postRouter.Use(oauth.Authorize("Edrlab-Rocks", nil))
			postRouter.Post("/", api.createUserHandler)
		})
		userRouter.Route("/{id}", func(idRouter chi.Router) {
			idRouter.Use(api.userCtx)
			idRouter.Get("/", api.getUserHandler)
			idRouter.Group(func(idRouterGroup chi.Router) {
				idRouterGroup.Use(oauth.Authorize("Edrlab-Rocks", nil))
				idRouterGroup.Put("/", api.updateUserHandler)
				idRouterGroup.Delete("/", api.deleteUserHandler)
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
