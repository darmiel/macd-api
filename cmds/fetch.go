package cmds

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/models"
	"github.com/darmiel/macd-api/nasdaq"
	"github.com/darmiel/macd-api/pg"
	"github.com/darmiel/macd-api/yahoo"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm/clause"
	"sync"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name:    "fetch",
		Aliases: []string{"fa"},
		Subcommands: []*cli.Command{
			{
				Name: "symbols",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "save", Value: true},
				},
				Action: func(ctx *cli.Context) (err error) {
					// Flags
					var (
						StgSave = ctx.Bool("save")
					)

					// Database
					db := pg.MustPostgres(pg.FromCLI(ctx))
					if err = db.AutoMigrate(&models.Symbol{}); err != nil {
						return
					}

					// Fetch
					fmt.Println(common.Info(), "Fetching ...")
					var sbl []nasdaq.SecurityModel
					if sbl, err = nasdaq.FetchAll(); err != nil {
						return
					}

					bar := pb.Full.Start64(int64(len(sbl)))
					var mds []*models.Symbol
					for _, m := range sbl {
						mds = append(mds, nasdaq.ToSymbolModel(m))
						bar.Increment()
					}
					bar.Finish()
					fmt.Println(common.Info(), "Fetched", len(sbl), "symbols")

					if StgSave {
						fmt.Println(common.Info(), "Saving to db ...")
						tx := db.Clauses(clause.OnConflict{
							Columns:   []clause.Column{{Name: "symbol"}},
							UpdateAll: true,
						}).CreateInBatches(mds, 1024)
						if tx.Error != nil {
							fmt.Println(common.Error(), tx.Error)
						} else {
							fmt.Println(common.Info(), tx.RowsAffected, "rows affected")
						}
					}
					return
				},
			},
			{
				Name: "historical",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "range", Value: "90d"},
					&cli.StringFlag{Name: "interval", Value: "1d"},
					&cli.IntFlag{Name: "max", Value: -1, Usage: "Max models"},
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
					if err = db.AutoMigrate(&models.Historical{}); err != nil {
						panic(err)
					}

					// Fetch
					fmt.Println(common.Info(), "Fetching ...")
					var sbl []*models.Symbol
					if err = db.Where("etf = false AND symbol SIMILAR TO '[A-Z]{1,5}'").
						Find(&sbl).Error; err != nil {
						return
					}

					if StgMax > 0 && len(sbl) > StgMax {
						fmt.Println(common.Info(), "Got", len(sbl), "models ( >", StgMax, "). Shrinking to", StgMax)
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
					common.DistributedGoroutine(models.ConvertToGenericArray(sbl), StgGroupSize, func(arr []interface{}) {
						var historical []*models.Historical
						for _, a := range arr {
							m, o := a.(*models.Symbol)
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
								var c []*models.Historical
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
			},
		},
	})
}
