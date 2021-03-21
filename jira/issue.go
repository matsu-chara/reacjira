package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

type issueRequest struct {
	title       string
	description string
	reporter    *gojira.User
	issueType   string
	project     string
	epicID      string
}

func newIssueRequest(ticketReq TicketRequest, reporter *gojira.User, epicID string) issueRequest {
	return issueRequest{
		title:       ticketReq.Title,
		description: ticketReq.Description,
		reporter:    reporter,
		issueType:   ticketReq.IssueType,
		project:     ticketReq.Project,
		epicID:      epicID,
	}
}

func (request issueRequest) toGoJiraRequest() gojira.Issue {
	issue := gojira.Issue{
		Fields: &gojira.IssueFields{
			Reporter:    request.reporter,
			Description: request.description,
			Type:        gojira.IssueType{Name: request.issueType},
			Project:     gojira.Project{Key: request.project},
			Summary:     request.title,
		},
	}
	if request.epicID != "" {
		issue.Fields.Parent = &gojira.Parent{ID: request.epicID}
	}
	return issue
}

// FindEpicIdByEpicKey returns an epicId which is parent Issue.ID
func (myJiraImpl *myJiraImpl) CreateIssue(request issueRequest) (*gojira.Issue, error) {
	issueRequest := request.toGoJiraRequest()

	issue, resp, err := myJiraImpl.underlying.Issue.Create(&issueRequest)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while creating an issue: %w", gojira.NewJiraError(resp, err))
	}
	return issue, nil
}
