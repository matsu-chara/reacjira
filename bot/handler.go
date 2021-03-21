package bot

import (
	"fmt"
	"log"
	"math"
	"strings"

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
	for _, r := range commandHandler.botConfig.Reacjiras {
		if r.Emoji == reaction {
			return &r
		}
	}

	return nil
}

func (commandHandler *CommandHandler) createTicket(ev *slack.ReactionAddedEvent, reacjira config.Reacjira) {
	log.Printf("A slack reaction was received: channel: %s, timestamp: %s, user: %s", ev.Item.Channel, ev.Item.Timestamp, ev.User)
	log.Printf("prepare to create a ticket %+v", reacjira)

	// fetch slack message
	msg, err := findMessage(commandHandler.slack, ev)
	if err != nil {
		log.Printf("findMessage error: %+v", err)
		if msg != nil {
			commandHandler.slack.SendMessage(err.Error(), ev.Item.Channel, msg.ThreadTimestamp)
		}
		return
	}
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

	reporter, err := commandHandler.slack.SearchUser(msg.User)
	if err != nil {
		log.Printf("searchUser error: %+v", err)
		commandHandler.slack.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	channel, err := commandHandler.slack.SearchChannel(ev.Item.Channel)
	if err != nil {
		log.Printf("searchChannel error: %+v", err)
		commandHandler.slack.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	link, err := getPermLink(commandHandler.slack, ev, msg)
	if err != nil {
		log.Printf("getPermLink error: %+v", err)
		commandHandler.slack.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	title := strings.Replace(msg.Text, "\n", " ", -1)
	limit := int(math.Min(float64(len(title)), 200))
	title = title[0:limit]

	description := fmt.Sprintf(`auto generated
from: %s
at: %s
%s`, link, channel.Name, reacjira.Description)
	log.Printf("reporter:%s, channel: %s, title: %s", reporter.Name, channel.Name, title)

	log.Printf("attempt to create a ticket.")
	ticketRequest := service.TicketRequest{
		Project:       reacjira.Project,
		IssueType:     reacjira.IssueType,
		EpicKey:       reacjira.EpicKey,
		ReporterEmail: reporter.Profile.Email,
		Title:         title,
		Description:   description,
	}
	ticketURL, err := commandHandler.jira.CreateTicket(ticketRequest)
	log.Println("an ticket has been created: " + ticketURL)

	if err != nil {
		log.Printf("createTicket error: %+v", err)
		commandHandler.slack.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	commandHandler.slack.SendMessages([]string{ticketURL}, ev.Item.Channel, replyTo)
}

func findMessage(slack *service.SlackService, ev *slack.ReactionAddedEvent) (*slack.Message, error) {
	msg, err := slack.SearchMsg(ev.Item.Channel, ev.Item.Timestamp)
	if err != nil {
		log.Printf("findMessage error: %+v", err)
		return nil, err
	}

	log.Printf("slack message search finish: %s, %s, %s, reply:%d", msg.User, msg.Text, msg.Timestamp, msg.ReplyCount)
	return msg, nil
}

func getPermLink(slack *service.SlackService, ev *slack.ReactionAddedEvent, msg *slack.Message) (string, error) {
	link, err := slack.GetPermLink(ev.Item.Channel, msg.Timestamp)
	if err != nil {
		log.Printf("getPermLink error: %+v", err)
		return "", err
	}

	log.Printf("got slack permlink: %s", link)
	return link, nil
}
