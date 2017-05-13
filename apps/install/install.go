package install

import (
	"log"
	"net/http"

	"goji.io/pat"

	goji "goji.io"

	"github.com/google/go-querystring/query"
	"github.com/parnurzeal/gorequest"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"
)

type oauthSuccess struct {
	Ok          bool   `json:"ok"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TeamID      string `json:"team_id"`
}

// App is a simple app for installing the lunch bot
type App struct {
	config model.LunchConfig
	store  *sessions.CookieStore
}

// NewInstallApp creates an instance of InstallApp
func NewInstallApp(config model.LunchConfig, store *sessions.CookieStore) App {
	return App{config, store}
}

// NewMux creats a Goji SubMux router to mount the installer app
func (app App) NewMux() *goji.Mux {
	mux := goji.SubMux()

	mux.HandleFunc(pat.Get("/redirect"), app.slackRedirect())
	mux.HandleFunc(pat.Get("/logout"), app.logout())
	mux.HandleFunc(pat.Get("/"), app.homepage())

	mux.Use(support.Logging)
	return mux
}

func (app App) homepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := app.store.Get(r, "install-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get session: %v\n", err)
			installFailed(w)
			return
		}

		if session.Values["access_token"] == nil || session.Values["team_id"] == nil {
			addToSlackPage(w, app.config)
			return
		}

		loggedinUserPage(w, session)
	}
}

func (app App) slackRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := app.store.Get(r, "install-session")
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
			app.config.ClientID,
			app.config.ClientSecret,
			code,
			app.config.Hostname + "/install/redirect",
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
	}
}

func (app App) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.

		session, err := app.store.Get(r, "install-session")
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
}
