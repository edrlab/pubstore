// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/edrlab/pubstore/pkg/conf"
	_ "github.com/edrlab/pubstore/pkg/docs"
	"github.com/edrlab/pubstore/pkg/internal/auth"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Api struct {
	*conf.Config
	*stor.Store
}

func Init(c *conf.Config, s *stor.Store) Api {
	return Api{
		Config: c,
		Store:  s,
	}
}

// -------------------------------------------

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

func (a *Api) Router(r chi.Router) {

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
		a.Config.OAuthSeed,
		time.Second*120,
		&auth.UserVerifier{Store: a.Store},
		nil)

	r.Get("/api/v1/swagger/*", httpSwagger.WrapHandler)

	credentials := make(map[string]string)
	credentials[a.Config.UserName] = a.Config.Password

	// Create a publication using basic auth (used by the LCP encryption tool)
	r.Route("/api/v1/notify", func(r chi.Router) {
		r.Use(middleware.BasicAuth("restricted", credentials))
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post("/", a.createPublication)
	})

	/*
		 Generate a token using username & password
			POST http://localhost:8080/api/v1/token
			Content-Type: application/x-www-form-urlencoded
			grant_type=password&username=user01&password=12345
	*/
	/*
		Refresh a token
			POST http://localhost:8080/api/v1/token
			Content-Type: application/x-www-form-urlencoded
			grant_type=refresh_token&refresh_token={the refresh_token obtained in the previous response}
	*/
	r.Post("/api/v1/token", s.UserCredentials)
	r.Post("/api/v1/auth", s.ClientCredentials)

	r.Route("/api/v1/publications", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.With(paginate).Get("/", a.listPublications)
		r.With(paginate).Get("/search", a.searchPublications)
		r.Group(func(r chi.Router) {
			r.Use(oauth.Authorize(a.Config.OAuthSeed, nil))
			r.Post("/", a.createPublication)
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Use(a.publicationId)
			r.Get("/", a.getPublication)
			r.Group(func(r chi.Router) {
				r.Use(oauth.Authorize(a.Config.OAuthSeed, nil))
				r.Put("/", a.updatePublication)
				r.Delete("/", a.deletePublication)
			})
		})
	})
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.With(paginate).Get("/", a.listUsers)
		r.Group(func(r chi.Router) {
			r.Use(oauth.Authorize(a.Config.OAuthSeed, nil))
			r.Post("/", a.createUser)
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Use(a.userId)
			r.Get("/", a.getUser)
			r.Group(func(r chi.Router) {
				r.Use(oauth.Authorize(a.Config.OAuthSeed, nil))
				r.Put("/", a.updateUser)
				r.Delete("/", a.deleteUser)
			})
		})
	})

	// License gateway
	r.Route("/api/v1/licenses", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Route("/{id}", func(r chi.Router) {
			r.Use(a.licenseId)
			r.Get("/", a.getFreshLicense)
		})
	})

}

// userId middleware,
func (a *Api) userId(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		user, err := a.Store.GetUser(userID)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := newUserContext(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// publicationId middleware
func (a *Api) publicationId(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pubID := chi.URLParam(r, "id")
		pub, err := a.Store.GetPublication(pubID)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := newPubContext(r.Context(), pub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// licenseId middleware
func (a *Api) licenseId(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		licID := chi.URLParam(r, "id")
		trans, err := a.Store.GetTransactionByLicence(licID)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := newTransContext(r.Context(), trans)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Pagination defines a page and pageSize, used for api pagination
type Pagination struct {
	Page     int
	PageSize int
}

// paginate middleware,
func paginate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p Pagination
		var err error

		p.Page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			p.Page = 1 // default value
		}
		if p.Page < 1 {
			render.Render(w, r, ErrInvalidRequest(errors.New("page must be 1 or more")))
			return
		}
		p.PageSize, err = strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil {
			p.PageSize = 1000 // default value (limits the count of retrieved items)
		}
		if p.PageSize < 1 {
			render.Render(w, r, ErrInvalidRequest(errors.New("pageSize must be 1 or more")))
			return
		}

		ctx := newPaginateContext(r.Context(), &p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type key int

var userKey key
var pubKey key
var paginateKey key

// newUserContext returns a new Context that carries a User.
func newUserContext(ctx context.Context, u *stor.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// fromUserContext returns the User value stored in ctx, if any.
func fromUserContext(ctx context.Context) *stor.User {
	u, _ := ctx.Value(userKey).(*stor.User)
	return u
}

// newPubContext returns a new Context that carries a Publication.
func newPubContext(ctx context.Context, p *stor.Publication) context.Context {
	return context.WithValue(ctx, pubKey, p)
}

// fromPubContext returns the Publication value stored in ctx, if any.
func fromPubContext(ctx context.Context) *stor.Publication {
	p, _ := ctx.Value(pubKey).(*stor.Publication)
	return p
}

// newTransContext returns a new Context that carries a Transaction.
func newTransContext(ctx context.Context, p *stor.Transaction) context.Context {
	return context.WithValue(ctx, pubKey, p)
}

// fromTransContext returns the Transaction value stored in ctx, if any.
func fromTransContext(ctx context.Context) *stor.Transaction {
	p, _ := ctx.Value(pubKey).(*stor.Transaction)
	return p
}

// newPaginateContext returns a new Context that carries a Pagination.
func newPaginateContext(ctx context.Context, p *Pagination) context.Context {
	return context.WithValue(ctx, paginateKey, p)
}

// fromPaginateContext returns the Pagination value stored in ctx, if any.
func fromPaginateContext(ctx context.Context) *Pagination {
	p, _ := ctx.Value(paginateKey).(*Pagination)
	return p
}
