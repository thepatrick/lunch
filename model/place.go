package model

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Place is a "Place" to have lunch
type Place struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	TeamID      string        `json:"team_id"`
	ChannelID		string				`json:"channel_id"`
	ChannelName string				`json:"channel_name"`
	Name        string        `json:"name"`
	LastVisited time.Time     `json:"last_visited"`
	LastSkipped time.Time     `json:"last_skipped"`
	SkipCount   uint          `json:"skip_count"`
	VisitCount  uint          `json:"visit_count"`
}

// MarshalJSON Convert place to JSON document
func (place *Place) MarshalJSON() ([]byte, error) {
	var lastVisited string
	if !place.LastVisited.IsZero() {
		lastVisited = place.LastVisited.Format(time.RFC3339Nano)
	}
	var lastSkipped string
	if !place.LastSkipped.IsZero() {
		lastSkipped = place.LastSkipped.Format(time.RFC3339Nano)
	}
	type Alias Place
	return json.Marshal(&struct {
		*Alias
		LastVisited string `json:"last_visited,omitempty"`
		LastSkipped string `json:"last_skipped,omitempty"`
	}{
		Alias:       (*Alias)(place),
		LastVisited: lastVisited,
		LastSkipped: lastSkipped,
	})
}
