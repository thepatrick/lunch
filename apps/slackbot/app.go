package slackbot

import (
	goji "goji.io"
	"goji.io/pat"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"
)

// App is the API interface that slack uses to make the lunch bot work
type App struct {
	config model.LunchConfig
	places *model.Places
	store  *sessions.CookieStore
}

// NewApp creates a new instance of slackbot.App
func NewApp(config model.LunchConfig, places *model.Places, store *sessions.CookieStore) App {
	return App{config, places, store}
}

// NewMux creates a new goji SubMux suitable for mounting this app
func (app App) NewMux() *goji.Mux {
	mux := goji.SubMux()

	mux.Handle(pat.New("/action"), app.action())
	mux.Handle(pat.New("/command"), app.command())

	mux.Use(support.Logging)
	return mux
}
