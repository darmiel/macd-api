package cmds

import (
	"github.com/urfave/cli/v2"
)

func init() {
	App.Commands = append(App.Commands, &cli.Command{
		Name:    "fetch",
		Aliases: []string{"fa"},
		Subcommands: []*cli.Command{
			cmdFetchHistoric,
			cmdFetchSymbols,
		},
	})
}
