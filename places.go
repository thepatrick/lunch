package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"goji.io"
	"goji.io/pat"

	"github.com/thepatrick/lunch/support"
)

// Place is a "Place" to have lunch
type Place struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name"`
	LastVisited time.Time     `json:"last_visited"`
	LastSkipped time.Time     `json:"last_skipped"`
}

func newPlacesMux(session *mgo.Session) *goji.Mux {
	mux := goji.SubMux()

	ensurePlacesIndex(session)

	mux.HandleFunc(pat.Get("/"), allPlacesHTTP(session))
	mux.HandleFunc(pat.Get("/propose"), proposePlaceHTTP(session))
	mux.HandleFunc(pat.Post("/"), addPlaceHTTP(session))
	mux.HandleFunc(pat.Get("/:id"), placeByID(session))
	mux.HandleFunc(pat.Put("/:id"), updatePlace(session))
	mux.HandleFunc(pat.Post("/:id/visit"), visitPlaceHTTP(session))
	mux.HandleFunc(pat.Post("/:id/skip"), skipPlaceHTTP(session))
	mux.HandleFunc(pat.Delete("/:id"), deletePlace(session))
	mux.Use(support.Logging)
	return mux
}

func ensurePlacesIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("lunch").C("places")

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}
}

func allPlaces(s *mgo.Session) ([]Place, error) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("lunch").C("places")

	var places []Place
	err := c.Find(bson.M{}).All(&places)
	if err != nil {
		log.Println("Failed to get all books: ", err)
		return nil, fmt.Errorf("Database error")
	}

	return places, nil
}

func onlyProposablePlaces(vs []Place) []Place {
	vsf := make([]Place, 0)
	for _, v := range vs {
		if v.LastSkipped.Before(time.Now().Add(-time.Hour*6)) &&
			v.LastVisited.Before(time.Now().Add(-time.Hour*72)) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func sortProposablePlaces(vs []Place) []Place {
	// TODO: weight places so we prefer places we haven't been to / skipped recently
	for i := range vs {
		j := rand.Intn(i + 1)
		vs[i], vs[j] = vs[j], vs[i]
	}
	return vs
}

func proposePlace(s *mgo.Session) (Place, error) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("lunch").C("places")

	var places []Place
	err := c.Find(bson.M{}).All(&places)
	if err != nil {
		log.Println("Failed to get all places: ", err)
		return Place{}, fmt.Errorf("Database error")
	}

	places = onlyProposablePlaces(places)

	if len(places) == 0 {
		return Place{}, fmt.Errorf("There are no places that haven't been skipped or visited recently")
	}

	places = sortProposablePlaces(places)

	return places[0], nil
}

func addPlace(place Place, s *mgo.Session) (string, error) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("lunch").C("places")

	err := c.Insert(place)
	if err != nil {
		if mgo.IsDup(err) {
			return "", fmt.Errorf("A place with this name already exists")
		}

		log.Println("Failed to insert place: ", err)
		return "", fmt.Errorf("Database error")
	}

	return place.ID.Hex(), nil
}

func placeByID(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")
		c := session.DB("lunch").C("places")

		var place Place
		err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&place)
		if err != nil {
			support.ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Printf("Failed to find place %v error: %v\n", id, err)
			return
		}

		if place.ID == "" {
			support.ErrorWithJSON(w, "Place not found", http.StatusNotFound)
		}

		respBody, err := json.MarshalIndent(place, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		support.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func visitPlace(id string, s *mgo.Session) error {
	session := s.Copy()
	defer session.Close()

	c := session.DB("lunch").C("places")

	var place Place
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&place)
	if err != nil {
		log.Printf("Failed to find place %v error: %v\n", id, err)
		return fmt.Errorf("Database error")
	}

	place.LastVisited = time.Now()

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, &place)
	if err != nil {
		if mgo.IsDup(err) {
			return fmt.Errorf("A place with that name already exists")
		}
		switch err {
		case mgo.ErrNotFound:
			return fmt.Errorf("Place not found")
		default:
			log.Println("Failed to update place: ", err)
			return fmt.Errorf("Database error")
		}
	}

	return nil
}

func skipPlace(id string, s *mgo.Session) error {
	session := s.Copy()
	defer session.Close()

	c := session.DB("lunch").C("places")

	var place Place
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&place)
	if err != nil {
		log.Printf("Failed to find place %v error: %v\n", id, err)
		return fmt.Errorf("Database error")
	}

	place.LastSkipped = time.Now()

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, &place)
	if err != nil {
		if mgo.IsDup(err) {
			return fmt.Errorf("A place with this name already exists")
		}
		switch err {
		case mgo.ErrNotFound:
			return fmt.Errorf("Place not found")
		default:
			log.Println("Failed to update place: ", err)
			return fmt.Errorf("Database error")
		}
	}

	return nil
}

func updatePlace(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")

		var place Place
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&place)
		if err != nil {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		if place.Name == "" {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB("lunch").C("places")

		err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, &place)
		if err != nil {
			if mgo.IsDup(err) {
				support.ErrorWithJSON(w, "A place with this name already exists", http.StatusBadRequest)
				return
			}
			switch err {
			case mgo.ErrNotFound:
				support.ErrorWithJSON(w, "Place not found", http.StatusNotFound)
				return
			default:
				support.ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed to update place: ", err)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deletePlace(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")

		c := session.DB("lunch").C("places")

		query := bson.M{"_id": bson.ObjectIdHex(id)}
		log.Println("Looking for ", query)
		err := c.Remove(query)
		if err != nil {
			switch err {
			case mgo.ErrNotFound:
				support.ErrorWithJSON(w, "Place not found", http.StatusNotFound)
				return
			default:
				support.ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed to delete place: ", err)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// HTTP API functions

func allPlacesHTTP(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		places, err := allPlaces(s)

		if err != nil {
			support.ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed to get all books: ", err)
			return
		}

		respBody, err := json.MarshalIndent(places, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		support.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func proposePlaceHTTP(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		place, err := proposePlace(s)

		if err != nil {
			support.ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			return
		}

		respBody, err := json.MarshalIndent(place, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		support.ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func visitPlaceHTTP(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := pat.Param(r, "id")

		err := visitPlace(id, s)
		if err != nil {
			support.ErrorWithJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func skipPlaceHTTP(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")

		err := skipPlace(id, s)
		if err != nil {
			support.ErrorWithJSON(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to find place %v error: %v\n", id, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func addPlaceHTTP(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var place Place
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&place)
		if err != nil {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		placeID, err := addPlace(place, s)

		if err != nil {
			support.ErrorWithJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+string(placeID))
		w.WriteHeader(http.StatusCreated)
	}
}
