package main

import (
	"github.com/gookit/slog"
	"github.com/schaermu/docker-magic-dashboard/cmd"
)

func main() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
		f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}\n")
		f.TimeFormat = "2006-01-02 15:04:05"
	})

	cmd.Execute()
}
