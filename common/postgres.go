package common

import (
	"fmt"
	"github.com/darmiel/macd-api/yahoo"
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

type Postgres struct {
	*gorm.DB
}

func MustPostgres() *Postgres {
	var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		PostgresHost, PostgresUser, PostgresPass, PostgresDb, PostgresPort, PostgresTZ)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	return &Postgres{db}
}

func (p *Postgres) FindAllSymbols() (res []string, err error) {
	tx := p.Model(&yahoo.Historical{}).
		Distinct("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindAllSymbolsWithMinData(num int) (res []string, err error) {
	tx := p.Model(&yahoo.Historical{}).
		Select("symbol").
		Having("count(symbol) >= ?", num).
		Group("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindHistoricalsWithMinData(num int) (res map[string][]*yahoo.Historical, err error) {
	h := make([]*yahoo.Historical, 0)
	tx := p.Model(&yahoo.Historical{}).
		Raw("SELECT * FROM historicals h WHERE symbol IN (SELECT symbol FROM historicals GROUP BY symbol HAVING COUNT(symbol) >= ?)", num).
		Find(&h)
	if err = tx.Error; err != nil {
		return
	}
	res = make(map[string][]*yahoo.Historical)
	for _, x := range h {
		if _, o := res[x.Symbol]; !o {
			res[x.Symbol] = make([]*yahoo.Historical, 0)
		}
		res[x.Symbol] = append(res[x.Symbol], x)
	}
	return
}
