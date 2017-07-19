package slackbot

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/thepatrick/lunch/support"
)

func (app App) action() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		rawPayload := r.Form.Get("payload")

		var payload slackActionPayload

		decoder := json.NewDecoder(strings.NewReader(rawPayload))
		err = decoder.Decode(&payload)
		if err != nil {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		if len(payload.Actions) != 1 {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		action := payload.Actions[0].Value

		var response SlackResponse

		if action == "skip" {
			response = app.skipPlace(payload.Team.ID, payload.CallbackID)
		} else if action == "ok" {
			response = app.okPlace(payload.Team.ID, payload.CallbackID)
		} else {
			response = app.getHelp(SlackCommand{})
		}

		respBody, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		support.ResponseWithJSON(w, respBody, 200)
	}
}

func (app App) skipPlace(teamID string, placeID string) SlackResponse {
	err := app.places.SkipPlace(teamID, placeID)

	if err != nil {
		return errorResponse(err.Error())
	}

	place, err := app.places.ProposePlace(teamID)
	if err != nil {
		return errorResponse(err.Error())
	}

	message := "Ok, how about " + place.Name + "?"
	attachments := []SlackAttachment{attachmentForPlace(place)}

	return SlackResponse{"in_channel", message, attachments}
}

func (app App) okPlace(teamID string, placeID string) SlackResponse {
	place, err := app.places.VisitPlace(teamID, placeID)

	if err != nil {
		return errorResponse(err.Error())
	}

	message := "Ok, we're going to " + place.Name + " today!"
	return SlackResponse{"in_channel", message, []SlackAttachment{}}
}
