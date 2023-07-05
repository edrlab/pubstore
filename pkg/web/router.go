package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/edrlab/pubstore/pkg/config"
	"github.com/edrlab/pubstore/pkg/lcp"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Web struct {
	view *view.View
	stor *stor.Stor
}

func Init(s *stor.Stor, v *view.View) *Web {
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
		"userName":            "",
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

	user, err := web.stor.GetUserByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(password)) != nil {
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

func (web *Web) signup(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signin handler
	// This function will handle the "/signin" route

	if web.userIsAuthenticated(r) {
		http.Redirect(w, r, "/index", http.StatusFound)
	}

	signupGoview(w, false)
}

func signupGoview(w http.ResponseWriter, userCreationFailed bool) {

	err := goview.Render(w, http.StatusOK, "signup", goview.M{
		"pageTitle":           "pubstore - signin",
		"userIsAuthenticated": false,
		"userName":            "",
		"userCreationFailed":  userCreationFailed,
	})
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

func (web *Web) signupPostHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the signup handler
	// This function will handle the "/signup" route
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract form values
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	lcpPass := r.Form.Get("lcpPass")
	lcpHint := r.Form.Get("lcpHint")

	// Create a new User instance
	newUser := stor.User{
		UUID:        uuid.New().String(),
		Name:        name,
		Email:       email,
		Pass:        password,
		LcpPassHash: lcp.CreateLcpPassHash(lcpPass),
		LcpHintMsg:  lcpHint,
		SessionId:   "",
	}

	// Perform validation on newUser if required

	// Save newUser to the database using your storage function
	err = web.stor.CreateUser(&newUser)
	if err != nil {
		signupGoview(w, true)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/signin", http.StatusFound)
}

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
		printRights = config.PrintRights
	} else {
		if printRights < 0 {
			printRights = 0
		}
		if printRights >= 10000 {
			printRights = 10000
		}
	}
	if copyRights, err = strconv.Atoi(copyRightsString); err != nil {
		copyRights = config.CopyRights
	} else {
		if copyRights < 0 {
			copyRights = 0
		}
		if copyRights >= 10000 {
			copyRights = 10000
		}
	}

	storUser := web.getUserByCookie(r)
	userUUID := storUser.UUID
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

	licenceBytes, err = lcp.LicenceBuy(pubUUID, userUUID, userEmail, textHint, hexValue, printRights, copyRights)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	licenceId, publicationTitle, _, _, _, _, _, err := lcp.ParseLicenceLCPL(licenceBytes)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	err = web.stor.CreateTransactionWithUUID(pubUUID, userUUID, licenceId)
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

	printRightsString := r.URL.Query().Get("printRights")
	copyRightsString := r.URL.Query().Get("copyRights")
	startDateString := r.URL.Query().Get("startDate")
	endDateString := r.URL.Query().Get("endDate")

	// timeString := "2023-06-14T01:08:15+00:00"
	// layout := "2006-01-02T15:04:05Z07:00"
	// layout := "YYYY-MM-DDTHH:mm:ss.nnnZ" // ISO 8601 extended
	layout := time.RFC3339 // JS toISOString
	startDate, err := time.Parse(layout, startDateString)
	if err != nil {
		fmt.Println(err.Error())
		startDate = time.Now()
	}
	endDate, err := time.Parse(layout, endDateString)
	if err != nil {
		fmt.Println(err.Error())
		endDate = time.Now().AddDate(0, 0, 7)
	}

	if printRights, err = strconv.Atoi(printRightsString); err != nil {
		printRights = config.PrintRights
	} else {
		if printRights < 0 {
			printRights = 0
		}
		if printRights >= 10000 {
			printRights = 10000
		}
	}
	if copyRights, err = strconv.Atoi(copyRightsString); err != nil {
		copyRights = config.CopyRights
	} else {
		if copyRights < 0 {
			copyRights = 0
		}
		if copyRights >= 10000 {
			copyRights = 10000
		}
	}

	storUser := web.getUserByCookie(r)
	userUUID := storUser.UUID
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

	licenceBytes, err = lcp.LicenceLoan(pubUUID, userUUID, userEmail, textHint, hexValue, printRights, copyRights, startDate, endDate)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	licenceId, publicationTitle, _, _, _, _, _, err := lcp.ParseLicenceLCPL(licenceBytes)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	err = web.stor.CreateTransactionWithUUID(pubUUID, userUUID, licenceId)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}
}

func (web *Web) publicationFreshLicenceHandler(w http.ResponseWriter, r *http.Request) {

	pubUUID := chi.URLParam(r, "id")

	publication, _ := web.stor.GetPublicationByUUID(pubUUID)
	user := web.getUserByCookie(r)
	transaction, err := web.stor.GetTransactionByUserAndPublication(user.ID, publication.ID)
	if err != nil {
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

	licenceBytes, err = lcp.GenerateFreshLicenceFromLcpServer(transaction.LicenceId, user.Email, user.LcpHintMsg, user.LcpPassHash)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}

	_, publicationTitle, _, _, _, _, _, err = lcp.ParseLicenceLCPL(licenceBytes)
	if err != nil {
		message += err.Error()
		errorWasHappend = true
	}
}

func (web *Web) bookshelfHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the bookshelf handler
	// This function will handle the "/user/bookshelf" route

	user := web.getUserByCookie(r)
	if user == nil {
		fmt.Fprintf(w, "bookshelf error")
		return
	}
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
		"userName":            user.Name,
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
		pageSizeInt = config.NumberOfPublicationsPerPage
	}

	var facet string = ""
	var value string = ""
	if len(queryStr) > 0 {
		facet = "search"
		value = queryStr
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

	facetsView := web.view.GetCatalogFacetsView()
	pubsView, count := web.view.GetCatalogPublicationsView(facet, value, pageInt, pageSizeInt)
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

func (web *Web) publicationHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for the publication handler
	// This function will handle the "/catalog/publication/{id}" route

	pubUUID := chi.URLParam(r, "id")
	errLcp := r.URL.Query().Get("err")
	userStor := web.getUserByCookie(r)
	licenseOK := false

	if publicationStor, err := web.stor.GetPublicationByUUID(pubUUID); err != nil {
		http.ServeFile(w, r, "static/404.html")
		w.WriteHeader(http.StatusNotFound)
	} else {
		viewTransaction := view.TransactionView{}
		userName := ""
		if userStor != nil {
			userName = userStor.Name
			transaction, err := web.stor.GetTransactionByUserAndPublication(userStor.ID, publicationStor.ID)
			if err == nil {
				viewTransaction := *web.view.GetTransactionViewFromTransactionStor(transaction)
				if viewTransaction.LicenseStatusCode == "ready" || viewTransaction.LicenseStatusCode == "active" {
					licenseOK = true
				}
			}
		}

		publicationView := web.view.GetPublicationViewFromPublicationStor(publicationStor)
		goviewModel := goview.M{
			"pageTitle":             fmt.Sprintf("pubstore - %s", publicationView.Title),
			"userIsAuthenticated":   web.userIsAuthenticated(r),
			"userName":              userName,
			"errLcp":                errLcp,
			"licenseFoundAndActive": licenseOK,
			"title":                 publicationView.Title,
			"uuid":                  publicationView.UUID,
			"datePublication":       publicationView.DatePublication,
			"description":           publicationView.Description,
			"coverUrl":              publicationView.CoverUrl,
			"authors":               publicationView.Author,
			"publishers":            publicationView.Publisher,
			"languages":             publicationView.Language,
			"categories":            publicationView.Category,
			"licenseFound":          bool(viewTransaction != view.TransactionView{}),
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

func (web *Web) Rooter(r chi.Router) {

	// Serve static files from the "static" directory
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
		r.Get("/signup", web.signup)
		r.Post("/signup", web.signupPostHandler)
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
}
