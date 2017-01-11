package support

import (
	"log"

	"fmt"

	"github.com/google/go-querystring/query"
	"github.com/parnurzeal/gorequest"
)

// AccessTokenResponse contains an AccessToken from Slack
type AccessTokenResponse struct {
	Ok          bool   `json:"ok"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TeamID      string `json:"team_id"`
}

// GetAccessToken converts an authorization token into an AccessTokenResponse
func GetAccessToken(clientID string, clientSecret string, redirect string, code string) (AccessTokenResponse, error) {
	data := struct {
		ClientID     string `url:"client_id"`
		ClientSecret string `url:"client_secret"`
		Code         string `url:"code"`
		Redirect     string `url:"redirect_uri"`
	}{
		clientID,
		clientSecret,
		code,
		redirect,
	}
	v, _ := query.Values(data)

	slackURL := "https://slack.com/api/oauth.access?" + v.Encode()

	log.Printf("Slack url... %v\n", slackURL)

	var success AccessTokenResponse

	resp, _, errs := gorequest.New().Get(slackURL).EndStruct(&success)

	if errs != nil {
		return AccessTokenResponse{}, fmt.Errorf("Failed to get slack oauth.access: %v", errs)
	}

	if resp.StatusCode != 200 {
		return AccessTokenResponse{}, fmt.Errorf("Failed to get slack oauth.success: %v", resp.StatusCode)
	}

	if !success.Ok {
		return AccessTokenResponse{}, fmt.Errorf("Failed to get slack oauth.success: ok is false")
	}

	return success, nil
}
