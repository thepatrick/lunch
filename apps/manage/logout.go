package manage

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func (app App) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "places-session")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to create session: %v\n", err)
			app.failed(w)
			return
		}

		session.Options.MaxAge = -1
		sessions.Save(r, w)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
