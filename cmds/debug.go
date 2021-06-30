package cmds

import (
	"fmt"
	"github.com/darmiel/macd-api/pg"
	"github.com/urfave/cli/v2"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name: "debug",
		Action: func(ctx *cli.Context) (err error) {
			db := pg.MustPostgres(pg.FromCLI(ctx))
			fmt.Println("fetching symbols ...")
			data, err := db.FindHistoricalsWithMinData(90)
			if err != nil {
				panic(err)
			}
			fmt.Println("DATA:")
			fmt.Println()
			for k, v := range data {
				fmt.Println(k, "::")
				for _, a := range v {
					fmt.Printf("%+v | ", a)
				}
				fmt.Println()
				fmt.Println("  -> len:", len(v))
				fmt.Println()
			}
			fmt.Println("len:", len(data))
			return
		},
		Flags: []cli.Flag{},
	})
}
