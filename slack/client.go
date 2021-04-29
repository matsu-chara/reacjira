package slack

import (
	goslack "github.com/slack-go/slack"
)

type SlackClient struct {
	rtm *goslack.RTM
}

func New(rtm *goslack.RTM) *SlackClient {
	return &SlackClient{rtm}
}
