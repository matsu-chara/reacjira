package bot

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"golang.org/x/xerrors"

	"reacjira/config"
	"reacjira/service"
)

type CommandHandler struct {
	slack      *service.SlackService
	jira       *service.JiraService
	botProfile Profile
	botConfig  BotConfig
}

func NewCommandHandler(
	rtm *slack.RTM,
	botConfig BotConfig,
	botProfile Profile,
) (*CommandHandler, error) {
	slack := service.NewSlack(botConfig.Slack, rtm)

	jira, err := service.NewJira(botConfig.Jira)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize service.MyJiraService: %w", err)
	}

	return &CommandHandler{
		slack:      slack,
		jira:       jira,
		botProfile: botProfile,
		botConfig:  botConfig,
	}, nil
}

func (commandHandler *CommandHandler) HandleCommand(ev *slack.ReactionAddedEvent) {
	reacjira := commandHandler.searchReacjira(ev.Reaction)
	if reacjira == nil {
		// if a reacjira was not found, ignore this event.
		return
	}

	commandHandler.createTicket(ev, *reacjira)
}

// search reacjira which has the same name with the argument.
func (commandHandler *CommandHandler) searchReacjira(reaction string) *config.Reacjira {
	for _, r := range commandHandler.botConfig.Reacjiras.Values {
		if r.Emoji == reaction {
			return &r
		}
	}

	return nil
}

func (commandHandler *CommandHandler) createTicket(ev *slack.ReactionAddedEvent, reacjira config.Reacjira) {
	log.Printf("A slack reaction was received: channel: %s, timestamp: %s, user: %s", ev.Item.Channel, ev.Item.Timestamp, ev.User)
	log.Printf("prepare to create a ticket %+v", reacjira)

	msg, err := commandHandler.slack.FindMessage(ev)
	if err != nil {
		log.Printf("FindMessage error: %+v", err)
		if msg != nil {
			commandHandler.slack.SendError(err, ev.Item.Channel, msg.ThreadTimestamp)
		}
		return
	}
	log.Printf("slack message search finish: %s, %s, %s, reply:%d", msg.User, msg.Text, msg.Timestamp, msg.ReplyCount)

	var replyTo = msg.ThreadTimestamp
	if msg.ThreadTimestamp == "" {
		replyTo = msg.Timestamp
	}

	for _, reaction := range msg.Reactions {
		if reaction.Name == reacjira.Emoji && reaction.Count > 1 {
			log.Printf("multiple reacjira reactions were found(%s, %d). skip", reaction.Name, reaction.Count)
			return
		}
	}

	reacted, err := commandHandler.slack.CollectReactedMessageAttributes(msg, ev)
	if err != nil {
		log.Printf("CollectReactedMessageAttributes error: %+v", err)
		if msg != nil {
			commandHandler.slack.SendError(err, ev.Item.Channel, msg.ThreadTimestamp)
		}
		return
	}
	description := fmt.Sprintf(`auto generated
from: %s
at: %s
%s`, reacted.Message.Link, reacted.Message.Channel.Name, reacjira.Description)
	log.Printf("reporter:%s, channel: %s, title: %s", reacted.ReactedUser.Name, reacted.Message.Channel.Name, reacted.Message.Text)

	log.Printf("attempt to create a ticket.")
	ticketRequest := service.TicketRequest{
		Project:       reacjira.Project,
		IssueType:     reacjira.IssueType,
		EpicKey:       reacjira.EpicKey,
		ReporterEmail: reacted.ReactedUser.Profile.Email,
		Title:         reacted.Message.Text,
		Description:   description,
	}
	ticketURL, err := commandHandler.jira.CreateTicket(ticketRequest)
	log.Println("an ticket has been created: " + ticketURL)

	if err != nil {
		log.Printf("createTicket error: %+v", err)
		commandHandler.slack.SendError(err, ev.Item.Channel, replyTo)
		return
	}

	commandHandler.slack.SendMessage(ticketURL, ev.Item.Channel, replyTo)
}
