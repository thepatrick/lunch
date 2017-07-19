package manage

import (
	"log"
	"net/http"

	"github.com/thepatrick/lunch/support"
)

func (app App) slackRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := app.store.Get(r, "places-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to get session: %v\n", err)
			app.failed(w)
			return
		}

		code := r.URL.Query().Get("code")

		accessToken, err := support.GetAccessToken(app.config.ClientID, app.config.ClientSecret, app.config.Hostname+app.root+"api/redirect", code)

		if err != nil {
			app.failed(w)
			log.Printf("Failed to get slack oauth.access: %v\n", err)
			return
		}

		session.Values["access_token"] = accessToken.AccessToken
		session.Values["team_id"] = accessToken.TeamID
		session.Save(r, w)

		log.Printf("Created session with access_token %v and team_id %v\n", accessToken.AccessToken, accessToken.TeamID)

		// session.Values["back"] default to /places/
		http.Redirect(w, r, app.root, http.StatusFound)
	}
}
