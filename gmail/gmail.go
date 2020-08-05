package gmail

import (
	"fmt"
	"sync"

	"github.com/levigross/grequests"
)

type messageRes struct {
	Messages []struct {
		ID       string `json:"id"`
		ThreadID string `json:"threadId"`
	} `json:"messages"`
	ResultSizeEstimate int `json:"resultSizeEstimate"`
}

type apiError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message      string `json:"message"`
			Domain       string `json:"domain"`
			Reason       string `json:"reason"`
			Location     string `json:"location"`
			LocationType string `json:"locationType"`
		} `json:"errors"`
		Status string `json:"status"`
	} `json:"error"`
}

type apiToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// Payload is a struct to hold all items we send to the server
type Payload struct {
	User         string
	AccessToken  string
	Wg           *sync.WaitGroup
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func (p *Payload) authenticate() {
	data, err := grequests.Post(
		"https://www.googleapis.com/oauth2/v4/token",
		&grequests.RequestOptions{
			JSON: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     p.ClientID,
				"client_secret": p.ClientSecret,
				"refresh_token": p.RefreshToken,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	var at apiToken
	err = data.JSON(&at)
	if err != nil {
		panic(err)
	}
	p.AccessToken = at.AccessToken
}

// GetUnread returns a list of unread emails
func (p *Payload) GetUnread() int {
	defer p.Wg.Done()
	//
	data := p.getUnreadEmails()

	var e apiError

	if err := data.JSON(&e); err == nil {
		if e.Error.Code == 401 {
			p.authenticate()
			data = p.getUnreadEmails()
		}
	}

	var msg messageRes
	err := data.JSON(&msg)
	if err != nil {
		panic(err)
	}

	return msg.ResultSizeEstimate
}

func (p *Payload) getUnreadEmails() *grequests.Response {
	resp, err := grequests.Get(
		fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/%s/messages?labelIds=UNREAD", p.User),
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", p.AccessToken),
			},
		},
	)
	if err != nil {
		panic(err)
	}
	return resp
}
