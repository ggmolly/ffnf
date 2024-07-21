package orm

import (
	"os"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	GormDB *gorm.DB
)

func InitDatabase() {
	var err error
	GormDB, err = gorm.Open(sqlite.Open(os.Getenv("DB_DSN")), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("failed to connect database " + err.Error())
	}

	err = GormDB.AutoMigrate()
	if err != nil {
		panic("failed to migrate database " + err.Error())
	}
}
