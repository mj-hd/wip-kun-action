package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	GithubToken string `envconfig:"TOKEN"`
	GithubSHA   string `envconfig:"SHA"`
	GithubRepo  string `envconfig:"REPO"`
	GithubOwner string `envconfig:"OWNER"`
	WIPLabel    string `envconfig:"LABEL"`
}

func New() (*Config, error) {
	conf := &Config{}
	return conf, envconfig.Process("INPUT", conf)
}
