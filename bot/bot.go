package bot

import (
	"fmt"
	"log"

	"reacjira/bot/handler"
	"reacjira/config"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// GetSlackBotInfo get slack bot info
func GetSlackBotInfo(api *slack.Client) (config.BotProfile, error) {
	resp, err := api.AuthTest()
	if err != nil {
		return config.BotProfile{}, fmt.Errorf("an error occurred during authTest api call: %w", err)
	}
	return config.BotProfile{
		BotUserID: resp.UserID,
		BotName:   resp.User}, nil
}

// Run activate bot routine
func Run(slackClient *slack.Client, slackCtx config.SlackCtx, jiraCtx config.JiraCtx, botProfile config.BotProfile, reacjiras []config.Reacjira) int {
	slackSocket := socketmode.New(
		slackClient,
		socketmode.OptionDebug(false),
	)

	commandHandler := handler.NewCommandHandler(
		slackClient,
		slackCtx,
		jiraCtx,
		botProfile,
		reacjiras,
	)

	go func() {
		for evt := range slackSocket.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				log.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				log.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				log.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("failed to cast Data to EventsAPIEvent. %+v\n", evt)
					continue
				}
				slackSocket.Ack(*evt.Request) // we must return an ack within 3 seconds or Slack will retry

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.ReactionAddedEvent:
						commandHandler.HandleCommand(ev)
					}
				default:
					log.Printf("unsupported Events API event received: %v\n", eventsAPIEvent)
				}
			default:
				// 全く関係ないイベントなども全部ログに出てしまうので普段はコメントアウトしておく
				// log.Printf("Unexpected Event(Ignored): %v\n", msg.Data)
			}
		}
	}()

	err := slackSocket.Run()
	if err != nil {
		log.Printf("failed to run slack socket mode: %v\n", err)
		return 1
	}
	return 0
}
