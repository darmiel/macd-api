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
	"strings"
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
			cd := make([][]string, len(symbolCounts)+2)

			/// HEADER
			cd[0] = []string{"Symbol", "Date"}
			for i := 7; i >= 0; i-- {
				cd[0] = append(cd[0], "EMA10T"+strconv.Itoa(i))
			}
			for i := 7; i >= 0; i-- {
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

				cd[i+2] = append(cd[i+2], c.Symbol)
				cd[i+2] = append(cd[i+2], quarter[89].DayDate.Format("02.01.2006"))
				for _, x := range calc.EMA(10, quarter).Seven().Reverse() {
					cd[i+2] = append(cd[i+2], formatFloat(x, ","))
				}
				for _, x := range calc.EMA(35, quarter).Seven().Reverse() {
					cd[i+2] = append(cd[i+2], formatFloat(x, ","))
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

func formatFloat(f float64, sep string) string {
	return strings.ReplaceAll(strconv.FormatFloat(f, 'f', 10, 64), ".", sep)
}
