package config

type Profile struct {
	ID   string
	Name string
}

type SlackCtx struct {
	Host string
}

type JiraCtx struct {
	Host  string
	Email string
	Token string
}
