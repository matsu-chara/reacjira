package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

// a jira client wrapper which uses go-jira
type MyJiraClient struct {
	underlying *gojira.Client
}

func New(apiHost string, apiEmail string, apiToken string) (*MyJiraClient, error) {
	if apiHost == "" {
		return nil, xerrors.Errorf("an error occurred. an apiHost was empty")
	}
	if apiEmail == "" {
		return nil, xerrors.Errorf("an error occurred. an apiEmail was empty")
	}
	if apiToken == "" {
		return nil, xerrors.Errorf("an error occurred. an apiToken was empty")
	}

	tp := gojira.BasicAuthTransport{
		Username: apiEmail,
		Password: apiToken,
	}

	jiraClient, err := gojira.NewClient(tp.Client(), apiHost)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while creating the go-jira client: %w", err)
	}

	myJiraClient := MyJiraClient{jiraClient}
	return &myJiraClient, nil
}
