package manage

import "github.com/google/go-querystring/query"

func (app App) slackAuthorizeURL(scope string, redirect string) string {
	urlData := struct {
		Scope    string `url:"scope"`
		ClientID string `url:"client_id"`
		Redirect string `url:"redirect_uri"`
	}{
		scope,
		app.config.ClientID,
		app.config.Hostname + redirect,
	}
	v, _ := query.Values(urlData)

	return "https://slack.com/oauth/authorize?" + v.Encode()
}
