package cmd

import (
	"errors"
	"net/http"
	"time"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/schaermu/wonderboard/internal/docker"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

type Server struct {
	Http *echo.Echo
}

type ContainerResponse struct {
	Items []docker.DockerInfo `json:"items"`
}

type ContainerGroup struct {
	Name  string              `json:"name"`
	Items []docker.DockerInfo `json:"items"`
}

type ContainerGroupedResponse struct {
	Groups []ContainerGroup `json:"groups"`
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
		return c.JSON(200, ContainerResponse{Items: harvester.GetCurrentContainers()})
	})

	e.GET("/api/grouped", func(c echo.Context) error {
		groupParam := c.QueryParams().Get("by")
		if len(groupParam) > 0 {
			output := make([]ContainerGroup, 0)
			current := harvester.GetCurrentContainers()
			switch groupParam {
			case "project":
				for _, v := range current {
					idx := slices.IndexFunc(output, func(c ContainerGroup) bool { return c.Name == v.Project })
					if idx > -1 {
						output[idx].Items = append(output[idx].Items, v)
					} else {
						output = append(output, ContainerGroup{Name: v.Project, Items: []docker.DockerInfo{v}})
					}
				}
			}

			return c.JSON(200, ContainerGroupedResponse{Groups: output})
		} else {
			return c.JSON(400, errors.New("please supply by argument"))
		}
	})

	return &Server{
		Http: e,
	}, nil
}
