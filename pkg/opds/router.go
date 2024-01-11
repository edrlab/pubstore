// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package opds

import (
	"context"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
)

// Router creates OPDS routes
func (o *Opds) Router(r chi.Router) {

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
		o.Config.OAuthSeed,
		time.Second*3600,
		&auth.UserVerifier{Store: o.Store},
		nil)

	r.Post("/opds/token", s.UserCredentials)
	r.Get("/401", GetAuthenticationDoc)
	r.Route("/opds", func(r chi.Router) {
		r.Get("/catalog", o.GetCatalog)
		r.Route("/publication/{id}", func(r chi.Router) {
			// TODO: check the use of authorizePassthrough
			r.Use(authorizePassthrough(o.Config.OAuthSeed, nil))
			r.Use(o.publicationCtx)
			r.Get("/", o.GetPublication)
			r.Get("/loan", o.GetPublicationLoan)
			r.Get("/borrow", o.GetPublicationBorrow)
			r.Get("/license", o.GetPublicationLicense)
		})
		r.Group(func(r chi.Router) {
			// TODO: check why a custom authorize is required
			r.Use(authorize(o.Config.OAuthSeed, nil))
			r.Get("/bookshelf", o.GetBookshelf)
		})
	})

}

type key int

var pubKey key

func (o *Opds) publicationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pubID := chi.URLParam(r, "id")
		pub, err := o.Store.GetPublication(pubID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), pubKey, pub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
