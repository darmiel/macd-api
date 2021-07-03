package cmds

import (
	"errors"
	"fmt"
	"github.com/darmiel/macd-api/calc"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/pg"
	"github.com/urfave/cli/v2"
	"strings"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name:      "ema",
		Category:  "Calculation",
		ArgsUsage: "[symbol]",
		Action: func(ctx *cli.Context) error {

			var (
				StgSymbol     = ctx.Args().Slice()
				StgSampleSize = ctx.Int("sample-size")
				StgDays       = ctx.IntSlice("days")
			)
			if len(StgSymbol) == 0 {
				return errors.New("symbol/s required")
			}

			// connect to db
			db := pg.MustPostgres(pg.FromCLI(ctx))

			for _, symbol := range StgSymbol {
				symbol = strings.ToUpper(symbol)

				// fetch historical data
				fmt.Print(common.Info(), " Fetching symbol ", symbol, " ...")
				data, err := db.GetHistorical90Data(symbol)
				if err != nil {
					fmt.Println("<->", common.Error(), err)
					continue
				}

				fmt.Println("<->", common.Info(), len(data), "records")

				// check size
				if len(data) < StgSampleSize {
					fmt.Println(common.Error(), len(data), "records returned, but", StgSampleSize, "required")
					continue
				}

				for _, day := range StgDays {
					ema, err := calc.EMA(day, data)
					if err != nil {
						fmt.Println(common.Error(), "Failed to calculate for", day, "days:", err)
						continue
					}
					fmt.Println(common.Info(), "EMA", day, "|", ema.Seven(), "..", ema.Min(14), ema.Min(0))
				}
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "sample-size", Value: 90},
			&cli.IntSliceFlag{Name: "days", Value: cli.NewIntSlice(10, 20, 35)},
		},
	})
}
