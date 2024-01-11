package stor

import (
	"fmt"
	"os"
	"testing"
)

// Global store
var store Store

func TestMain(m *testing.M) {

	var err error
	store, err = Init("sqlite3://file::memory:?cache=shared")
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Run the tests
	exitCode := m.Run()

	fmt.Println("ExitCode", exitCode)
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}
