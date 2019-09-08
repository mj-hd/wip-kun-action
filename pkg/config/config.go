package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	GithubToken string `envconfig:"GITHUB_TOKEN"`
	GithubSHA   string `envconfig:"GITHUB_SHA"`
	GithubRepo  string `envconfig:"GITHUB_REPO"`
	GithubOwner string `envconfig:"GITHUB_OWNER"`
	WIPLabel    string `envconfig:"WIP_LABEL"`
}

func New() (*Config, error) {
	conf := &Config{}
	return conf, envconfig.Process("", conf)
}
