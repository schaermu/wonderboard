/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"sync"
	"time"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/schaermu/docker-magic-dashboard/internal/docker"
	"github.com/spf13/cobra"
)

var wg *sync.WaitGroup

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		harvester := docker.New(
			docker.WithInterval(5*time.Second),
			docker.WithBaseUrl("http://localhost"),
			docker.WithTraefikApi("http://localhost:8080"),
		)

		if err := harvester.Start(); err != nil {
			slog.Panic(err)
		}

		slog.Info("started harvesting docker api...")

		e := echo.New()
		e.HideBanner = true
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogStatus:       true,
			LogURI:          true,
			LogMethod:       true,
			LogResponseSize: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				slog.Infof("request [%v] %v - %v (%v bytes)", v.Method, v.URI, v.Status, v.ResponseSize)
				return nil
			},
		}))
		e.GET("/api/current", func(c echo.Context) error {
			return c.JSON(200, harvester.GetCurrentContainers())
		})

		slog.Fatal(e.Start(":3000"))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
