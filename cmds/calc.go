package cmds

import (
	"encoding/csv"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/darmiel/macd-api/calc"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/model"
	"github.com/darmiel/macd-api/pg"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name:      "ema",
		Category:  "Calculation",
		ArgsUsage: "[symbol]",
		Action: func(ctx *cli.Context) (err error) {
			// connect to db
			db := pg.MustPostgres(pg.FromCLI(ctx))

			//// symbols
			type symbolCount struct {
				Symbol string
				Count  int
			}
			var symbolCounts []*symbolCount

			// SELECT symbol, COUNT(symbol) FROM historics GROUP BY symbol HAVING COUNT(symbol) >= 90
			tx := db.Model(&model.Historic{}).
				Select("symbol, COUNT(symbol) AS count").
				Group("symbol").
				Having("COUNT(symbol) >= 90").
				Find(&symbolCounts)
			if tx.Error != nil {
				panic(tx.Error)
			}
			////

			var csvFile *os.File
			if csvFile, err = os.Create("current-90-day-ema.csv"); err != nil {
				panic(err)
			}
			defer csvFile.Close()

			bar := pb.Full.Start(len(symbolCounts))
			cd := make([][]string, len(symbolCounts)+1)

			/// HEADER
			cd[0] = []string{"Symbol"}
			for i := 0; i < 8; i++ {
				cd[0] = append(cd[0], "EMA10T"+strconv.Itoa(i))
			}
			for i := 0; i < 8; i++ {
				cd[0] = append(cd[0], "EMA35T"+strconv.Itoa(i))
			}
			///

			for i, c := range symbolCounts {
				bar.Increment()

				var quarter model.Quarter
				if quarter, err = db.GetHistorical90Data(c.Symbol); err != nil {
					fmt.Println(common.Error(), c.Symbol, "::", err)
					continue
				}

				cd[i+1] = append(cd[i+1], c.Symbol)
				for _, x := range calc.EMA(10, quarter).Seven() {
					cd[i+1] = append(cd[i+1], strconv.FormatFloat(x, 'f', 10, 64))
				}
				for _, x := range calc.EMA(35, quarter).Seven() {
					cd[i+1] = append(cd[i+1], strconv.FormatFloat(x, 'f', 10, 64))
				}
			}
			bar.Finish()
			fmt.Println("<<")

			writer := csv.NewWriter(csvFile)
			writer.Comma = '\t'
			if err = writer.WriteAll(cd); err != nil {
				panic(err)
			}

			fmt.Println("done writing")
			return
		},
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "sample-size", Value: 90},
			&cli.IntSliceFlag{Name: "days", Value: cli.NewIntSlice(10, 20, 35)},
		},
	})
}
