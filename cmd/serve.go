package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/schaermu/wonderboard/internal/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var wg *sync.WaitGroup

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the application on port PORT or 3000",
	Long: `This command starts both the data harvester using the Docker and Traefik API's.

You can configure the behaviour using the following environment variables:
- BASE_URL = http://<YOUR_HOMELAB_IP_OR_HOSTNAME>
- TRAEFIK_API_URL = http://<YOUR_TRAFIK_SERVICENAME>
`,
	Run: func(cmd *cobra.Command, args []string) {
		harvester := docker.New(
			docker.WithInterval(5*time.Second),
			docker.WithBaseUrl(viper.GetString("BASE_URL")),
			docker.WithTraefikApi(viper.GetString("TRAEFIK_API_URL")),
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

		slog.Fatal(e.Start(fmt.Sprintf(":%v", viper.GetString("PORT"))))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
