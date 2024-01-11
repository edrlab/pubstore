// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Use of this source code is governed by a BSD-style license
// specified in the Github project LICENSE file.

// The stor package manages the database storage of pubstore entities.
package stor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Store defines a generic store with a gorm db
type Store struct {
	db *gorm.DB
}

// Init initializes the database
func Init(dsn string) (Store, error) {
	str := Store{}

	if len(dsn) == 0 {
		return str, errors.New("database source name is missing")
	}
	dialect, cnx := dbFromURI(dsn)
	if dialect == "error" {
		return str, fmt.Errorf("incorrect database source name: %q", dsn)
	}

	// the use of time.Time fields for mysql requires parseTime
	if dialect == "mysql" && !strings.Contains(cnx, "parseTime") {
		return str, fmt.Errorf("incomplete mysql database source name, parseTime required: %q", dsn)
	}
	// Any constraint for other databases?

	// database logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Warn, // Log level (Silent, Error, Warn, Info)
			//LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true, // Ignore ErrRecordNotFound error for logger
			Colorful:                  true, // Enable color
		},
	)

	db, err := gorm.Open(GormDialector(cnx), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Printf("Failed connecting to the database: %v", err)
		return str, err
	}

	err = performDialectSpecific(db, dialect)
	if err != nil {
		log.Printf("Failed performing dialect specific database init: %v", err)
		return str, err
	}

	// db = db.Session(&gorm.Session{FullSaveAssociations: true})

	err = db.AutoMigrate(&Language{}, &Publisher{}, &Author{}, &Category{}, &Publication{}, &User{}, &Transaction{})
	if err != nil {
		log.Printf("Failed performing database automigrate: %v", err)
		return str, err
	}

	str.db = db
	return str, nil
}

// dbFromURI
func dbFromURI(uri string) (string, string) {
	parts := strings.Split(uri, "://")
	if len(parts) != 2 {
		return "error", ""
	}
	return parts[0], parts[1]
}

// performDialectSpecific
func performDialectSpecific(db *gorm.DB, dialect string) error {
	switch dialect {
	case "sqlite3":
		err := db.Exec("PRAGMA journal_mode = WAL").Error
		if err != nil {
			return err
		}
		err = db.Exec("PRAGMA foreign_keys = ON").Error
		if err != nil {
			return err
		}
	case "mysql":
		// nothing , so far
	case "postgres":
		// nothing , so far
	case "mssql":
		// nothing , so far
	default:
		return fmt.Errorf("invalid dialect: %s", dialect)
	}
	return nil
}
