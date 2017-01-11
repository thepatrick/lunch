package main

import (
	"log"
	"net/http"

	"goji.io/pat"

	goji "goji.io"

	"github.com/google/go-querystring/query"
	"github.com/nlopes/slack"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/support"
)

func managePlacesPage(w http.ResponseWriter, session *sessions.Session, places *Places) {
	accessToken := session.Values["access_token"].(string)

	api := slack.New(accessToken)
	api.SetDebug(true)

	auth, err := api.AuthTest()

	if err != nil {
		log.Printf("Failed to get test auth: %v\n", err)
		support.Render(w, "500.html", nil)
		return
	}

	log.Printf("Got user identity!: %v\n", auth.User)

	user, err := api.GetUserIdentity()

	if err != nil {
		log.Printf("Failed to get user identity: %v\n", err)
		support.Render(w, "500.html", nil)
		return
	}

	log.Printf("Got user identity!: %v\n", user.Team.Name)

	allPlaces, err := places.allPlaces(user.Team.ID)
	if err != nil {
		log.Printf("Failed to get all places: %v\n", err)
		support.Render(w, "500.html", nil)
		return
	}

	data := struct {
		TeamName string
		Places   []Place
	}{
		user.Team.Name,
		allPlaces,
	}

	support.Render(w, "places/index.html", data)
	return
}

func slackAuthorizeURL(config LunchConfig, scope string, redirect string) string {
	urlData := struct {
		Scope    string `url:"scope"`
		ClientID string `url:"client_id"`
		Redirect string `url:"redirect_uri"`
	}{
		scope,
		config.ClientID,
		config.Hostname + redirect,
	}
	v, _ := query.Values(urlData)

	return "https://slack.com/oauth/authorize?" + v.Encode()
}

func manageHomepage(config LunchConfig, places *Places) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := store.Get(r, "places-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get session: %v\n", err)
			support.Render(w, "500.html", nil)
			return
		}

		if session.Values["access_token"] == nil || session.Values["team_id"] == nil {
			authorizeURL := slackAuthorizeURL(config, "identity.basic,identity.team", "/places/redirect")
			log.Printf("Authorize URL: %v\n", authorizeURL)
			http.Redirect(w, r, authorizeURL, http.StatusFound)
			return
		}

		managePlacesPage(w, session, places)
	}
}

func manageLogout(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.

	session, err := store.Get(r, "places-session")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create session: %v\n", err)
		support.Render(w, "500.html", nil)
		return
	}

	session.Options.MaxAge = -1
	sessions.Save(r, w)

	http.Redirect(w, r, "/places/", http.StatusFound)
}

func manageSlackRedirect(config LunchConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := store.Get(r, "places-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get session: %v\n", err)
			support.Render(w, "500.html", nil)
			return
		}

		code := r.URL.Query().Get("code")

		accessToken, err := support.GetAccessToken(config.ClientID, config.ClientSecret, config.Hostname+"/places/redirect", code)

		if err != nil {
			support.Render(w, "500.html", nil)
			log.Printf("Failed to get slack oauth.access: %v\n", err)
			return
		}

		session.Values["access_token"] = accessToken.AccessToken
		session.Values["team_id"] = accessToken.TeamID
		session.Save(r, w)

		log.Printf("Created session with access_token %v and team_id %v\n", accessToken.AccessToken, accessToken.TeamID)

		http.Redirect(w, r, "/places/", http.StatusFound)
	}
}

func newManageMux(config LunchConfig, places *Places) *goji.Mux {
	mux := goji.SubMux()

	mux.HandleFunc(pat.Get("/redirect"), manageSlackRedirect(config))
	mux.HandleFunc(pat.Get("/logout"), manageLogout)
	mux.HandleFunc(pat.Get("/"), manageHomepage(config, places))

	mux.Use(support.Logging)
	return mux
}
