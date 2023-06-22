package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Web struct {
	view *view.View
	stor *stor.Stor
}

func Init(s *stor.Stor) *Web {
	v := view.Init(s)
	return &Web{stor: s, view: v}
}

func (web *Web) userIsAuthenticated(r *http.Request) bool {

	if cookie, err := r.Cookie("session"); err == nil {
		sessionId := cookie.Value

		fmt.Println(sessionId)
		if _, err = web.stor.GetUserBySessionId(sessionId); err == nil {
			// redirect to /index the user is logged
			return true
		}
	}
	return false
}

func (web *Web) signin(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signin handler
	// This function will handle the "/signin" route

	if web.userIsAuthenticated(r) {
		http.Redirect(w, r, "/index", http.StatusFound)
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	fmt.Println(email)
	fmt.Println(password)

	signinGoview(w, false)
}

func (web *Web) signout(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signin handler
	// This function will handle the "/signin" route

	cookie := &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}

	http.SetCookie(w, cookie) // Send the cookie in the response to remove it

	http.Redirect(w, r, "/index", http.StatusFound)
}

func signinGoview(w http.ResponseWriter, userNotFound bool) {

	err := goview.Render(w, http.StatusOK, "signin", goview.M{
		"pageTitle":           "pubstore - signin",
		"userIsAuthenticated": false,
		"userNotFound":        userNotFound,
	})
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

func (web *Web) signinPost(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signin handler
	// This function will handle the "/signin" route

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	fmt.Println(email)
	fmt.Println(password)

	var err error
	user, err := web.stor.GetUserByEmailAndPass(email, password)
	if err != nil {
		signinGoview(w, true)
		return
	}

	sessionId := uuid.New().String()
	user.SessionId = sessionId

	if err := web.stor.UpdateUser(user); err != nil {
		signinGoview(w, true)
		return
	}

	cookie := &http.Cookie{
		Name:    "session",
		Value:   sessionId,
		Expires: time.Now().Add(24 * time.Hour), // Set cookie expiration time
		Path:    "/",
	}

	http.SetCookie(w, cookie) // Send the cookie in the response

	http.Redirect(w, r, "/index", http.StatusFound)

}

// func signup(w http.ResponseWriter, r *http.Request) {
// 	// Implementation for the signup handler
// 	// This function will handle the "/signup" route
// }

// func userInfos(w http.ResponseWriter, r *http.Request) {
// 	// Implementation for the userInfos handler
// 	// This function will handle the "/user/infos" and "/user/bookshelf" routes
// }

// func publicationBuyAction(w http.ResponseWriter, r *http.Request) {
// 	// Implementation for the publicationBuyAction handler
// 	// This function will handle the "/catalog/publication/{id}/buy" route
// }

// func publicationLoanAction(w http.ResponseWriter, r *http.Request) {
// 	// Implementation for the publicationLoanAction handler
// 	// This function will handle the "/catalog/publication/{id}/loan" route
// }

func (web *Web) bookshelfHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "bookshelf")
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
		"pageTitle":           "pubstore - catalog",
		"userIsAuthenticated": web.userIsAuthenticated(r),
		"authors":             (*catalogView).Authors,
		"publishers":          (*catalogView).Publishers,
		"languages":           (*catalogView).Languages,
		"categories":          (*catalogView).Categories,
		"publications":        (*catalogView).Publications,
	}

	err := goview.Render(w, http.StatusOK, "catalog", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

func (web *Web) publicationHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the publication handler
	// This function will handle the "/catalog/publication/{id}" route

	pubUUID := chi.URLParam(r, "id")

	if publication, err := web.view.GetPublicationView(pubUUID); err != nil {
		http.ServeFile(w, r, "static/404.html")
		w.WriteHeader(http.StatusNotFound)
	} else {
		goviewModel := goview.M{
			"pageTitle":           fmt.Sprintf("pubstore - %s", publication.Title),
			"userIsAuthenticated": web.userIsAuthenticated(r),
			"title":               publication.Title,
			"uuid":                publication.UUID,
			"datePublication":     publication.DatePublication,
			"description":         publication.Description,
			"coverUrl":            publication.CoverUrl,
			"authors":             publication.Author,
			"publishers":          publication.Publisher,
			"languages":           publication.Language,
			"categories":          publication.Category,
		}
		err := goview.Render(w, http.StatusOK, "publication", goviewModel)
		if err != nil {
			fmt.Fprintf(w, "Render index error: %v!", err)
		}
	}
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
			goviewModel := goview.M{
				"pageTitle":           "pubstore",
				"userIsAuthenticated": web.userIsAuthenticated(r),
			}
			err := goview.Render(w, http.StatusOK, "index", goviewModel)
			if err != nil {
				fmt.Fprintf(w, "Render index error: %v!", err)
			}
		})
		r.Get("/catalog", web.catalogHangler)
		r.Get("/catalog/publication/{id}", web.publicationHandler)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/404.html")
			w.WriteHeader(http.StatusNotFound)
		})
	})

	// Public signin/signout/signup
	r.Group(func(r chi.Router) {
		r.Get("/signin", web.signin)
		r.Post("/signin", web.signinPost)
		r.Get("/signout", web.signout)
		// r.Get("/signup", signup)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(web.AuthMiddleware)
		// r.Get("/user/infos", userInfos)
		r.Get("/user/bookshelf", web.bookshelfHandler)
		// r.Get("/catalog/publication/{id}/buy", publicationBuyAction)
		// r.Post("/catalog/publication/{id}/loan", publicationLoanAction)
	})

	return r
}
