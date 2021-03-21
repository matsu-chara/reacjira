package bot

import (
	"log"

	"reacjira/bot/handler"
	"reacjira/config"

	"github.com/slack-go/slack"
)

type Bot struct {
	commandHandler *handler.CommandHandler
}

// Run activate bot routine
func Run(api *slack.Client, slackCtx config.SlackCtx, jiraCtx config.JiraCtx, reacjiras []config.Reacjira) int {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// initialized when connected
	var bot *Bot

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			log.Print("hello")
		case *slack.ConnectedEvent:
			log.Print("connected")
			handler, err := handler.NewCommandHandler(
				rtm,
				slackCtx,
				jiraCtx,
				config.Profile{
					ID:   ev.Info.User.ID,
					Name: ev.Info.User.Name,
				},
				reacjiras,
			)
			if err != nil {
				log.Fatalf("an error occurred while initilizing CommandHandler. %+v", err)
			}
			bot = &Bot{commandHandler: handler}
		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials. %v", ev)
			return 1
		case *slack.ReactionAddedEvent:
			bot.commandHandler.HandleCommand(ev)
		}
	}
	return 1
}
