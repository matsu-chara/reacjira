package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

// an implementation of myJira using go-jira
type myJiraImpl struct {
	underlying *gojira.Client
}

func createMyJira(apiHost string, apiEmail string, apiToken string) (myJira, error) {
	tp := gojira.BasicAuthTransport{
		Username: apiEmail,
		Password: apiToken,
	}

	jiraClient, err := gojira.NewClient(tp.Client(), apiHost)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred: %w", err)
	}

	myJiraImpl := myJiraImpl{jiraClient}
	return &myJiraImpl, nil
}
