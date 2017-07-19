package manage

import (
	"net/http"

	"github.com/thepatrick/lunch/support"
)

func (app App) failed(w http.ResponseWriter) {
	data := struct {
		Title     string
		LogoutURL string
	}{
		Title:     "Manage Lunch Bot",
		LogoutURL: app.root + "/api/logout",
	}
	support.Render(w, "500.html", data)
}
