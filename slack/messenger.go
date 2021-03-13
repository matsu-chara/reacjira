package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// Messenger handle slack rtm
type Messenger struct {
	rtm *slack.RTM
}

// New returns instance
func New(rtm *slack.RTM) *Messenger {
	return &Messenger{rtm}
}

// SendMessages send messages
func (messenger *Messenger) SendMessages(msgs []string, channel string, threadTs string) {
	for _, msg := range msgs {
		messenger.sendMessage(msg, channel, threadTs)
	}
}

func (messenger *Messenger) sendMessage(msg string, channel string, threadTs string) {
	outgoingMessage := messenger.rtm.NewOutgoingMessage(msg, channel)
	if threadTs != "" {
		outgoingMessage.ThreadTimestamp = threadTs
	}
	messenger.rtm.SendMessage(outgoingMessage)
}

func (messenger *Messenger) SearchMsg(channel string, timestamp string) (*slack.Message, error) {
	result, _, _, err := messenger.rtm.GetConversationReplies(
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
	return messenger.rtm.GetUserInfo(userId)
}

func (messenger *Messenger) SearchChannel(channelId string) (*slack.Channel, error) {
	return messenger.rtm.GetConversationInfo(channelId, false)
}

func (messenger *Messenger) GetPermLink(channelID string, messageTs string) (string, error) {
	return messenger.rtm.GetPermalink(&slack.PermalinkParameters{Channel: channelID, Ts: messageTs})
}
