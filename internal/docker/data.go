package docker

type DockerInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Project   string `json:"project"`
	TargetUrl string `json:"targetUrl"`
}

type TraefikRouterResponse struct {
	EntryPoints []string
	Service     string
	Rule        string
	Status      string
	Using       []string
	Name        string
	Provider    string
}
