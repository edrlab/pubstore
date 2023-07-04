//go:build PGSQL
// +build PGSQL

package stor

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GormDialector(sqliteDsn string) gorm.Dialector {

	return postgres.Open(sqliteDsn)
}
