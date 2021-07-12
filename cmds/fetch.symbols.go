package cmds

import (
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/model"
	"github.com/darmiel/macd-api/pg"
	"github.com/darmiel/macd-api/security"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm/clause"
)

var cmdFetchSymbols = &cli.Command{
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
		if err = db.AutoMigrate(&model.Symbol{}); err != nil {
			return
		}

		// Fetch
		fmt.Println(common.Info(), "Fetching ...")
		var all []*model.Symbol
		if all, err = security.FetchAll(); err != nil {
			return
		}

		fmt.Println(common.Info(), "Fetched", len(all), "symbols")

		// Save to database
		if StgSave {
			fmt.Println(common.Info(), "Saving to db ...")
			tx := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "symbol"}},
				UpdateAll: true,
			}).CreateInBatches(all, 1024)
			if tx.Error != nil {
				fmt.Println(common.Error(), tx.Error)
			} else {
				fmt.Println(common.Info(), tx.RowsAffected, "rows affected")
			}
		}
		return
	},
}
