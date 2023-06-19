package web

import (
	"fmt"
	"net/http"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Web struct {
	view *view.View
	stor *stor.Stor
}

func Init(s *stor.Stor) *Web {
	v := view.Init(s)
	return &Web{stor: s, view: v}
}

func signin(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signin handler
	// This function will handle the "/signin" route
}

func signout(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signout handler
	// This function will handle the "/signout" route
}

func signup(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signup handler
	// This function will handle the "/signup" route
}

func userInfos(w http.ResponseWriter, r *http.Request) {
	// Implementation for the userInfos handler
	// This function will handle the "/user/infos" and "/user/bookshelf" routes
}

func publicationBuyAction(w http.ResponseWriter, r *http.Request) {
	// Implementation for the publicationBuyAction handler
	// This function will handle the "/catalog/publication/{id}/buy" route
}

func publicationLoanAction(w http.ResponseWriter, r *http.Request) {
	// Implementation for the publicationLoanAction handler
	// This function will handle the "/catalog/publication/{id}/loan" route
}

func (web *Web) catalogHangler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the catalog handler
	// This function will handle the "/catalog" route

	author := r.URL.Query().Get("author")
	language := r.URL.Query().Get("language")
	publisher := r.URL.Query().Get("publisher")
	category := r.URL.Query().Get("category")

	var facet string = ""
	var value string = ""
	if len(author) > 0 {
		facet = "author"
		value = author
	} else if len(publisher) > 0 {
		facet = "publisher"
		value = publisher
	} else if len(language) > 0 {
		facet = "language"
		value = language
	} else if len(category) > 0 {
		facet = "category"
		value = category
	}

	facetsView := web.view.GetFacetsView()
	pubsView := web.view.GetPublicationsView(facet, value)
	catalogView := view.GetCatalogView(pubsView, facetsView)

	goviewModel := goview.M{
		"authors":      (*catalogView).Authors,
		"publishers":   (*catalogView).Publishers,
		"languages":    (*catalogView).Languages,
		"categories":   (*catalogView).Categories,
		"publications": (*catalogView).Publications,
	}

	err := goview.Render(w, http.StatusOK, "catalog", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

func publication(w http.ResponseWriter, r *http.Request) {
	// Implementation for the publication handler
	// This function will handle the "/catalog/publication/{id}" route
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated
		if !IsUserAuthenticated(r) {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		// If authenticated, call the next handler
		next.ServeHTTP(w, r)
	})
}

func IsUserAuthenticated(r *http.Request) bool {
	// Implement your authentication logic here
	// Check if the user is authenticated based on the session, token, or any other authentication mechanism

	// For the purpose of this example, we assume the user is authenticated if there is a session cookie
	_, err := r.Cookie("session")
	return err == nil
}

func (web *Web) Rooter() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files from the "static" directory
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/index", http.StatusFound)
		})
		r.Get("/index", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/index.html")
		})
		r.Get("/catalog", web.catalogHangler)
		r.Get("/catalog/publication/{id}", publication)
	})

	// Public signin/signout/signup
	r.Group(func(r chi.Router) {
		r.Get("/signin", signin)
		r.Get("/signout", signout)
		r.Get("/signup", signup)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/user/infos", userInfos)
		r.Get("/user/bookshelf", userInfos)
		r.Get("/catalog/publication/{id}/buy", publicationBuyAction)
		r.Post("/catalog/publication/{id}/loan", publicationLoanAction)
	})

	return r
}
