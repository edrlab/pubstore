package web

import (
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/foolin/goview"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// getUserByCookie retrieves a user via a cookie
func (web *Web) getUserByCookie(r *http.Request) *stor.User {

	if cookie, err := r.Cookie("session"); err == nil {
		sessionId := cookie.Value

		if user, err := web.Store.GetUserBySession(sessionId); err == nil {
			return user
		}
	}
	return nil
}

// userIsAuthenticated checks if the user cookie is valid
func (web *Web) userIsAuthenticated(r *http.Request) bool {

	if err := web.getUserByCookie(r); err != nil {
		return true
	}
	return false
}

// signin checks the email and password passed as parameters, initiates a session and redirects to the signin view
func (web *Web) signin(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := web.Store.GetUserByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.HPassword), []byte(password)) != nil {
		signinGoview(w, true)
		return
	}

	sessionId := uuid.New().String()
	user.SessionId = sessionId

	if err := web.Store.UpdateUser(user); err != nil {
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

// signinCheck checks if a user is authenticated and, if not, redirects to the signin view
func (web *Web) signinCheck(w http.ResponseWriter, r *http.Request) {

	if web.userIsAuthenticated(r) {
		http.Redirect(w, r, "/index", http.StatusFound)
	}

	signinGoview(w, false)
}

// signinGoview displays the signin view
func signinGoview(w http.ResponseWriter, userNotFound bool) {

	err := goview.Render(w, http.StatusOK, "signin", goview.M{
		"pageTitle": "pubstore - signin",
		//"userIsAuthenticated": false,
		//"userName":            "",
		"userNotFound": userNotFound,
	})
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

// signout resets the used cookie and redirects to the home
func (web *Web) signout(w http.ResponseWriter, r *http.Request) {

	cookie := &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}

	http.SetCookie(w, cookie) // Send the cookie in the response to remove it

	http.Redirect(w, r, "/index", http.StatusFound)
}

// signup checks the form parameters, creates a user and redirects to the signin view
func (web *Web) signup(w http.ResponseWriter, r *http.Request) {

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// extract form values
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	lcpPass := r.Form.Get("lcpPass")
	lcpHint := r.Form.Get("lcpHint")

	// create a new User instance
	newUser := stor.User{
		UUID:       uuid.New().String(),
		Name:       name,
		Email:      email,
		Password:   password,
		TextHint:   lcpHint,
		Passphrase: lcpPass,
	}

	// perform validation.
	// the minimum length of the passphrase is 3 characters
	if utf8.RuneCountInString(newUser.Passphrase) < 3 {
		signupGoview(w, true)
		return
	}

	// save newUser to the database using your storage function
	err = web.Store.CreateUser(&newUser)
	if err != nil {
		signupGoview(w, true)
		return
	}
	http.Redirect(w, r, "/signin", http.StatusFound)
}

// signupCheck checks if a user is authenticated and, if not, redirects to the signup view
func (web *Web) signupCheck(w http.ResponseWriter, r *http.Request) {

	if web.userIsAuthenticated(r) {
		http.Redirect(w, r, "/index", http.StatusFound)
	}

	signupGoview(w, false)
}

// signupGoview displays the signup view
func signupGoview(w http.ResponseWriter, userCreationFailed bool) {

	err := goview.Render(w, http.StatusOK, "signup", goview.M{
		"pageTitle": "pubstore - signup",
		//"userIsAuthenticated": false,
		//"userName":            "",
		"userCreationFailed": userCreationFailed,
	})
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}
