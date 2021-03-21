package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

// FindEpicIdByEpicKey returns an epicId which is parent Issue.ID
func (myJiraImpl *myJiraImpl) FindEpicIDByEpicKey(epicKey string) (string, error) {
	issue, resp, err := myJiraImpl.underlying.Issue.Get(epicKey, nil)
	if err != nil {
		return "", xerrors.Errorf("an error occurred while searching for an issue with %s as an epicKey: %w", epicKey, gojira.NewJiraError(resp, err))
	}

	if issue == nil {
		return "", xerrors.Errorf("an error occurred. an issue with %s as an epicKey was not found.", epicKey)
	}
	return issue.ID, nil
}
