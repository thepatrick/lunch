package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/apps/install"
	"github.com/thepatrick/lunch/apps/manage"
	"github.com/thepatrick/lunch/apps/slackbot"
	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"
	"goji.io"
	"goji.io/pat"
)

var store *sessions.CookieStore

func main() {
	fmt.Printf("Let's get lunch!\n")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(fmt.Errorf("Invalid port: %v", err))
	}

	config := model.LunchConfig{
		ClientID:     os.Getenv("LUNCH_CLIENT_ID"),
		ClientSecret: os.Getenv("LUNCH_CLIENT_SECRET"),
		MongoURL:     os.Getenv("LUNCH_MONGO_URL"),
		DatabaseName: os.Getenv("LUNCH_MONGO_DB"),
		Hostname:     os.Getenv("LUNCH_HOSTNAME"),
		Port:         port,
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("LUNCH_SESSION_SECRET")))

	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)

	places := model.NewPlaces(session, config.DatabaseName)

	mux := goji.NewMux()

	slackbotApp := slackbot.NewApp(config, places, store)
	mux.Handle(pat.New("/slack/*"), slackbotApp.NewMux())

	installApp := install.NewInstallApp(config, store)
	mux.Handle(pat.New("/install/*"), installApp.NewMux())
	mux.Handle(pat.New("/install"), simpleRedirect("/install/"))

	manageApp := manage.NewApp("/manage/", config, places, store)
	mux.Handle(pat.New("/manage/api/*"), manageApp.NewMux())

	mux.Handle(pat.New("/places/*"), simpleRedirect("/manage/"))
	mux.Handle(pat.New("/places"), simpleRedirect("/manage/"))

	mux.Handle(pat.New("/manage/*"), simpleFile("./static/manage/index.html"))
	mux.Handle(pat.New("/manage"), simpleFile("./static/manage/index.html"))

	mux.Handle(pat.New("/static/*"), http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	mux.HandleFunc(pat.New("/"), homepage)

	mux.Use(support.Logging)

	fmt.Printf("Listening on %v\n", config.Port)
	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(config.Port), mux)
	if err != nil {
		panic(err)
	}

}

func homepage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Lunch Bot",
	}

	support.Render(w, "index.html", data)
}

func simpleRedirect(toURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, toURL, http.StatusFound)
	}
}

func simpleFile(fileName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fileName)
	}
}
