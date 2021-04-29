package main

import (
	"log"
	"os"
	"reacjira/bot"
	"reacjira/config"

	"golang.org/x/xerrors"
)

func main() {
	log.Println("start reacjira.")

	config, err := loadTomlConfig("REACJIRA_CONFIG_NAME", "config.toml")
	if err != nil {
		log.Fatalf("an error occurred while loading a config. %+v", err)
	}

	reacjiras, err := loadReacjiraToml("REACJIRA_REACJIRA_NAME", "reacjira.toml")
	if err != nil {
		log.Fatalf("an error occurred while loading reacjiras. %+v", err)
	}

	botConfig := bot.BotConfig{
		Slack:     config.ToSlackConfig(),
		Jira:      config.ToJiraConfig(),
		Reacjiras: reacjiras,
	}
	os.Exit(bot.Run(botConfig))
}

func loadTomlConfig(configEnvName string, defaultFileName string) (*config.TomlConfig, error) {
	configName := envGetOrElse(configEnvName, defaultFileName)

	log.Printf("start reading a config from %s.", configName)
	loadedConfig, err := config.LoadConfigToml(configName)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while reading toml config: %w", err)
	}
	log.Printf("finish reading a config from %s.", configName)
	return loadedConfig, nil
}

func loadReacjiraToml(configEnvName string, defaultFileName string) (*config.Reacjiras, error) {
	reacjiraName := envGetOrElse(configEnvName, defaultFileName)

	log.Printf("start reading reacjiras from %s.", reacjiraName)
	reacjiras, err := config.LoadReacjiraToml(reacjiraName)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while reading reacjiras toml: %w", err)
	}
	log.Printf("finish reading reacjiras from %s.", reacjiraName)
	return reacjiras, nil
}

func envGetOrElse(envName string, defaultValue string) string {
	v := os.Getenv(envName)
	if v == "" {
		return defaultValue
	}
	return v
}
