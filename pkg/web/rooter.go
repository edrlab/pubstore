package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/edrlab/pubstore/pkg/lcp"
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

func (web *Web) getUserByCookie(r *http.Request) *stor.User {

	if cookie, err := r.Cookie("session"); err == nil {
		sessionId := cookie.Value

		if user, err := web.stor.GetUserBySessionId(sessionId); err == nil {
			// redirect to /index the user is logged
			return user
		}
	}
	return nil
}

func (web *Web) userIsAuthenticated(r *http.Request) bool {

	if err := web.getUserByCookie(r); err != nil {
		return true
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

func (web *Web) publicationBuyHandler(w http.ResponseWriter, r *http.Request) {
	// This function will handle the "/catalog/publication/{id}/buy" route

	var printRights int
	var copyRights int
	var err error

	pubUUID := chi.URLParam(r, "id")

	printRightsString := r.URL.Query().Get("printRigths")
	copyRightsString := r.URL.Query().Get("copyRights")

	if printRights, err = strconv.Atoi(printRightsString); err != nil {
		printRights = 10
	} else {
		if printRights < 0 {
			printRights = 0
		}
		if printRights >= 10000 {
			printRights = 10000
		}
	}
	if copyRights, err = strconv.Atoi(copyRightsString); err != nil {
		copyRights = 10
	} else {
		if copyRights < 0 {
			copyRights = 0
		}
		if copyRights >= 10000 {
			copyRights = 10000
		}
	}

	storUser := web.getUserByCookie(r)
	userId := storUser.UUID
	userEmail := storUser.Email
	textHint := storUser.LcpHintMsg
	hexValue := storUser.LcpPassHash

	message := "something went wrong with buy function : "

	var licenceBytes []byte
	var errorWasHappend bool = false
	var publicationTitle string

	defer func() {
		if errorWasHappend {

			http.Redirect(w, r, fmt.Sprintf("/catalog/publication/%s?err=%s", pubUUID, url.QueryEscape(message)), http.StatusFound)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+publicationTitle+".lcpl")
		w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(licenceBytes)))

		io.Copy(w, bytes.NewReader(licenceBytes))
	}()

	licenceBytes, err = lcp.LicenceBuy(pubUUID, userId, userEmail, textHint, hexValue, printRights, copyRights)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	licenceId, publicationTitle, err := lcp.ParseLicenceLCPL(licenceBytes)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	err = web.stor.CreateTransactionWithUUID(pubUUID, userId, licenceId)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	if errorWasHappend {

		http.Redirect(w, r, fmt.Sprintf("/catalog/publication/%s?err=%s", pubUUID, url.QueryEscape(message)), http.StatusFound)
		return
	}
}

func (web *Web) publicationLoanHandler(w http.ResponseWriter, r *http.Request) {
	// This function will handle the "/catalog/publication/{id}/buy" route

	var printRights int
	var copyRights int
	var err error

	pubUUID := chi.URLParam(r, "id")

	printRightsString := r.URL.Query().Get("printRigths")
	copyRightsString := r.URL.Query().Get("copyRights")
	startDateString := r.URL.Query().Get("startDate") + ":00+00:00"
	endDateString := r.URL.Query().Get("endDate") + ":00+00:00"

	// timeString := "2023-06-14T01:08:15+00:00"
	layout := "2006-01-02T15:04:05Z07:00"
	startDate, err := time.Parse(layout, startDateString)
	if err != nil {
		startDate = time.Now()
	}
	endDate, err := time.Parse(layout, endDateString)
	if err != nil {
		endDate = time.Now()
	}

	if printRights, err = strconv.Atoi(printRightsString); err != nil {
		printRights = 10
	} else {
		if printRights < 0 {
			printRights = 0
		}
		if printRights >= 10000 {
			printRights = 10000
		}
	}
	if copyRights, err = strconv.Atoi(copyRightsString); err != nil {
		copyRights = 10
	} else {
		if copyRights < 0 {
			copyRights = 0
		}
		if copyRights >= 10000 {
			copyRights = 10000
		}
	}

	storUser := web.getUserByCookie(r)
	userId := storUser.UUID
	userEmail := storUser.Email
	textHint := storUser.LcpHintMsg
	hexValue := storUser.LcpPassHash

	message := "something went wrong with loan function : "

	var licenceBytes []byte
	var errorWasHappend bool = false
	var publicationTitle string

	defer func() {
		if errorWasHappend {

			http.Redirect(w, r, fmt.Sprintf("/catalog/publication/%s?err=%s", pubUUID, url.QueryEscape(message)), http.StatusFound)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+publicationTitle+".lcpl")
		w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(licenceBytes)))

		io.Copy(w, bytes.NewReader(licenceBytes))
	}()

	licenceBytes, err = lcp.LicenceLoan(pubUUID, userId, userEmail, textHint, hexValue, printRights, copyRights, startDate, endDate)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	licenceId, publicationTitle, err := lcp.ParseLicenceLCPL(licenceBytes)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	err = web.stor.CreateTransactionWithUUID(pubUUID, userId, licenceId)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}
}

func (web *Web) getLicenceIdTransaction(publication *stor.Publication, user *stor.User) string {

	if publication == nil || user == nil {
		return ""
	}

	userID := user.ID
	pubID := publication.ID
	transaction, _ := web.stor.GetTransactionByUserAndPublication(userID, pubID)

	if transaction == nil {
		return ""
	}

	licenceUUID := transaction.LicenceId
	return licenceUUID
}

func (web *Web) publicationFreshLicenceHandler(w http.ResponseWriter, r *http.Request) {

	pubUUID := chi.URLParam(r, "id")

	publication, _ := web.stor.GetPublicationByUUID(pubUUID)
	user := web.getUserByCookie(r)
	licenceID := web.getLicenceIdTransaction(publication, user)
	if len(licenceID) == 0 {
		http.Redirect(w, r, "/catalog/publication/"+pubUUID, http.StatusFound)
		return
	}

	message := "something went wrong to generate a fresh license : "
	errorWasHappend := false

	var licenceBytes []byte
	var publicationTitle string

	defer func() {
		if errorWasHappend {

			http.Redirect(w, r, fmt.Sprintf("/catalog/publication/%s?err=%s", pubUUID, url.QueryEscape(message)), http.StatusFound)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+publicationTitle+".lcpl")
		w.Header().Set("Content-Type", "application/vnd.readium.lcp.license.v1.0+json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(licenceBytes)))

		io.Copy(w, bytes.NewReader(licenceBytes))
	}()

	licenceBytes, err := lcp.GenerateFreshLicenceFromLcpServer(licenceID, user.Email, user.LcpHintMsg, user.LcpPassHash)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	_, publicationTitle, err = lcp.ParseLicenceLCPL(licenceBytes)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}
}

func (web *Web) bookshelfHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the bookshelf handler
	// This function will handle the "/user/bookshelf" route

	user := web.getUserByCookie(r)
	transactions, err := web.stor.GetTransactionsByUserID(user.ID)
	if err != nil {

		fmt.Fprintf(w, "bookshelf error")
		return
	}

	var transactionsView []*view.TransactionView = make([]*view.TransactionView, len(*transactions))
	for i, transactionStor := range *transactions {
		transactionsView[i] = web.view.GetTransactionViewFromTransactionStor(&transactionStor)
	}

	goviewModel := goview.M{
		"pageTitle":           "pubstore - bookshelf",
		"userIsAuthenticated": true,
		"transactions":        transactionsView,
	}

	err = goview.Render(w, http.StatusOK, "bookshelf", goviewModel)
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
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

	facetsView := web.view.GetCatalogFacetsView()
	pubsView := web.view.GetCatalogPublicationsView(facet, value)
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
	errLcp := r.URL.Query().Get("err")
	user := web.getUserByCookie(r)

	if publicationStor, err := web.stor.GetPublicationByUUID(pubUUID); err != nil {
		http.ServeFile(w, r, "static/404.html")
		w.WriteHeader(http.StatusNotFound)
	} else {
		licenceId := web.getLicenceIdTransaction(publicationStor, user)
		fmt.Println(licenceId)
		licenceIdFound := len(licenceId) > 0

		publicationView := web.view.GetPublicationViewFromPublicationStor(publicationStor)
		goviewModel := goview.M{
			"pageTitle":           fmt.Sprintf("pubstore - %s", publicationView.Title),
			"userIsAuthenticated": web.userIsAuthenticated(r),
			"errLcp":              errLcp,
			"licenseIdFound":      licenceIdFound,
			"title":               publicationView.Title,
			"uuid":                publicationView.UUID,
			"datePublication":     publicationView.DatePublication,
			"description":         publicationView.Description,
			"coverUrl":            publicationView.CoverUrl,
			"authors":             publicationView.Author,
			"publishers":          publicationView.Publisher,
			"languages":           publicationView.Language,
			"categories":          publicationView.Category,
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
		r.Get("/catalog/publication/{id}/buy", web.publicationBuyHandler)
		r.Get("/catalog/publication/{id}/loan", web.publicationLoanHandler)
		r.Get("/catalog/publication/{id}/license", web.publicationFreshLicenceHandler)
	})

	return r
}
