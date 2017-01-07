package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/thepatrick/lunch/support"
	"goji.io"
	"goji.io/pat"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

// LunchConfig encapsulate all of the config needed to run this Slack App
type LunchConfig struct {
	ClientID     string
	ClientSecret string
	MongoURL     string
	Hostname     string
	Port         int
}

func main() {
	fmt.Printf("Let's get lunch!\n")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(fmt.Errorf("Invalid port: %v", err))
	}

	config := LunchConfig{
		ClientID:     os.Getenv("LUNCH_CLIENT_ID"),
		ClientSecret: os.Getenv("LUNCH_CLIENT_SECRET"),
		MongoURL:     os.Getenv("LUNCH_MONGO_URL"),
		Hostname:     os.Getenv("LUNCH_HOSTNAME"),
		Port:         port,
	}

	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)

	ensurePlacesIndex(session)

	mux := goji.NewMux()

	mux.Handle(pat.New("/slack/*"), newSlackMux(config, session))
	mux.Handle(pat.New("/install/*"), newInstallMux(config))

	// mux.HandleFunc(pat.Get("/"), allPlacesHTTP(session))
	// mux.HandleFunc(pat.Get("/:id"), placeByID(session))
	// mux.HandleFunc(pat.Put("/:id"), updatePlace(session))
	// mux.HandleFunc(pat.Delete("/:id"), deletePlace(session))

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
