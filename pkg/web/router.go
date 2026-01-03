// Copyright 2025 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

package web

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
)

type Web struct {
	*conf.Config
	*stor.Store
	*view.View
}

func Init(c *conf.Config, s *stor.Store, v *view.View) Web {

	// Configure goview to retrieve views at the proper location
	gvConf := goview.DefaultConfig
	gvConf.Root = filepath.Join(c.RootDir, "views")
	gv := goview.New(gvConf)
	goview.Use(gv)

	return Web{
		Config: c,
		Store:  s,
		View:   v,
	}
}

func (web *Web) Router(r chi.Router) {

	// Serve static files from the "static" directory. The '*' means that sub-routes are served from sub-directories
	filesDir := http.Dir(filepath.Join(web.Config.RootDir, "static"))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/index", http.StatusFound)
		})
		r.Get("/index", func(w http.ResponseWriter, r *http.Request) {
			userStor := web.getUserByCookie(r)
			userName := ""
			if userStor != nil {
				userName = userStor.Name
			}
			goviewModel := goview.M{
				"pageTitle":           "pubstore",
				"userIsAuthenticated": web.userIsAuthenticated(r),
				"userName":            userName,
			}
			err := goview.Render(w, http.StatusOK, "index", goviewModel)
			if err != nil {
				fmt.Fprintf(w, "Render index error: %v!", err)
			}
		})
		r.Get("/catalog", web.catalogHandler)
		r.Get("/catalog/publications/{id}", web.publicationHandler)
		r.Get("/catalog/syncpub", web.syncpubHandler)
		r.Get("/bookshelf/publications/{id}", web.publicationHandler)
		r.Get("/bookshelf/licenses/{id}", web.freshLicenseHandler)
		r.Get("/bookshelf/licenses/{id}/register", web.licenseRegisterHandler)
		r.Get("/bookshelf/licenses/{id}/renew", web.licenseRenewHandler)
		r.Get("/bookshelf/licenses/{id}/return", web.licenseReturnHandler)
		r.Get("/bookshelf/licenses/{id}/revoke", web.licenseRevokeHandler)
	})

	// Public signin/signout/signup
	r.Group(func(r chi.Router) {
		r.Get("/signin", web.signinCheck)
		r.Post("/signin", web.signin)
		r.Get("/signout", web.signout)
		r.Get("/signup", web.signupCheck)
		r.Post("/signup", web.signup)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(web.AuthMiddleware)
		// r.Get("/user/infos", userInfos)
		r.Get("/user/bookshelf", web.bookshelfHandler)
		r.Get("/catalog/publications/{id}/buy", web.createLicense)
		r.Get("/catalog/publications/{id}/loan", web.createLicense)
	})

	// 404 Handler - doit être défini au niveau du router principal, pas dans un groupe
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, filepath.Join(web.Config.RootDir, "static", "404.html"))
	})
}

func (web *Web) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated
		if !web.userIsAuthenticated(r) {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		// If authenticated, call the next handler
		next.ServeHTTP(w, r)
	})
}
