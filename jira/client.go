package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

// a jira client wrapper which uses go-jira
type JiraClient struct {
	underlying *gojira.Client
}

func New(apiHost string, apiEmail string, apiToken string) (*JiraClient, error) {
	if apiHost == "" {
		return nil, xerrors.Errorf("an error occurred. An apiHost was empty")
	}
	if apiEmail == "" {
		return nil, xerrors.Errorf("an error occurred. An apiEmail was empty")
	}
	if apiToken == "" {
		return nil, xerrors.Errorf("an error occurred. An apiToken was empty")
	}

	tp := gojira.BasicAuthTransport{
		Username: apiEmail,
		Password: apiToken,
	}

	gojiraClient, err := gojira.NewClient(tp.Client(), apiHost)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while creating the go-jira client: %w", err)
	}

	jiraClient := JiraClient{gojiraClient}
	return &jiraClient, nil
}
