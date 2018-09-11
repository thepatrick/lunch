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
			response = app.skipPlace(payload.User.Name, payload.Team.ID, payload.Channel.ID, payload.CallbackID)
		} else if action == "ok" {
			response = app.okPlace(payload.User.Name, payload.Team.ID, payload.Channel.ID, payload.CallbackID)
		} else if action == "cancel" {
			response = app.cancelPlace(payload.User.Name)
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

func (app App) cancelPlace(userName string) SlackResponse {
	message := ":wave: Ok, " + userName + "! Cancel acknowledged."
	attachments := []SlackAttachment{}

	return SlackResponse{"in_channel", message, attachments}
}

func (app App) skipPlace(userName string, teamID string, channelID string, placeID string) SlackResponse {
	err := app.places.SkipPlace(teamID, channelID, placeID)

	if err != nil {
		return errorResponse(err.Error())
	}

	previousPlace, err := app.places.FindByID(teamID, channelID, placeID)
	if err != nil {
		return errorResponse(err.Error())
	}

	if err != nil {
		return errorResponse(err.Error())
	}

	message := userName + " said no to " + previousPlace.Name + " :cry:, how about *" + place.Name + "*?"
	attachments := []SlackAttachment{attachmentForPlace(place)}

	return SlackResponse{"in_channel", message, attachments}
}

func (app App) okPlace(userName string, teamID string, channelID string, placeID string) SlackResponse {
	place, err := app.places.VisitPlace(teamID, channelID, placeID)

	if err != nil {
		return errorResponse(err.Error())
	}

	message := ":tada: " + userName + ", let's go to *" + place.Name + "* today!"
	return SlackResponse{"in_channel", message, []SlackAttachment{}}
}
