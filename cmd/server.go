package cmd

import (
	"net/http"
	"time"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/schaermu/wonderboard/internal/docker"
	"github.com/spf13/viper"
)

type Server struct {
	Http *echo.Echo
}

type ContainerResponse struct {
	Current []docker.DockerInfo `json:"current"`
}

func NewServer(embedFS http.FileSystem) (server *Server, err error) {
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
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: embedFS,
		HTML5:      true,
	}))
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

	// routes
	e.File("/", "public/index.html")
	e.GET("/api/current", func(c echo.Context) error {
		return c.JSON(200, ContainerResponse{Current: harvester.GetCurrentContainers()})
	})

	return &Server{
		Http: e,
	}, nil
}
