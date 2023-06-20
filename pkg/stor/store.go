package stor

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Stor struct {
	db *gorm.DB
}

func Init(sqliteDsn string) *Stor {

	if len(sqliteDsn) == 0 {
		sqliteDsn = "pub.db"
	}

	db, err := gorm.Open(sqlite.Open(sqliteDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// db = db.Session(&gorm.Session{FullSaveAssociations: true})

	// Migrate the schema
	db.AutoMigrate(&Language{}, &Publisher{}, &Author{}, &Category{}, &Publication{}, &User{}, &Transaction{})

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
