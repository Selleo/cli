package console

import "io"

type App struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
}

var Regions = map[string][2]string{
	"eu-central-1": {"🇩🇪", "Frankfurt"},
	"ap-east-1":    {"🇭🇰", "Hong Kong"},
	"eu-west-3":    {"🇫🇷", "Paris"},
}

func New(out io.Writer, err io.Writer, in io.Reader) *App {
	return &App{
		stdout: out,
		stderr: err,
		stdin:  in,
	}
}
