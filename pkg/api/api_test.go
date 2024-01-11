package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
)

var testapi Api

func TestMain(m *testing.M) {

	config := conf.Config{OAuthSeed: "EDRLAB_Rocks"}

	store, err := stor.Init("sqlite3://file::memory:?cache=shared")
	if err != nil {
		panic("Database setup failed.")
	}

	testapi = Init(&config, &store)

	// Run the tests
	exitCode := m.Run()

	fmt.Println("ExitCode", exitCode)
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}
