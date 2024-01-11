package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/edrlab/pubstore/pkg/view"
	"github.com/go-chi/chi/v5"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

var web Web

func TestMain(m *testing.M) {

	config := conf.Config{}

	store, err := stor.Init("sqlite3://file::memory:?cache=shared")
	if err != nil {
		panic("Database setup failed.")
	}

	view := view.Init(&config, &store)
	web = Init(&config, &store, &view)

	// Run the tests
	exitCode := m.Run()

	fmt.Println("ExitCode", exitCode)
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func TestSign(t *testing.T) {

	// Initialize the router
	r := chi.NewRouter()
	r.Group(web.Router)

	// create a new user, directly in the store
	testUser := &stor.User{
		UUID:       gofakeit.UUID(),
		Name:       "Pierre ler",
		Email:      gofakeit.Email(),
		Password:   "password",
		TextHint:   "hint",
		Passphrase: "passphrase",
	}

	err := web.Store.CreateUser(testUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, testUser.ID)

	// validate the user
	err = testUser.Validate()
	assert.NoError(t, err)

	// signin requires an email and a password
	form := url.Values{}
	form.Add("email", testUser.Email)
	form.Add("password", testUser.Password)

	// sign in the user
	req := httptest.NewRequest("POST", "/signin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusFound, recorder.Code)

	// sign out the user
	req = httptest.NewRequest("GET", "/signout", nil)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusFound, recorder.Code)

}
