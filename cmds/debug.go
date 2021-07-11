package cmds

import (
	"fmt"
	"github.com/darmiel/macd-api/model"
	"github.com/darmiel/macd-api/yahoo"
	"github.com/urfave/cli/v2"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name: "debug",
		Action: func(ctx *cli.Context) (err error) {
			var h []*model.Historical
			if h, err = yahoo.RequestHistorical("ABST", "1d", "90d"); err != nil {
				return
			}
			for _, v := range h {
				fmt.Printf("%30s | H %8e | L %8e | C %8e | O %8e | V %8d\n",
					v.DayDate, v.High, v.Low, v.Close, v.Open, v.Volume)
			}
			fmt.Println(len(h), "records.")
			return
		},
		Flags: []cli.Flag{},
	})
}
