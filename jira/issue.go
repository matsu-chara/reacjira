package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

type IssueRequest struct {
	Title       string
	Description string
	Reporter    *User
	IssueType   string
	Project     string
	Epic        *Epic
}

func (request IssueRequest) toGoJiraRequest() (*gojira.Issue, error) {
	if request.Reporter == nil {
		return nil, xerrors.Errorf("request.repoter was nil.")
	}

	issue := gojira.Issue{
		Fields: &gojira.IssueFields{
			Reporter:    &request.Reporter.user,
			Description: request.Description,
			Type:        gojira.IssueType{Name: request.IssueType},
			Project:     gojira.Project{Key: request.Project},
			Summary:     request.Title,
		},
	}
	if request.Epic != nil {
		issue.Fields.Parent = &gojira.Parent{ID: request.Epic.ID}
	}
	return &issue, nil
}

type Issue struct {
	ID  string
	Key string
}

func fromGoJiraIssue(issue *gojira.Issue) *Issue {
	return &Issue{
		ID:  issue.ID,
		Key: issue.Key,
	}
}

// FindEpicIdByEpicKey returns an epicId which is parent Issue.ID
func (myJiraClient *MyJiraClient) CreateIssue(request IssueRequest) (*Issue, error) {
	issueRequest, err := request.toGoJiraRequest()
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while creating a GoJiraRequest: %w", err)
	}

	issue, resp, err := myJiraClient.underlying.Issue.Create(issueRequest)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while creating an issue: %w", gojira.NewJiraError(resp, err))
	}

	if err != nil {
		return nil, xerrors.Errorf("an error occurred. creating ticket was succeeded but an issue in a response was nil")
	}

	return fromGoJiraIssue(issue), nil
}
