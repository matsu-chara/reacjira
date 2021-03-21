package config

import (
	"log"
	"reacjira/service"

	toml "github.com/sioncojp/tomlssm"
	"golang.org/x/xerrors"
)

type TomlConfig struct {
	// slack
	SlackHost  string `toml:"slack_host"`
	SlackToken string `toml:"slack_token"`

	// jira
	JiraHost  string `toml:"jira_host"`
	JiraEmail string `toml:"jira_email"`
	JiraToken string `toml:"jira_token"`
}

func LoadConfigToml(filename string) (*TomlConfig, error) {
	log.Printf("try to load config from %s", filename)

	var config TomlConfig
	if _, err := toml.DecodeFile(filename, &config, "ap-northeast-1"); err != nil {
		return nil, xerrors.Errorf("failed to load config toml from %s: %w", filename, err)
	}

	return &config, nil
}

func (c TomlConfig) ToJiraConfig() service.JiraConfig {
	return service.JiraConfig{
		Host:  c.JiraHost,
		Email: c.JiraEmail,
		Token: c.JiraToken,
	}
}

func (c TomlConfig) ToSlackConfig() service.SlackConfig {
	return service.SlackConfig{
		Host: c.SlackHost,
		Token: c.SlackToken,
	}
}
