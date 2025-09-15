package web

import (
	"crypto/sha256"
	"encoding/hex"
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
	"golang.org/x/crypto/bcrypt"

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

	// generate a hash of the user password
	userPassword := "user-password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// generate a hash of the lcp passphrase
	var hashedPassphrase string
	hash := sha256.Sum256([]byte("lcpPassphrase"))
	hashedPassphrase = hex.EncodeToString(hash[:])

	// create a new user, directly in the store
	testUser := &stor.User{
		UUID:        gofakeit.UUID(),
		Name:        "Pierre ler",
		Email:       gofakeit.Email(),
		HPassword:   string(hashedPassword),
		HPassphrase: hashedPassphrase,
		TextHint:    "hint",
	}

	err = web.Store.CreateUser(testUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, testUser.ID)

	// validate the user
	err = testUser.Validate()
	assert.NoError(t, err)

	// signin requires an email and a password
	form := url.Values{}
	form.Add("email", testUser.Email)
	form.Add("password", userPassword)

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
