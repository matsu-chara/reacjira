package config

type BotProfile struct {
	BotUserID string
	BotName   string
}

type SlackCtx struct {
	Host string
}

type JiraCtx struct {
	Host  string
	Email string
	Token string
}
