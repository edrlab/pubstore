package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/lcp"
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

func (web *Web) publicationFreshLicenceHandler(w http.ResponseWriter, r *http.Request) {

	pubUUID := chi.URLParam(r, "id")

	publication, _ := web.Store.GetPublication(pubUUID)
	user := web.getUserByCookie(r)
	transaction, err := web.Store.GetTransactionByUserAndPublication(user.ID, publication.ID)
	if err != nil {
		http.Redirect(w, r, "/catalog/publication/"+pubUUID, http.StatusFound)
		return
	}

	var licenceBytes []byte
	var publicationTitle string

	licenceBytes, err = lcp.GetFreshLicense(web.Config.LCPServer, transaction)

	if err == nil {
		_, publicationTitle, _, _, _, _, _, err = lcp.ParseLicense(licenceBytes)
	}

	if err != nil {
		message := url.QueryEscape("Failed to generate a fresh license: " + err.Error())
		http.Redirect(w, r, fmt.Sprintf("/catalog/publication/%s?err=%s", pubUUID, message), http.StatusFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+publicationTitle+".lcpl")
	w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(licenceBytes)))

	io.Copy(w, bytes.NewReader(licenceBytes))
}

func (web *Web) bookshelfHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the bookshelf handler
	// This function will handle the "/user/bookshelf" route

	user := web.getUserByCookie(r)
	if user == nil {
		fmt.Fprintf(w, "bookshelf error")
		return
	}
	transactions, err := web.Store.FindTransactionsByUser(user.ID)
	if err != nil {
		fmt.Fprintf(w, "bookshelf error")
		return
	}

	var transactionsView []*view.TransactionView = make([]*view.TransactionView, len(*transactions))
	for i, transactionStor := range *transactions {
		transactionsView[i] = web.View.GetTransactionViewFromTransactionStor(&transactionStor)
	}

	goviewModel := goview.M{
		"pageTitle":           "pubstore - bookshelf",
		"userIsAuthenticated": true,
		"userName":            user.Name,
		"transactions":        transactionsView,
	}

	err = goview.Render(w, http.StatusOK, "bookshelf", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

func (web *Web) catalogHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the catalog handler
	// This function will handle the "/catalog" route

	format := r.URL.Query().Get("format")
	author := r.URL.Query().Get("author")
	language := r.URL.Query().Get("language")
	publisher := r.URL.Query().Get("publisher")
	category := r.URL.Query().Get("category")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	q := r.URL.Query().Get("q")
	query := r.URL.Query().Get("query")
	queryStr := query
	if len(query) == 0 {
		queryStr = q
	}

	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 || pageInt > 1000 {
		pageInt = 1
	}
	pageSizeInt, _ := strconv.Atoi(pageSize)
	if pageSizeInt < 1 || pageSizeInt > 1000 {
		pageSizeInt = web.Config.PageSize
	}

	var facet string = ""
	var value string = ""
	if len(queryStr) > 0 {
		facet = "search"
		value = queryStr
	} else if len(format) > 0 {
		facet = "format"
		value = format
	} else if len(author) > 0 {
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

	facetsView := web.View.GetCatalogFacetsView()
	pubsView, count := web.View.GetCatalogPublicationsView(facet, value, pageInt, pageSizeInt)
	catalogView := view.GetCatalogView(pubsView, facetsView)

	var pageRange []string = make([]string, pageInt)
	for i := 0; i < pageInt; i++ {
		pageRange[i] = fmt.Sprintf("%d", i+1)
	}
	userStor := web.getUserByCookie(r)
	userName := ""
	if userStor != nil {
		userName = userStor.Name
	}

	goviewModel := goview.M{
		"pageTitle":           "pubstore - catalog",
		"userIsAuthenticated": web.userIsAuthenticated(r),
		"userName":            userName,
		"currentFacetType":    facet,
		"currentFacetValue":   value,
		"currentPageSize":     fmt.Sprintf("%d", pageSizeInt),
		"currentPage":         fmt.Sprintf("%d", pageInt),
		"pageRange":           pageRange,
		"publicationCount":    fmt.Sprintf("%d", count),
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

// publicationHandler handles a request for a publication page
func (web *Web) publicationHandler(w http.ResponseWriter, r *http.Request) {

	pubUUID := chi.URLParam(r, "id")
	errLcp := r.URL.Query().Get("err")
	userStor := web.getUserByCookie(r)
	licenseOK := false

	if publicationStor, err := web.Store.GetPublication(pubUUID); err != nil {
		http.ServeFile(w, r, "static/404.html")
		w.WriteHeader(http.StatusNotFound)
	} else {
		viewTransaction := view.TransactionView{}
		userName := ""
		if userStor != nil {
			userName = userStor.Name
			transaction, err := web.Store.GetTransactionByUserAndPublication(userStor.ID, publicationStor.ID)
			if err == nil {
				viewTransaction = *web.View.GetTransactionViewFromTransactionStor(transaction)
				if viewTransaction.LicenseStatusCode == "ready" || viewTransaction.LicenseStatusCode == "active" {
					licenseOK = true
				}
			}
		}

		publicationView := web.View.GetPublicationViewFromPublicationStor(publicationStor)
		goviewModel := goview.M{
			"pageTitle":             fmt.Sprintf("pubstore - %s", publicationView.Title),
			"host":                  strings.Split(web.Config.PublicBaseUrl, "://")[1],
			"userIsAuthenticated":   web.userIsAuthenticated(r),
			"userName":              userName,
			"errLcp":                errLcp,
			"licenseFoundAndActive": licenseOK,
			"title":                 publicationView.Title,
			"uuid":                  publicationView.UUID,
			"format":                publicationView.Format,
			"datePublished":         publicationView.DatePublished,
			"description":           publicationView.Description,
			"coverUrl":              publicationView.CoverUrl,
			"authors":               publicationView.Author,
			"publishers":            publicationView.Publisher,
			"languages":             publicationView.Language,
			"categories":            publicationView.Category,
			"licenseFound":          bool(viewTransaction.PublicationUUID != ""),
			"transaction":           viewTransaction,
		}
		err = goview.Render(w, http.StatusOK, "publication", goviewModel)
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

func (web *Web) Router(r chi.Router) {

	// Serve static files from the "static" directory. The '*' means that sub-routes are served from sub-directories
	filesDir := http.Dir(filepath.Join(web.Config.RootDir, "static"))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))

	// Serve resources from a configurable directory (used for cover images)
	//r.Handle("/resources/*", http.StripPrefix("/resources/", http.FileServer(http.Dir(web.Config.Resources))))
	//fmt.Println("Resources fetched from ", web.Config.Resources)

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
		r.Get("/catalog/publication/{id}", web.publicationHandler)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/404.html")
			w.WriteHeader(http.StatusNotFound)
		})
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
		r.Get("/catalog/publication/{id}/buy", web.createLicense)
		r.Get("/catalog/publication/{id}/loan", web.createLicense)
		r.Get("/catalog/publication/{id}/license", web.publicationFreshLicenceHandler)
	})
}
