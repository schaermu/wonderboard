package main

import (
	"fmt"
	"os"

	"github.com/gookit/slog"
	"github.com/schaermu/wonderboard/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func main() {
	initConfig()

	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
		f.SetTemplate("[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}\n")
		f.TimeFormat = "2006-01-02 15:04:05"
	})

	server, err := cmd.NewServer(getFileSystem())
	if err != nil {
		slog.Fatal(err)
	}

	slog.Fatal(server.Http.Start(fmt.Sprintf(":%v", viper.GetString("PORT"))))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name "config.yml" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("BASE_URL", "http://localhost")
	viper.SetDefault("TRAEFIK_API_URL", "http://localhost:8080")
	viper.SetDefault("PORT", "3000")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
