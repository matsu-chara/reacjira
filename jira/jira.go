package jira

import (
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/andygrunwald/go-jira"
	"github.com/davecgh/go-spew/spew"
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
		logError(res, err, "user.find")

		return nil, fmt.Errorf("user.find error: %w", err)
	}
	if len(reporters) != 1 {
		return nil, fmt.Errorf("user.find error. found multiple users or zero users on email address search. len=%d", len(reporters))
	}

	epicId := ""
	if epicKey != "" {
		issue, res, err := jiraClient.Issue.Get(epicKey, nil)
		if err != nil {
			logError(res, err, "issuge.get")

			return nil, fmt.Errorf("epic issue.get error: %w", err)
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

	log.Printf("calling api to create a jira issue %+v", spew.Sdump(i))
	issue, res, err := jiraClient.Issue.Create(&i)
	if err != nil {
		logError(res, err, "issue.create")

		return nil, fmt.Errorf("issue.create returned error: %w", err)
	}

	ticketURL := fmt.Sprintf("%s/browse/%s", base, issue.Key)
	log.Println("created: " + ticketURL)
	return &ticketURL, nil
}

// 以下は github.com/andygrunwald/go-jira の user.go をコピペしつつ改変したもの
// GDPR対応用のAPIに対応してなかったので拡張している
func (myjira *MyJira) neoFindUserByEmail(jiraClient *jira.Client, email string) ([]jira.User, *jira.Response, error) {
	u, _ := url.Parse("/rest/api/2/user/search")
	queryString := u.Query()
	queryString.Set("query", email)
	u.RawQuery = queryString.Encode()

	req, err := jiraClient.NewRequest("GET", u.String(), nil)

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

func logError(res *jira.Response, err error, apiName string) {
	if res == nil {
		log.Printf("%s returned error. err=%s\n", apiName, err.Error())

		return
	}

	b, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		log.Printf("%s returned error, failed to parse responseBody. err1=%s, err2=%s\n", apiName, err.Error(), err2.Error())

		return
	}

	log.Printf("%s returned error err1=%s, responseBody=%s", apiName, err.Error(), string(b))
}
