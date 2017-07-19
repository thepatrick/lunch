package manage

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	"github.com/thepatrick/lunch/support"
)

func (app App) placeUpdate() http.HandlerFunc {
	return app.withValidSession(func(w http.ResponseWriter, r *http.Request, session validSession) {
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

		err = app.places.UpdatePlace(session.user.Team.ID, id, updateBody)
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
