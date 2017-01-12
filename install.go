package main

import (
	"log"
	"net/http"

	"goji.io/pat"

	goji "goji.io"

	"github.com/google/go-querystring/query"
	"github.com/nlopes/slack"
	"github.com/parnurzeal/gorequest"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/support"
)

func addToSlackPage(w http.ResponseWriter, config LunchConfig) {

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

func installHomepage(config LunchConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := store.Get(r, "install-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get session: %v\n", err)
			installFailed(w)
			return
		}

		if session.Values["access_token"] == nil || session.Values["team_id"] == nil {
			addToSlackPage(w, config)
			return
		}

		loggedinUserPage(w, session)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.

	session, err := store.Get(r, "install-session")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create session: %v\n", err)
		installFailed(w)
		return
	}

	session.Options.MaxAge = -1
	sessions.Save(r, w)

	http.Redirect(w, r, "/install/", http.StatusFound)
}

type oauthSuccess struct {
	Ok          bool   `json:"ok"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TeamID      string `json:"team_id"`
}

func slackRedirect(config LunchConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := store.Get(r, "install-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get session: %v\n", err)
			installFailed(w)
			return
		}

		code := r.URL.Query().Get("code")

		data := struct {
			ClientID     string `url:"client_id"`
			ClientSecret string `url:"client_secret"`
			Code         string `url:"code"`
			Redirect     string `url:"redirect_uri"`
		}{
			config.ClientID,
			config.ClientSecret,
			code,
			config.Hostname + "/install/redirect",
		}
		v, _ := query.Values(data)

		slackURL := "https://slack.com/api/oauth.access?" + v.Encode()

		log.Printf("Slack url... %v\n", slackURL)

		var oauthSuccess oauthSuccess

		resp, _, errs := gorequest.New().Get(slackURL).EndStruct(&oauthSuccess)

		if errs != nil {
			installFailed(w)
			log.Printf("Failed to get slack oauth.access: %v\n", errs)
			return
		}

		if resp.StatusCode != 200 {
			log.Printf("Failed to get slack oauth.success: %v\n", resp.StatusCode)
			installFailed(w)
			return
		}

		if !oauthSuccess.Ok {
			log.Printf("Failed to get slack oauth.success: ok is false\n")
			installFailed(w)
			return
		}

		session.Values["access_token"] = oauthSuccess.AccessToken
		session.Values["team_id"] = oauthSuccess.TeamID
		session.Save(r, w)

		log.Printf("Created session with access_token %v and team_id %v\n", oauthSuccess.AccessToken, oauthSuccess.TeamID)

		http.Redirect(w, r, "/install/", http.StatusFound)

		// GET https://slack.com/api/oauth.access?client_id=CLIENT_ID&client_secret=CLIENT_SECRET&code=XXYYZZ

		// {
		// 	"ok": true,
		// 	"access_token": "xoxp-1111827399-16111519414-20367011469-5f89a31i07",
		// 	"scope": "identity.basic",
		// 	"team_id": "T0G9PQBBK"
		// }
	}
}

func newInstallMux(config LunchConfig) *goji.Mux {
	mux := goji.SubMux()

	mux.HandleFunc(pat.Get("/redirect"), slackRedirect(config))
	mux.HandleFunc(pat.Get("/logout"), logout)
	mux.HandleFunc(pat.Get("/"), installHomepage(config))

	mux.Use(support.Logging)
	return mux
}
