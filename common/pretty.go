package common

import (
	"github.com/muesli/termenv"
	"strings"
)

var (
	prf = termenv.ColorProfile()
)

func Smallify(in string, max int) string {
	in = strings.ReplaceAll(in, "\n", "$n$")
	if len(in) > max {
		in = in[:max] + "..."
	}
	return in
}

func Prefix(str string) termenv.Style {
	return termenv.String(" " + str + " ").
		Foreground(prf.Color("0")).
		Background(prf.Color("#D290E4"))
}

func Info() termenv.Style {
	return termenv.String(" INF ").Foreground(prf.Color("0")).Background(prf.Color("#3498db"))
}
func Error() termenv.Style {
	return termenv.String(" ERR ").Foreground(prf.Color("0")).Background(prf.Color("#E88388"))
}
