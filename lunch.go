package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/thepatrick/lunch/support"
	"goji.io"
	"goji.io/pat"
)

func main() {
	fmt.Printf("Let's get lunch!\n")

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// session.SetMode(mgo.Monotonic, true)

	mux := goji.NewMux()

	mux.Handle(pat.New("/places/*"), newPlacesMux(session))

	mux.Handle(pat.New("/slack"), handleSlack(session))

	mux.HandleFunc(pat.New("/"), homepage)

	mux.Use(support.Logging)
	http.ListenAndServe("localhost:8080", mux)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	support.Render(w, "index.html")
}
