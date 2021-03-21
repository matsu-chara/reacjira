package service

import (
	"reacjira/slack"

	goslack "github.com/slack-go/slack"
)

type SlackService struct {
	host   string
	client slackClient
}

// an interface of a slack client
type slackClient interface {
	// TODO: add method
}

func NewSlack(host string, rtm *goslack.RTM) *SlackService {
	client := slack.New(rtm)
	return &SlackService{host, client}
}

// TODO: add method
