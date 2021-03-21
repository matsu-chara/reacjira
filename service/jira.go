package service

import (
	"fmt"
	"reacjira/jira"

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
	FindUserByEmail(email string) (*jira.User, error)
	FindEpicIDByEpicKey(epicKey string) (*jira.Epic, error)
	CreateIssue(request jira.IssueRequest) (*jira.Issue, error)
}

func NewJira(host string, email string, token string) (*MyJiraService, error) {
	jira, err := jira.New(host, email, token)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred: %w", err)
	}
	return &MyJiraService{host, email, token, jira}, nil
}

type TicketRequest struct {
	Project       string
	IssueType     string
	EpicKey       string
	ReporterEmail string
	Title         string
	Description   string
}

func (ticketRequest TicketRequest) toIssueRequest(reporter *jira.User, epic *jira.Epic) jira.IssueRequest {
	return jira.IssueRequest{
		Title:       ticketRequest.Title,
		Description: ticketRequest.Description,
		Reporter:    reporter,
		IssueType:   ticketRequest.IssueType,
		Project:     ticketRequest.Project,
		Epic:        epic,
	}
}

// CreateTicket creates a jira ticket and returns the URL of the created ticket
func (myjiraService *MyJiraService) CreateTicket(ticketRequest TicketRequest) (string, error) {
	reporter, err := myjiraService.client.FindUserByEmail(ticketRequest.ReporterEmail)
	if err != nil {
		return "", xerrors.Errorf("an error occurred in FindUserByEmail: %w", err)
	}

	epicID, err := myjiraService.performEpicIdSearchOptional(ticketRequest.EpicKey)
	if err != nil {
		return "", xerrors.Errorf("an error occurred in performEpicIdSearchOptional: %w", err)
	}

	issueReq := ticketRequest.toIssueRequest(reporter, epicID)
	issue, err := myjiraService.client.CreateIssue(issueReq)
	if err != nil {
		return "", xerrors.Errorf("an error occurred in CreateIssue: %w", err)
	}

	ticketURL := fmt.Sprintf("%s/browse/%s", myjiraService.host, issue.Key)
	return ticketURL, nil
}

// performEpicIdSearchOptional performs an epicId search only if an epic key is specified.
// this method may return (nil, nil) as a normal scenario when the epicKey was an empty string.
func (myjiraService *MyJiraService) performEpicIdSearchOptional(epicKey string) (*jira.Epic, error) {
	if epicKey == "" {
		return nil, nil
	}

	epicID, err := myjiraService.client.FindEpicIDByEpicKey(epicKey)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred in FindEpicIDByEpicKey: %w", err)
	}
	return epicID, err
}
