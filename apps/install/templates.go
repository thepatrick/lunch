package install

import (
	"log"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/gorilla/sessions"
	"github.com/nlopes/slack"

	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"
)

func addToSlackPage(w http.ResponseWriter, config model.LunchConfig) {

	urlData := struct {
		Scope    string `url:"scope"`
		ClientID string `url:"client_id"`
		Redirect string `url:"redirect_uri"`
	}{
		"commands",
		config.ClientID,
		config.Hostname + "/install/redirect",
	}
	v, _ := query.Values(urlData)

	data := struct {
		Title  string
		AddURL string
	}{
		Title:  "Install Lunch Bot",
		AddURL: "https://slack.com/oauth/authorize?" + v.Encode(),
	}

	support.Render(w, "install/prompt.html", data)
}

func installFailed(w http.ResponseWriter) {

	data := struct {
		Title     string
		LogoutURL string
	}{
		Title:     "Install Lunch Bot",
		LogoutURL: "/install/logout",
	}
	support.Render(w, "500.html", data)
}

func loggedinUserPage(w http.ResponseWriter, session *sessions.Session) {
	accessToken := session.Values["access_token"].(string)

	api := slack.New(accessToken)
	api.SetDebug(true)

	auth, err := api.AuthTest()

	if err != nil {
		log.Printf("Failed to get user identity: %v\n", err)
		installFailed(w)
		return
	}

	log.Printf("Got user identity!: %v\n", auth.User)

	support.Render(w, "install/installed.html", nil)
	return
}
