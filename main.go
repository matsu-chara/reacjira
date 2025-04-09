package main

import (
	"log"
	"os"
	"reacjira/bot"
	"reacjira/config"

	"github.com/slack-go/slack"
)

func main() {
	log.Println("start reacjira.")

	configName := os.Getenv("REACJIRA_CONFIG_NAME")
	if configName == "" {
		configName = "config.toml"
	}

	log.Printf("start reading %s.", configName)
	loadedConfig, err := config.LoadConfigToml(configName)
	if err != nil {
		log.Printf("failed to load config: %s", err)
		os.Exit(1)
	}
	log.Printf("finish reading %s.", configName)

	reacjiraToml := os.Getenv("REACJIRA_REACJIRA_NAME")
	if reacjiraToml == "" {
		reacjiraToml = "reacjira.toml"
	}
	log.Printf("start reading %s.", reacjiraToml)
	reacjiras, err := config.LoadReacjiraToml(reacjiraToml)
	if err != nil {
		log.Printf("failed to load reacjira: %s", err)
		os.Exit(1)
	}
	log.Printf("finish reading %s.", configName)

	slackClient := slack.New(
		loadedConfig.SlackBotToken,
		slack.OptionAppLevelToken(loadedConfig.SlackAppToken),
		slack.OptionDebug(false),
	)

	botProfile, err := bot.GetSlackBotInfo(slackClient)
	if err != nil {
		log.Printf("failed to get bot info: %v", err)
		os.Exit(1)
	}
	log.Printf("got a bot info %v\n", botProfile)

	slackCtx := loadedConfig.ToSlackCtx()
	jiraCtx := loadedConfig.ToJiraCtx()

	os.Exit(bot.Run(slackClient, slackCtx, jiraCtx, botProfile, reacjiras))
}
