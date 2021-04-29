package service

import (
	"reacjira/slack"

	goslack "github.com/slack-go/slack"
	"golang.org/x/xerrors"
)

type SlackConfig struct {
	Host  string
	Token string
}

type SlackService struct {
	config SlackConfig
	client slackClient
}

type Reacted struct {
	ReactedUser *goslack.User
	Message     Message
}
type Message struct {
	Text    string
	Link    string
	Channel *goslack.Channel
}

type slackClient interface {
	SendMessages(msgs []string, channel string, threadTimestamp string)
	SearchMessage(channel string, timestamp string) (*goslack.Message, error)
	SearchUser(userId string) (*goslack.User, error)
	SearchChannel(channelId string) (*goslack.Channel, error)
	GetPermLink(channelID string, messageTimestamp string) (string, error)
}

func NewSlack(config SlackConfig, rtm *goslack.RTM) *SlackService {
	client := slack.New(rtm)
	return &SlackService{config, client}
}

// TODO: client goslackを返り値に入れないようにラップする
// TODO: handler.goをusecaseに切り出す
// TODO: handler.goの長いメソッドなんとかする

func (slackService *SlackService) FindMessage(ev *goslack.ReactionAddedEvent) (*goslack.Message, error) {
	msg, err := slackService.client.SearchMessage(ev.Item.Channel, ev.Item.Timestamp)
	if err != nil {
		return nil, xerrors.Errorf("FindMessage error: %w", err)
	}
	return msg, nil
}

func (slackService *SlackService) CollectReactedMessageAttributes(msg *goslack.Message, ev *goslack.ReactionAddedEvent) (*Reacted, error) {
	reactedUser, err := slackService.client.SearchUser(msg.User)
	if err != nil {
		return nil, xerrors.Errorf("can't find reactedUser: %w", err)
	}

	channel, err := slackService.client.SearchChannel(ev.Item.Channel)
	if err != nil {
		return nil, xerrors.Errorf("can't find reacted channel: %w", err)
	}

	link, err := slackService.client.GetPermLink(ev.Item.Channel, msg.Timestamp)
	if err != nil {
		return nil, xerrors.Errorf("can't find permLink: %w", err)
	}
	return &Reacted{
		reactedUser,
		Message{
			msg.Text,
			link,
			channel,
		},
	}, nil
}

func (slackService *SlackService) SendError(err error, channel string, threadTimestamp string) {
	slackService.client.SendMessages([]string{err.Error()}, channel, threadTimestamp)
}

func (slackService *SlackService) SendMessage(str string, channel string, threadTimestamp string) {
	slackService.client.SendMessages([]string{str}, channel, threadTimestamp)
}
