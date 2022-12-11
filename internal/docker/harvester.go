package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/slog"
)

var TRAEFIK_HOST_REGEX = regexp.MustCompile("\\(`(.+)`\\)")

type harvester struct {
	mu               sync.Mutex
	current          []DockerInfo
	ticker           *time.Ticker
	docker           *client.Client
	interval         time.Duration
	baseUrl          string
	traefikApiClient *resty.Client
}

func New(options ...func(*harvester)) *harvester {
	harv := &harvester{current: []DockerInfo{}}
	for _, o := range options {
		o(harv)
	}
	return harv
}

func (h *harvester) Start() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	h.ticker = time.NewTicker(h.interval)
	h.docker = cli

	go h.run()
	return nil
}

func (h *harvester) Stop() {
	h.docker.Close()
	h.ticker.Stop()
}

func WithInterval(interval time.Duration) func(*harvester) {
	return func(h *harvester) {
		h.interval = interval
	}
}

func WithBaseUrl(baseUrl string) func(*harvester) {
	return func(h *harvester) {
		h.baseUrl = baseUrl
	}
}

func WithTraefikApi(apiUrl string) func(*harvester) {
	return func(h *harvester) {
		h.traefikApiClient = resty.New().SetBaseURL(apiUrl)
		slog.Infof("initialized Traefik API client at %v", h.traefikApiClient.BaseURL)
	}
}

func (h *harvester) GetCurrentContainers() []DockerInfo {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.current
}

func (h *harvester) run() {
	for {
		select {
		case <-h.ticker.C:
			h.fetchCurrentData()
		}
	}
}

func (h *harvester) fetchCurrentData() {
	ctx := context.Background()
	listFilters := filters.NewArgs()
	listFilters.Add("status", "running")
	running, err := h.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: listFilters})
	if err != nil {
		slog.Panic(err)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

syncloop:
	for idx, curr := range h.current {
		for _, running := range running {
			if curr.ID == running.ID {
				break syncloop
			}
		}

		// container at idx not running, remove
		slog.Infof("detected stale container id %v", curr.ID)

		h.current = append(h.current[:idx], h.current[idx+1:]...)
	}

running:
	for _, container := range running {
		// sanitize shit
		name := strings.TrimLeft(container.Names[0], "/")

		for idx, val := range h.current {
			if val.ID == container.ID {
				// update existing record and break
				h.current[idx].Name = name
				h.current[idx].Image = container.Image

				if url, err := h.getTargetUrl(&container); err != nil {
					slog.Warnf("%v", err)
				} else {
					h.current[idx].TargetUrl = url
				}

				continue running
			}
		}

		record := DockerInfo{
			ID:    container.ID,
			Name:  name,
			Image: container.Image,
		}

		if url, err := h.getTargetUrl(&container); err != nil {
			slog.Warnf("%v", err)
		} else {
			record.TargetUrl = url
		}

		if project, ok := container.Labels["com.docker.compose.project"]; ok {
			record.Project = project
		}

		if service, ok := container.Labels["com.docker.compose.service"]; ok {
			record.Service = service
		}

		if version, ok := container.Labels["com.docker.compose.version"]; ok {
			record.Version = version
		}

		slog.Infof("added new container %v", record.Name)
		h.current = append(h.current, record)
	}
}

func (h *harvester) getTargetUrl(container *types.Container) (string, error) {
	if val, ok := container.Labels["traefik.enable"]; ok {
		if enabled, err := strconv.ParseBool(val); err == nil && enabled {
			// traefik-routed container, search explicit hostname
			for _, v := range container.Labels {
				if strings.Index(v, ".rule=Host(") > -1 {
					return fmt.Sprintf("http://%v", TRAEFIK_HOST_REGEX.FindAllStringSubmatch(val, -1)[0][1]), nil
				}
			}

			// no explicit hostname found, query traefik api
			var traefikHost string
			if h.traefikApiClient != nil {
				var result TraefikRouterResponse
				var svcName string
				if svcName, ok = container.Labels["com.docker.compose.service"]; !ok {
					slog.Warnf("could not find com.docker.compose.service label during traefik host lookup")
				}
				var targetUrl = fmt.Sprintf("/api/http/routers/%v@docker", svcName)
				if res, err := h.traefikApiClient.R().EnableTrace().Get(targetUrl); err != nil {
					slog.Errorf("request: %v", err)
				} else {
					if err := json.Unmarshal(res.Body(), &result); err != nil {
						slog.Errorf("unmarshal: %v", err)
					} else if len(result.Rule) > 0 {
						traefikHost = fmt.Sprintf("http://%v", TRAEFIK_HOST_REGEX.FindAllStringSubmatch(result.Rule, -1)[0][1])
					}
				}
			}

			if len(traefikHost) > 0 {
				return traefikHost, nil
			}
		}
	} else if len(container.Ports) == 1 {
		// only 1 port is mapped, use it with baseUrl
		port := container.Ports[0].PublicPort
		if port == 0 {
			port = container.Ports[0].PrivatePort
		}
		return fmt.Sprintf("http://%v:%v", h.baseUrl, port), nil
	} else {
		// cannot determine target url implicitly, check for explicit label
		if port, ok := container.Labels["dashboard.target.port"]; ok {
			return fmt.Sprintf("http://%v:%v", h.baseUrl, port), nil
		}
	}
	return "", fmt.Errorf("could not determine target host, configure using label dashboard.target.port")
}
