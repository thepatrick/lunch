package manage

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thepatrick/lunch/support"
)

func (app App) placesAll() http.HandlerFunc {
	return app.withValidSession(func(w http.ResponseWriter, r *http.Request, session validSession) {
		allPlaces, err := app.places.AllPlaces(session.user.Team.ID)
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
	})
}
