package manage

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thepatrick/lunch/support"
)

func (app App) whoami() http.HandlerFunc {
	return app.withValidSession(func(w http.ResponseWriter, r *http.Request, session validSession) {
		respBody, err := json.MarshalIndent(session.user, "", "  ")
		if err != nil {
			log.Fatal(err)
			support.ErrorWithJSON(w, "Failed to generate JSON", http.StatusInternalServerError)
			return
		}

		support.ResponseWithJSON(w, respBody, http.StatusOK)
	})
}
