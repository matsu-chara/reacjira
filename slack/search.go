package slack

import (
	goslack "github.com/slack-go/slack"
	"golang.org/x/xerrors"
)

func (client *SlackClient) SearchMessage(channel string, timestamp string) (*goslack.Message, error) {
	result, _, _, err := client.rtm.GetConversationReplies(
		&goslack.GetConversationRepliesParameters{
			ChannelID: channel,
			Timestamp: timestamp,
			Inclusive: true,
			Latest:    timestamp,
			Limit:     1,
		},
	)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while calling GetConversationReplies: %w", err)
	}
	if len(result) != 1 {
		return nil, xerrors.Errorf("an error occurred while searching slack with channel=%s, timestamp=%s. but result length was %d", channel, timestamp, len(result))
	}

	return &result[0], nil
}

func (client *SlackClient) SearchUser(userId string) (*goslack.User, error) {
	return client.rtm.GetUserInfo(userId)
}

func (client *SlackClient) SearchChannel(channelId string) (*goslack.Channel, error) {
	return client.rtm.GetConversationInfo(channelId, false)
}

func (client *SlackClient) GetPermLink(channelID string, messageTimestamp string) (string, error) {
	return client.rtm.GetPermalink(&goslack.PermalinkParameters{Channel: channelID, Ts: messageTimestamp})
}
