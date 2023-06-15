package stor

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init() {

	var _db, err = gorm.Open(sqlite.Open("pub.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = _db

	// Migrate the schema
	db.AutoMigrate(&Language{}, &Publisher{}, &Author{}, &Category{}, &Publication{})

}

func Step() {

	publication := Publication{
		Title:           "Sample Publication",
		UUID:            "123456",
		DatePublication: time.Now(),
		Description:     "Sample description",
		Language: []Language{
			{Code: "en"},
			{Code: "fr"},
		},
		Publisher: []Publisher{
			{Name: "Publisher A"},
			{Name: "Publisher B"},
		},
		Author: []Author{
			{Name: "Author A"},
			{Name: "Author B"},
		},
		Category: []Category{
			{Name: "Category A"},
			{Name: "Category B"},
		},
	}
	db.Create(&publication)

	// var publication Publication
	db.First(&publication, 1) // Assuming the publication with ID 1 exists

	// Access the publication's attributes
	fmt.Println(publication.Title)
	fmt.Println(publication.UUID)
	fmt.Println(publication.DatePublication)
	fmt.Println(publication.Description)
	fmt.Println(publication.Language)
	fmt.Println(publication.Publisher)
}

func Stop() {

	// Close the database connection
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to close database connection")
	}
	sqlDB.Close()

}
