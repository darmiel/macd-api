package pg

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	*gorm.DB
}

func MustPostgres(dsn string) *Postgres {
	// var dsn =
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	return &Postgres{db}
}

///

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "pg-dsn"},
		&cli.StringFlag{Name: "pg-host"},
		&cli.StringFlag{Name: "pg-user"},
		&cli.StringFlag{Name: "pg-pass"},
		&cli.StringFlag{Name: "pg-db"},
		&cli.StringFlag{Name: "pg-port"},
		&cli.StringFlag{Name: "pg-tz"},
	}
}

func FromCLI(ctx *cli.Context) (dsn string) {
	// manual dsn
	if q := ctx.String("pg-dsn"); q != "" {
		return q
	}

	// default values
	var (
		PostgresHost = "localhost"
		PostgresUser = "postgres"
		PostgresPass = "123456"
		PostgresDb   = "postgres"
		PostgresPort = "45432"
		PostgresTZ   = "Europe/Berlin"
	)

	if q := ctx.String("pg-host"); q != "" {
		PostgresHost = q
	}
	if q := ctx.String("pg-user"); q != "" {
		PostgresUser = q
	}
	if q := ctx.String("pg-pass"); q != "" {
		PostgresPass = q
	}
	if q := ctx.String("pg-db"); q != "" {
		PostgresDb = q
	}
	if q := ctx.String("pg-port"); q != "" {
		PostgresPort = q
	}
	if q := ctx.String("pg-tz"); q != "" {
		PostgresTZ = q
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		PostgresHost, PostgresUser, PostgresPass, PostgresDb, PostgresPort, PostgresTZ)
}
