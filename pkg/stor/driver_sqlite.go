//go:build !PGSQL
// +build !PGSQL

package stor

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GormDialector(sqliteDsn string) gorm.Dialector {

	return sqlite.Open(sqliteDsn)
}
