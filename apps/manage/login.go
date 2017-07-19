package manage

import (
	"log"
	"net/http"
)

func (app App) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// session.Values["back"] = r.URL.Query().Get("back")
		// session.Save(r, w)
		authorizeURL := app.slackAuthorizeURL("identity.basic,identity.team", app.root+"api/redirect")
		log.Printf("Authorize URL: %v\n", authorizeURL)
		http.Redirect(w, r, authorizeURL, http.StatusFound)
	}
}
