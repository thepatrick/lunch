package manage

import (
	goji "goji.io"
	"goji.io/pat"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"
)

// App is the application that provides the API for the management GUI
type App struct {
	root   string
	config model.LunchConfig
	places *model.Places
	store  *sessions.CookieStore
}

//

// NewApp creates a new instance of slackbot.App
func NewApp(root string, config model.LunchConfig, places *model.Places, store *sessions.CookieStore) App {
	return App{root, config, places, store}
}

// NewMux creates a new goji SubMux suitable for mounting this app
func (app App) NewMux() *goji.Mux {
	mux := goji.SubMux()

	mux.HandleFunc(pat.Get("/logout"), app.logout())
	mux.HandleFunc(pat.Get("/redirect"), app.slackRedirect())
	mux.HandleFunc(pat.Get("/whoami"), app.whoami())
	mux.HandleFunc(pat.Post("/places/:id"), app.placeUpdate())
	mux.HandleFunc(pat.Delete("/places/:id"), app.deletePlace())
	mux.HandleFunc(pat.Get("/places"), app.placesAll())
	mux.HandleFunc(pat.Get("/login"), app.login())

	mux.HandleFunc(pat.Get("/graphql"), app.placesGraphql())

	mux.HandleFunc(pat.Delete("/places/:id"), app.deletePlace())

	mux.Use(support.Logging)
	return mux
}
