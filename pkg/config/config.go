package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GithubToken     string `envconfig:"INPUT_TOKEN"`
	GithubEventName string `envconfig:"GITHUB_EVENT_NAME"`
	GithubEventPath string `envconfig:"GITHUB_EVENT_PATH"`
	GithubRepo      string `envconfig:"INPUT_REPO"`
	GithubOwner     string `envconfig:"INPUT_OWNER"`
	WIPLabel        string `envconfig:"INPUT_LABEL"`
	WIPTitle        string `envconfig:"INPUT_TITLE"`
	wipCommits      string `envconfig:"INPUT_COMMITS"`
}

func New() (*Config, error) {
	conf := &Config{}
	return conf, envconfig.Process("", conf)
}

func (c *Config) WIPCommits() []string {
	return strings.Split(c.wipCommits, ",")
}
