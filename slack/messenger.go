package slack

func (client *SlackClient) SendMessage(msg string, channel string, threadTimestamp string) {
	outgoingMessage := client.rtm.NewOutgoingMessage(msg, channel)
	if threadTimestamp != "" {
		outgoingMessage.ThreadTimestamp = threadTimestamp
	}
	client.rtm.SendMessage(outgoingMessage)
}

func (client *SlackClient) SendMessages(msgs []string, channel string, threadTimestamp string) {
	for _, msg := range msgs {
		client.SendMessage(msg, channel, threadTimestamp)
	}
}
