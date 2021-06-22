package cmds

import (
	"github.com/urfave/cli/v2"
)

var App = &cli.App{
	Name:                 "macds-api",
	Version:              "0.1.0",
	Description:          "Moving Average Convergence/Divergence",
	EnableBashCompletion: true,
	Authors: []*cli.Author{
		{
			Name:  "darmiel",
			Email: "hi@d2a.io",
		},
	},
}
