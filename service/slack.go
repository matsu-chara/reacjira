package service

import (
	"reacjira/slack"

	goslack "github.com/slack-go/slack"
)

type SlackConfig struct {
	Host string
	Token string
}

type SlackService struct {
	config SlackConfig
	client slackClient
}

// an interface of a slack client
type slackClient interface {
	// TODO: add method
}

func NewSlack(config SlackConfig, rtm *goslack.RTM) *SlackService {
	client := slack.New(rtm)
	return &SlackService{config, client}
}

// TODO: client goslackを返り値に入れないようにラップする
// TODO: add method
// TODO: channel, linkをあわせてdescriptionにする？
// TODO: handler.goをusecaseに切り出す
// TODO: handler.goの長いメソッドなんとかする


func (slackService *SlackService) FindReporter() {
}

func (slackService *SlackService) FindTitle() {

}

func (slackService *SlackService) FindChannel() {

}

func (slackService *SlackService) FindLink() {

}
