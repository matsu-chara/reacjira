package jira

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/andygrunwald/go-jira"
)

type MyJira struct {
	host  string
	email string
	token string
}

func New(host string, email string, token string) *MyJira {
	return &MyJira{host, email, token}
}

func (myjira *MyJira) CreateTicket(
	project string,
	issueType string,
	epicKey string,
	reporterEmail string,
	title string,
	description string,
) (*string, error) {
	base := myjira.host
	tp := jira.BasicAuthTransport{
		Username: myjira.email,
		Password: myjira.token,
	}
	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		return nil, fmt.Errorf("error. during jiraClient. error=%s", err.Error())
	}

	reporters, res, err := myjira.neoFindUserByEmail(jiraClient, reporterEmail)
	if err != nil {
		b, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			log.Printf("user.find is error and failed to read response body. err1=%s, err2=%s", err.Error(), err2.Error())
			return nil, err2
		}

		log.Printf("user.find error response repoters %v", string(b))
		return nil, err
	}
	if len(reporters) != 1 {
		err = fmt.Errorf("error. found multiple users on email address search. len=%d", len(reporters))
		return nil, err
	}

	epicId := ""
	if epicKey != "" {
		issue, res, err := jiraClient.Issue.Get(epicKey, nil)
		if err != nil {
			b, err2 := ioutil.ReadAll(res.Body)
			if err2 != nil {
				log.Printf("epic issue.get is error and failed to read response body. err1=%s, err2=%s", err.Error(), err2.Error())
				return nil, err2
			}

			log.Printf("epic issue.get error response repoters %v", string(b))
			return nil, err
		}
		epicId = issue.ID
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Reporter:    &reporters[0],
			Description: description,
			Type:        jira.IssueType{Name: issueType},
			Project:     jira.Project{Key: project},
			Summary:     title,
		},
	}
	if epicId != "" {
		i.Fields.Parent = &jira.Parent{ID: epicId}
	}

	issue, res, err := jiraClient.Issue.Create(&i)
	if err != nil {
		b, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			log.Printf("issue.create is error and failed to read response body. err1=%s, err2=%s", err.Error(), err2.Error())
			return nil, err2
		}
		log.Printf("issue.create returned error response %v", string(b))
		return nil, err
	}

	ticketURL := fmt.Sprintf("%s/browse/%s", base, issue.Key)
	log.Println("created: " + ticketURL)
	return &ticketURL, nil
}

func (myjira *MyJira) neoFindUserByEmail(jiraClient *jira.Client, email string) ([]jira.User, *jira.Response, error) {
	var queryString = "query=" + email

	apiEndpoint := fmt.Sprintf("/rest/api/2/user/search?%s", queryString[:len(queryString)-1])
	req, err := jiraClient.NewRequest("GET", apiEndpoint, nil)

	if err != nil {
		return nil, nil, err
	}

	users := []jira.User{}
	resp, err := jiraClient.Do(req, &users)

	if err != nil {
		return nil, resp, jira.NewJiraError(resp, err)
	}

	return users, resp, nil
}
