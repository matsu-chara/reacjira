package slack

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack"
)

// Messenger handle slack client
type Messenger struct {
	client *slack.Client
}

// New returns instance
func New(client *slack.Client) *Messenger {
	return &Messenger{client}
}

func (messenger *Messenger) SendMessage(msg string, channel string, threadTs string) error {
	msgOpt := slack.MsgOptionCompose(slack.MsgOptionText(msg, false))
	if threadTs != "" {
		msgOpt = slack.MsgOptionCompose(msgOpt, slack.MsgOptionTS(threadTs))
	}
	_, _, err := messenger.client.PostMessage(channel, msgOpt)
	return err
}

// SendMessages send messages
func (messenger *Messenger) SendMessages(msgs []string, channel string, threadTs string) error {
	errs := make([]error, 0, len(msgs))
	for _, msg := range msgs {
		err := messenger.SendMessage(msg, channel, threadTs)
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (messenger *Messenger) SearchMsg(channel string, timestamp string) (*slack.Message, error) {
	result, _, _, err := messenger.client.GetConversationReplies(
		&slack.GetConversationRepliesParameters{
			ChannelID: channel,
			Timestamp: timestamp,
			Inclusive: true,
			Latest:    timestamp,
			Limit:     1,
		})
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, fmt.Errorf("slack search with channel=%s, ts=%s. but result length is %d", channel, timestamp, len(result))
	}

	return &result[0], nil
}

func (messenger *Messenger) SearchUser(userId string) (*slack.User, error) {
	return messenger.client.GetUserInfo(userId)
}

func (messenger *Messenger) SearchChannel(channelId string) (*slack.Channel, error) {
	return messenger.client.GetConversationInfo(&slack.GetConversationInfoInput{ChannelID: channelId})
}

func (messenger *Messenger) GetPermLink(channelID string, messageTs string) (string, error) {
	return messenger.client.GetPermalink(&slack.PermalinkParameters{Channel: channelID, Ts: messageTs})
}
