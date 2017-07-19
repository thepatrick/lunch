package manage

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nlopes/slack"
	"github.com/thepatrick/lunch/support"
)

type handlerFuncWithSession func(w http.ResponseWriter, r *http.Request, session validSession)

type validSession struct {
	session *sessions.Session
	api     *slack.Client
	user    *slack.UserIdentityResponse
}

func (app App) withValidSession(ok handlerFuncWithSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "places-session")

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
