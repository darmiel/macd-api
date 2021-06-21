package main

import (
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/nasdaq"
	"github.com/darmiel/macd-api/yahoo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

func main() {

	fmt.Println("Connecting to database ...")
	const dsn = "host=localhost user=postgres password=123456 dbname=postgres port=45432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Auto migrating ...")
	if err = db.AutoMigrate(&yahoo.Historical{}); err != nil {
		panic(err)
	}

	fmt.Println("Fetching securities ...")
	var sec []*nasdaq.NASDAQSecurity
	if sec, err = nasdaq.FetchNASDAQ(); err != nil {
		panic(err)
	}

	var seci = make([]interface{}, len(sec))
	for i, s := range sec {
		seci[i] = s
	}

	var (
		wrk int
		mu  sync.Mutex
	)

	common.DistributedGoroutine(seci, 5, func(i []interface{}) {
		mu.Lock()
		wrk++
		w := wrk
		mu.Unlock()

		for _, d := range i {
			s, o := d.(*nasdaq.NASDAQSecurity)
			if !o {
				continue
			}
			if s.ETF {
				continue
			}

			var hist []*yahoo.Historical
			if hist, err = yahoo.RequestHistorical(s.Symbol, "1d", "1y"); err != nil {
				fmt.Println(s.Symbol, w, "ERROR ::", err)
				continue
			}

			if len(hist) > 0 {
				db.Create(hist)
				fmt.Println(s.Symbol, w, ":", len(hist))
			}
		}
	})

}
