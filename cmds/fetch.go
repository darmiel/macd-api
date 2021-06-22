package cmds

import (
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/nasdaq"
	"github.com/urfave/cli/v2"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name:     "fetch-all",
		Aliases:  []string{"fa"},
		Category: "fetch",
		Action: func(ctx *cli.Context) (err error) {
			fmt.Println(common.Info(), "Fetching ...")
			var models []nasdaq.SecurityModel
			if models, err = nasdaq.FetchAll(); err != nil {
				return
			}
			var c int
			for _, m := range models {
				if m.IsETF() || !nasdaq.IsSymbolValid(m) {
					continue
				}
				if m.Exchange() != "NASDAQ" && m.Exchange() != nasdaq.NYSE {
					continue
				}
				fmt.Println(common.Prefix(fmt.Sprintf("%-5s", m.Symbol())), m.Exchange(), "::", m.SecurityName())
				c++
			}
			fmt.Println(common.Info(), "Fetched", len(models), "security models.", c, "valid.")
			return nil
		},
	})
}
