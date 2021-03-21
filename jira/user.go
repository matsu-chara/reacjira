package jira

import (
	gojira "github.com/andygrunwald/go-jira"
	"golang.org/x/xerrors"
)

func (myJiraImpl *myJiraImpl) FindUserByEmail(email string) (*gojira.User, error) {
	users, resp, err := myJiraImpl.underlying.User.Find(email)
	if err != nil {
		return nil, xerrors.Errorf("an error occurred: %w", gojira.NewJiraError(resp, err))
	}

	if len(users) == 0 {
		err = xerrors.Errorf("an error occurred. users were not found on the email address search.")
		return nil, err
	} else if len(users) != 1 {
		err = xerrors.Errorf("an error occurred. multiple users were found on the email address search. len=%d", len(users))
		return nil, err
	}
	return &users[0], nil
}
