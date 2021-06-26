package common

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	PostgresHost = "localhost"
	PostgresUser = "postgres"
	PostgresPass = "123456"
	PostgresDb   = "postgres"
	PostgresPort = 45432
	PostgresTZ   = "Europe/Berlin"
)

func MustPostgres() *gorm.DB {
	var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		PostgresHost, PostgresUser, PostgresPass, PostgresDb, PostgresPort, PostgresTZ)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	return db
}
