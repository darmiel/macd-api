package cmds

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/model"
	"github.com/darmiel/macd-api/pg"
	"github.com/darmiel/macd-api/yahoo"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm/clause"
	"sync"
)

var cmdFetchHistoric = &cli.Command{
	Name: "historic",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "range", Value: "90d"},
		&cli.StringFlag{Name: "interval", Value: "1d"},
		&cli.IntFlag{Name: "max", Value: -1, Usage: "Max model"},
		&cli.IntFlag{Name: "gsize", Value: 7, Usage: "Group Size (The lower, the more threads: len(stocks) / gsize)"},
		&cli.BoolFlag{Name: "save", Value: true},
		&cli.BoolFlag{Name: "no-today", Value: false},
	},
	Action: func(ctx *cli.Context) (err error) {
		// Flags
		var (
			StgRange     = ctx.String("range")
			StgInterval  = ctx.String("interval")
			StgMax       = ctx.Int("max")
			StgGroupSize = ctx.Int("gsize")
			StgSave      = ctx.Bool("save")
			StgNoToday   = ctx.Bool("no-today")
		)

		// Database
		db := pg.MustPostgres(pg.FromCLI(ctx))
		if err = db.AutoMigrate(&model.Historic{}); err != nil {
			panic(err)
		}

		// Fetch
		fmt.Println(common.Info(), "Fetching ...")
		var sbl []*model.Symbol
		if err = db.Where("etf = false AND symbol SIMILAR TO '[A-Z]{1,5}'").
			Find(&sbl).Error; err != nil {
			return
		}

		if StgMax > 0 && len(sbl) > StgMax {
			fmt.Println(common.Info(), "Got", len(sbl), "model ( >", StgMax, "). Shrinking to", StgMax)
			sbl = sbl[:StgMax]
		}

		var (
			num     common.AtomicInt64
			dbmu    sync.Mutex
			skipped uint64
		)

		var errarr []string

		// start progressbar
		bar := pb.Full.Start(len(sbl))
		common.DistributedGoroutine(model.ConvertToGenericArray(sbl), StgGroupSize, func(arr []interface{}) {
			var historical []*model.Historic
			for _, a := range arr {
				m, o := a.(*model.Symbol)
				if !o {
					panic("o was no Symbol")
				}

				// request yahoo api
				historical, err = yahoo.RequestHistorical(m.Symbol, StgInterval, StgRange)
				if err != nil {
					msg := fmt.Sprintln(common.Error(), "Symbol", m.Symbol, "invalid response:", err)
					errarr = append(errarr, msg)
					if len(errarr) > 30 {
						// print if more than 30 errors
						fmt.Print(msg)
					}
					continue
				}
				// pb, save historical values
				bar.Increment()

				// remove from today
				if StgNoToday {
					var c []*model.Historic
					for _, h := range historical {
						if common.IsToday(h.DayDate) {
							skipped++
						} else {
							c = append(c, h)
						}
					}
					historical = c
				}

				num.Add(int64(len(historical)))

				if StgSave {
					// save to db
					dbmu.Lock()
					tx := db.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "symbol"}, {Name: "day_date"}},
						UpdateAll: true,
					}).CreateInBatches(historical, 1024)
					if tx.Error != nil {
						fmt.Println(common.Error(), "sql ::", tx.Error)
					}
					dbmu.Unlock()
				}
			}
		})

		// end progressbar
		bar.Finish()

		if len(errarr) > 0 {
			fmt.Println(common.Info(), "Got", len(errarr), "error responses:")
			fmt.Println("---")
			for _, e := range errarr {
				fmt.Print(e)
			}
			fmt.Println("---")
		}

		fmt.Println(common.Info(), "Loaded", num.Get(), "historical data ( skipped", skipped, ")")
		return nil
	},
}
