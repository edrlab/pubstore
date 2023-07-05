package stor

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Stor struct {
	db *gorm.DB
}

func Init(dsn string) *Stor {

	if len(dsn) == 0 {
		dsn = "pub.db"
	}

	db, err := gorm.Open(GormDialector(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// db = db.Session(&gorm.Session{FullSaveAssociations: true})

	// Migrate the schema
	db.AutoMigrate(&Language{}, &Publisher{}, &Author{}, &Category{}, &Publication{}, &User{}, &Transaction{})

	// Check if the table is empty
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		// Insert initial record
		createdUser := &User{
			UUID:        uuid.New().String(),
			Name:        "admin",
			Email:       "admin@edrlab.org",
			Pass:        "admin",
			LcpHintMsg:  "Do not used it",
			LcpPassHash: "edrlab",
			SessionId:   uuid.New().String(),
		}
		result := db.Create(createdUser)
		if result.Error != nil {
			panic(fmt.Errorf("failed to insert initial record: %w", result.Error))
		}
	}

	return &Stor{db: db}

}

func (stor *Stor) Stop() {

	// Close the database connection
	sqlDB, err := stor.db.DB()
	if err != nil {
		panic("failed to close database connection")
	}
	sqlDB.Close()

}
