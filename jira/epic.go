package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

type Epic struct {
	ID string
}

// FindEpicIdByEpicKey returns an epicID which is parent Issue.ID
func (jiraClient *JiraClient) FindEpicIDByEpicKey(epicKey string) (*Epic, error) {
	issue, resp, err := jiraClient.underlying.Issue.Get(epicKey, nil)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while searching for an issue with %s as an epicKey: %w", epicKey, gojira.NewJiraError(resp, err))
	}

	if issue == nil {
		return nil, xerrors.Errorf("an error occurred. an issue with %s as an epicKey was not found.", epicKey)
	}
	return &Epic{issue.ID}, nil
}
