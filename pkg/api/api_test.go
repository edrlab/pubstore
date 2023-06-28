package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-playground/validator/v10"
)

var api *Api

func TestMain(m *testing.M) {

	validate = validator.New()

	s := stor.Init("file::memory:?cache=shared")

	api = &Api{stor: s}

	// Run the tests
	exitCode := m.Run()

	s.Stop()

	fmt.Println("ExitCode", exitCode)
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func TestSuite(t *testing.T) {
	t.Run("TestPublicationHandler", TestPublicationHandler)
	t.Run("TestUserHandler", TestUserHandler)
}
