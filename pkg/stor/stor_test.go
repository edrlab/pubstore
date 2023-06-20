package stor

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	// Set up the database connection
	var db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Run migrations
	err = db.AutoMigrate(&Publication{}, &Language{}, &Publisher{}, &Author{}, &Category{}, &User{}, &Transaction{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	stor.db = db

	// Run the tests
	exitCode := m.Run()

	fmt.Println("ExitCode", exitCode)
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func TestSuite(t *testing.T) {
	t.Run("PublicationCRUD", TestPublicationCRUD)
	t.Run("TestCreate2PublicationsWithSameCategory", TestCreate2PublicationsWithSameCategory)
	t.Run("GetPublicationByCategory", TestGetPublicationByCategory)
	t.Run("GetPublicationByLanguage", TestGetPublicationByLanguage)
	t.Run("GetPublicationByPublisher", TestGetPublicationByPublisher)
	t.Run("GetPublicationByAuthor", TestGetPublicationByAuthor)
	t.Run("GetPublicationByUUID", TestGetPublicationByUUID)
	t.Run("UserCRUD", TestUserCRUD)
	t.Run("TransactionCRUD", TestTransactionCRUD)
}
