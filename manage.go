package main

import (
	"encoding/json"
	"log"
	"net/http"

	"goji.io/pat"

	goji "goji.io"

	"github.com/google/go-querystring/query"
	"github.com/nlopes/slack"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"
)

func slackAuthorizeURL(config model.LunchConfig, scope string, redirect string) string {
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

func manageLogout(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "places-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to create session: %v\n", err)
			manageFailed(root, w)
			return
		}

		session.Options.MaxAge = -1
		sessions.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func manageSlackRedirect(root string, config model.LunchConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := store.Get(r, "places-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get session: %v\n", err)
			manageFailed(root, w)
			return
		}

		code := r.URL.Query().Get("code")

		accessToken, err := support.GetAccessToken(config.ClientID, config.ClientSecret, config.Hostname+root+"api/redirect", code)

		if err != nil {
			manageFailed(root, w)
			log.Printf("Failed to get slack oauth.access: %v\n", err)
			return
		}

		session.Values["access_token"] = accessToken.AccessToken
		session.Values["team_id"] = accessToken.TeamID
		session.Save(r, w)

		log.Printf("Created session with access_token %v and team_id %v\n", accessToken.AccessToken, accessToken.TeamID)

		// session.Values["back"] default to /places/
		http.Redirect(w, r, root, http.StatusFound)
	}
}

func manageFailed(root string, w http.ResponseWriter) {

	data := struct {
		Title     string
		LogoutURL string
	}{
		Title:     "Manage Lunch Bot",
		LogoutURL: root + "/api/logout",
	}
	support.Render(w, "500.html", data)
}

func manageLogin(root string, config model.LunchConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// session.Values["back"] = r.URL.Query().Get("back")
		// session.Save(r, w)
		authorizeURL := slackAuthorizeURL(config, "identity.basic,identity.team", root+"api/redirect")
		log.Printf("Authorize URL: %v\n", authorizeURL)
		http.Redirect(w, r, authorizeURL, http.StatusFound)
	}
}

func managePlacesAll(config model.LunchConfig, places *model.Places) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "places-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to create session: %v\n", err)
			support.ErrorWithJSON(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		if session.Values["access_token"] == nil || session.Values["team_id"] == nil {
			support.ErrorWithJSON(w, "Failed to get user identity", http.StatusUnauthorized)
			return
		}

		accessToken := session.Values["access_token"].(string)

		api := slack.New(accessToken)
		api.SetDebug(true)

		user, err := api.GetUserIdentity()

		if err != nil {
			log.Printf("Failed to get user identity: %v\n", err)
			support.ErrorWithJSON(w, "Failed to get user identity", http.StatusUnauthorized)
			return
		}

		allPlaces, err := places.AllPlaces(user.Team.ID)
		if err != nil {
			log.Printf("Failed to get all places: %v\n", err)
			support.ErrorWithJSON(w, "Failed to get places", http.StatusInternalServerError)
			return
		}

		respBody, err := json.MarshalIndent(allPlaces, "", "  ")
		if err != nil {
			log.Fatal(err)
			support.ErrorWithJSON(w, "Failed to generate JSON", http.StatusInternalServerError)
			return
		}

		support.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

type handlerFuncWithSession func(w http.ResponseWriter, r *http.Request, session validSession)

type validSession struct {
	session *sessions.Session
	api     *slack.Client
	user    *slack.UserIdentityResponse
}

func withValidSession(ok handlerFuncWithSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "places-session")

		if err != nil {
			log.Printf("Failed to create session: %v\n", err)
			support.ErrorWithJSON(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		if session.Values["access_token"] == nil || session.Values["team_id"] == nil {
			support.ErrorWithJSON(w, "Failed to get user identity", http.StatusUnauthorized)
			return
		}

		accessToken := session.Values["access_token"].(string)

		api := slack.New(accessToken)
		api.SetDebug(true)

		user, err := api.GetUserIdentity()
		if err != nil {
			log.Printf("Failed to get user identity: %v\n", err)
			support.ErrorWithJSON(w, "Failed to get user identity", http.StatusUnauthorized)
			return
		}

		ok(w, r, validSession{session, api, user})
	}
}

func manageWhoami(config model.LunchConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "places-session")
		if err != nil {
			log.Printf("Failed to create session: %v\n", err)
			support.ErrorWithJSON(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		if session.Values["access_token"] == nil || session.Values["team_id"] == nil {
			support.ErrorWithJSON(w, "Failed to get user identity", http.StatusUnauthorized)
			return
		}

		accessToken := session.Values["access_token"].(string)

		api := slack.New(accessToken)
		api.SetDebug(true)

		user, err := api.GetUserIdentity()

		if err != nil {
			log.Printf("Failed to get user identity: %v\n", err)
			support.ErrorWithJSON(w, "Failed to get user identity", http.StatusUnauthorized)
			return
		}

		respBody, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			log.Fatal(err)
			support.ErrorWithJSON(w, "Failed to generate JSON", http.StatusInternalServerError)
			return
		}

		support.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func manageUpdatePlace(config model.LunchConfig, places *model.Places) http.HandlerFunc {
	return withValidSession(func(w http.ResponseWriter, r *http.Request, session validSession) {
		id := pat.Param(r, "id")

		updateBody := struct {
			Name string `json:"name"`
		}{""}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&updateBody)
		if err != nil {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		if updateBody.Name == "" {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		err = places.UpdatePlace(session.user.Team.ID, id, updateBody)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err.Error() == "A place with that name already exists" {
				statusCode = http.StatusConflict
			}
			if err.Error() == "Place not found" {
				statusCode = http.StatusNotFound
			}
			support.ErrorWithJSON(w, err.Error(), statusCode)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

type ManageAPI struct {
	root   string
	config model.LunchConfig
	places *model.Places
}

func (manage ManageAPI) deletePlace() http.HandlerFunc {
	return withValidSession(func(w http.ResponseWriter, r *http.Request, session validSession) {
		id := pat.Param(r, "id")

		err := manage.places.DeletePlace(session.user.Team.ID, id)
		if err != nil {
			statusCode := http.StatusBadRequest
			if err.Error() == "Place not found" {
				statusCode = http.StatusNotFound
			}
			support.ErrorWithJSON(w, err.Error(), statusCode)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func newManageMux(root string, config model.LunchConfig, places *model.Places) *goji.Mux {
	mux := goji.SubMux()

	manage := ManageAPI{root, config, places}

	mux.HandleFunc(pat.Get("/logout"), manageLogout(root))
	mux.HandleFunc(pat.Get("/redirect"), manageSlackRedirect(root, config))
	mux.HandleFunc(pat.Get("/whoami"), manageWhoami(config))
	mux.HandleFunc(pat.Post("/places/:id"), manageUpdatePlace(config, places))
	mux.HandleFunc(pat.Get("/places"), managePlacesAll(config, places))
	mux.HandleFunc(pat.Get("/login"), manageLogin(root, config))

	mux.HandleFunc(pat.Delete("/places/:id"), manage.deletePlace())

	mux.Use(support.Logging)
	return mux
}
