package jira

import (
	"fmt"

	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

type MyJiraService struct {
	host   string
	email  string
	token  string
	client myJira
}

// an interface of a jira client wrapper
type myJira interface {
	FindUserByEmail(email string) (*gojira.User, error)
	FindEpicIDByEpicKey(epicKey string) (string, error) // returns an EpicId( = Issue.ID) as a string
	CreateIssue(request issueRequest) (*gojira.Issue, error)
}

func New(host string, email string, token string) (*MyJiraService, error) {
	jiraClient, err := createMyJira(host, email, token)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred in New: %w", err)
	}
	return &MyJiraService{host, email, token, jiraClient}, nil
}

type TicketRequest struct {
	Project       string
	IssueType     string
	EpicKey       string
	ReporterEmail string
	Title         string
	Description   string
}

// CreateTicket creates a jira ticket and returns the URL of the created ticket
func (myjiraService *MyJiraService) CreateTicket(ticketRequest TicketRequest) (*string, error) {
	reporter, err := myjiraService.client.FindUserByEmail(ticketRequest.ReporterEmail)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred in FindUserByEmail: %w", err)
	}

	epicID, err := myjiraService.performEpicIdSearch(ticketRequest.EpicKey)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred in performEpicIdSearch: %w", err)
	}

	issueReq := newIssueRequest(ticketRequest, reporter, epicID)
	issue, err := myjiraService.client.CreateIssue(issueReq)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred in CreateIssue: %w", err)
	}

	ticketURL := fmt.Sprintf("%s/browse/%s", myjiraService.host, issue.Key)
	return &ticketURL, nil
}

// performEpicIdSearchWhenSpecified performs an epicId search only if an epic key is specified.
func (myjiraService *MyJiraService) performEpicIdSearch(epicKey string) (string, error) {
	if epicKey == "" {
		return "", nil
	}

	epicID, err := myjiraService.client.FindEpicIDByEpicKey(epicKey)
	if err != nil {
		return "", xerrors.Errorf("an error occurred in FindEpicIDByEpicKey: %w", err)
	}
	return epicID, err
}
