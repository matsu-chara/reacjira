package handler

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/slack-go/slack"

	"reacjira/config"
	myjira "reacjira/jira"
	myslack "reacjira/slack"
)

type CommandHandler struct {
	slackMessenger *myslack.Messenger
	jira           *myjira.MyJira
	slackCtx       config.SlackCtx
	botProfile     config.Profile
	reacjiras      []config.Reacjira
}

func NewCommandHandler(
	rtm *slack.RTM,
	slackCtx config.SlackCtx,
	jiraCtx config.JiraCtx,
	botProfile config.Profile,
	reacjiras []config.Reacjira,
) *CommandHandler {
	return &CommandHandler{
		myslack.New(rtm),
		myjira.New(jiraCtx.Host, jiraCtx.Email, jiraCtx.Token),
		slackCtx,
		botProfile,
		reacjiras,
	}
}

func (commandHandler *CommandHandler) HandleCommand(ev *slack.ReactionAddedEvent) {
	// search reacjira
	var found config.Reacjira

	for _, r := range commandHandler.reacjiras {
		if r.Emoji == ev.Reaction {
			found = r
		}
	}

	if found.Emoji == "" {
		// if reacjira was not found, then ignore this command
		return
	}

	commandHandler.createTicket(ev, found)
}

func (commandHandler *CommandHandler) createTicket(ev *slack.ReactionAddedEvent, reacjira config.Reacjira) {
	log.Printf("A slack reaction was received: channel: %s, timestamp: %s, user: %s", ev.Item.Channel, ev.Item.Timestamp, ev.User)
	log.Printf("prepare to create a ticket %+v", reacjira)

	// fetch slack message
	msg, err := findMessage(commandHandler.slackMessenger, ev)
	if err != nil {
		log.Printf("findMessage error: %+v", err)
		if msg != nil {

			commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, msg.ThreadTimestamp)
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

	reporter, err := commandHandler.slackMessenger.SearchUser(msg.User)
	if err != nil {
		log.Printf("searchUser error: %+v", err)
		commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	channel, err := commandHandler.slackMessenger.SearchChannel(ev.Item.Channel)
	if err != nil {
		log.Printf("searchChannel error: %+v", err)
		commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	link, err := getPermLink(commandHandler.slackMessenger, ev, msg)
	if err != nil {
		log.Printf("getPermLink error: %+v", err)
		commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
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
	ticket, err := commandHandler.jira.CreateTicket(
		reacjira.Project,
		reacjira.IssueType,
		reacjira.EpicKey,
		reporter.Profile.Email,
		title,
		description,
	)
	if err != nil {
		log.Printf("createTicket error: %+v", err)
		commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	commandHandler.slackMessenger.SendMessages([]string{*ticket}, ev.Item.Channel, replyTo)
}

func findMessage(slackMessenger *myslack.Messenger, ev *slack.ReactionAddedEvent) (*slack.Message, error) {
	msg, err := slackMessenger.SearchMsg(ev.Item.Channel, ev.Item.Timestamp)
	if err != nil {
		log.Printf("findMessage error: %+v", err)
		return nil, err
	}

	log.Printf("slack message search finish: %s, %s, %s, reply:%d", msg.User, msg.Text, msg.Timestamp, msg.ReplyCount)
	return msg, nil
}

func getPermLink(slackMessenger *myslack.Messenger, ev *slack.ReactionAddedEvent, msg *slack.Message) (string, error) {
	link, err := slackMessenger.GetPermLink(ev.Item.Channel, msg.Timestamp)
	if err != nil {
		log.Println("getPermLink error: %+v", err)
		return "", err
	}

	log.Printf("got slack permlink: %s", link)
	return link, nil
}
