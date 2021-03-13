package config

import (
	"fmt"
	"log"

	toml "github.com/sioncojp/tomlssm"
)

type Config struct {
	// slack
	SlackHost  string `toml:"slack_host"`
	SlackToken string `toml:"slack_token"`

	// jira
	JiraHost  string `toml:"jira_host"`
	JiraEmail string `toml:"jira_email"`
	JiraToken string `toml:"jira_token"`
}

func LoadConfigToml(filename string) (*Config, error) {
	log.Printf("try to load config from %s", filename)

	var config Config
	if _, err := toml.DecodeFile(filename, &config, "ap-northeast-1"); err != nil {
		return nil, fmt.Errorf("failed to load config toml from %s: %w", filename, err)
	}

	return &config, nil
}

func (c Config) ToJiraCtx() JiraCtx {
	return JiraCtx{
		Host:  c.JiraHost,
		Email: c.JiraEmail,
		Token: c.JiraToken,
	}
}

func (c Config) ToSlackCtx() SlackCtx {
	return SlackCtx{
		Host: c.SlackHost,
	}
}
