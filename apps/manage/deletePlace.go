package manage

import (
	"net/http"

	"github.com/thepatrick/lunch/support"
	"goji.io/pat"
)

func (manage App) deletePlace() http.HandlerFunc {
	return manage.withValidSession(func(w http.ResponseWriter, r *http.Request, session validSession) {
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
