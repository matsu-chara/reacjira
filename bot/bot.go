package bot

import (
	"log"

	"reacjira/config"
	"reacjira/service"

	"github.com/slack-go/slack"
	goslack "github.com/slack-go/slack"
)

type Bot struct {
	config  BotConfig
	handler *CommandHandler
}

type BotConfig struct {
	Slack     service.SlackConfig
	Jira      service.JiraConfig
	Reacjiras *config.Reacjiras
}

type Profile struct {
	ID   string
	Name string
}

// Run returns int which means an error code.
func Run(config BotConfig) int {
	rtm := goslack.New(config.Slack.Token).NewRTM()
	go rtm.ManageConnection()

	var bot *Bot
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *goslack.HelloEvent:
			log.Print("A helloEvent was arrived")
		case *goslack.ConnectedEvent:
			log.Print("A connection to Slack has been established. Start initialization of the handler.")
			handler, err := NewCommandHandler(
				rtm,
				config,
				Profile{
					ID:   ev.Info.User.ID,
					Name: ev.Info.User.Name,
				},
			)
			if err != nil {
				log.Fatalf("an error occurred while initilizing CommandHandler. %+v", err)
			}
			bot = &Bot{config: config, handler: handler}
		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials. %+v", ev)
			return 1
		case *slack.ReactionAddedEvent:
			bot.handler.HandleCommand(ev)
		}
	}
	return 1
}
