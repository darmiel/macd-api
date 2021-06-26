package cmds

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/nasdaq"
	"github.com/darmiel/macd-api/yahoo"
	"github.com/urfave/cli/v2"
	"sync"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name:    "database",
		Aliases: []string{"db"},
		Action: func(ctx *cli.Context) (err error) {
			// settings
			var (
				StgRange     = ctx.String("range")
				StgInterval  = ctx.String("interval")
				StgMax       = ctx.Int("max")
				StgGroupSize = ctx.Int("gsize")
			)

			fmt.Println(common.Info(), "Fetching ...")
			var models []nasdaq.SecurityModel
			conn := nasdaq.MustConnection()
			if models, err = nasdaq.FetchAllAccepted(); err != nil {
				return
			}
			_ = conn.Quit() // quit FTP connection
			if len(models) > StgMax {
				fmt.Println(common.Info(), "Got", len(models), "models ( >", StgMax, "). Shrinking to", StgMax)
				models = models[:StgMax]
			}

			// copy models to interface{} array
			var ifa = make([]interface{}, len(models))
			for i, m := range models {
				ifa[i] = m
			}

			var (
				historicals []*yahoo.Historical
				hmu         sync.Mutex
				wg          sync.WaitGroup
			)

			var errarr []string

			// start progressbar
			bar := pb.Full.Start(len(ifa))
			common.DistributedGoroutine(ifa, StgGroupSize, func(arr []interface{}) {
				wg.Add(1)

				go func() {
					var historical []*yahoo.Historical
					for _, a := range arr {
						m, o := a.(nasdaq.SecurityModel)
						if !o {
							panic("o was no SecurityModel")
						}
						historical, err = yahoo.RequestHistorical(m.Symbol(), StgInterval, StgRange)
						if err != nil {
							msg := fmt.Sprintln(common.Error(), "Symbol", m.Symbol(), "invalid response:", err)
							errarr = append(errarr, msg)
							if len(errarr) > 30 {
								// print if more than 30 errors
								fmt.Print(msg)
							}
							continue
						}
						bar.Increment()
						hmu.Lock()
						historicals = append(historicals, historical...)
						hmu.Unlock()
					}
					wg.Done()
				}()
			})
			wg.Wait()

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

			fmt.Println(common.Info(), "Loaded", len(historicals), "historical data")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "range", Value: "90d"},
			&cli.StringFlag{Name: "interval", Value: "1d"},
			&cli.IntFlag{Name: "max", Value: 100, Usage: "Max models"},
			&cli.IntFlag{Name: "gsize", Value: 10, Usage: "Group Size (The lower, the more threads: len(stocks) / gsize)"},
		},
	})
}
