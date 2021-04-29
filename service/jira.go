package service

import (
	"fmt"
	"reacjira/jira"

	"golang.org/x/xerrors"
)

type JiraConfig struct {
	Host  string
	Email string
	Token string
}

type JiraService struct {
	config JiraConfig
	client jiraClient
}

type jiraClient interface {
	FindUserByEmail(email string) (*jira.User, error)
	FindEpicIDByEpicKey(epicKey string) (*jira.Epic, error)
	CreateIssue(request jira.IssueRequest) (*jira.Issue, error)
}

func NewJira(config JiraConfig) (*JiraService, error) {
	client, err := jira.New(config.Host, config.Email, config.Token)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred while initializiong jira client: %w", err)
	}
	return &JiraService{config, client}, nil
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
func (jiraService *JiraService) CreateTicket(ticketRequest TicketRequest) (string, error) {
	reporter, err := jiraService.client.FindUserByEmail(ticketRequest.ReporterEmail)
	if err != nil {
		return "", xerrors.Errorf("an error occurred in FindUserByEmail: %w", err)
	}

	epicID, err := jiraService.performEpicIdSearchOptional(ticketRequest.EpicKey)
	if err != nil {
		return "", xerrors.Errorf("an error occurred in performEpicIdSearchOptional: %w", err)
	}

	issueReq := ticketRequest.toIssueRequest(reporter, epicID)
	issue, err := jiraService.client.CreateIssue(issueReq)
	if err != nil {
		return "", xerrors.Errorf("an error occurred in CreateIssue: %w", err)
	}

	ticketURL := fmt.Sprintf("%s/browse/%s", jiraService.config.Host, issue.Key)
	return ticketURL, nil
}

// performEpicIdSearchOptional performs an epicId search only if an epic key is specified.
// this method may return (nil, nil) as a normal scenario when the epicKey was an empty string.
func (jiraService *JiraService) performEpicIdSearchOptional(epicKey string) (*jira.Epic, error) {
	if epicKey == "" {
		return nil, nil
	}

	epicID, err := jiraService.client.FindEpicIDByEpicKey(epicKey)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred in FindEpicIDByEpicKey: %w", err)
	}
	return epicID, err
}
