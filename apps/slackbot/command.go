package slackbot

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/goji/param"
	"github.com/thepatrick/lunch/model"
	"github.com/thepatrick/lunch/support"
)

func (app App) command() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		var command SlackCommand

		// Parse url.Values (in this case, r.PostForm) and
		// a pointer to our struct so that param can populate it.
		err = param.Parse(r.Form, &command)
		if err != nil {
			support.ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			log.Println("Failed to parse body from slack: ", err)
			return
		}

		var response SlackResponse

		words := strings.Fields(command.Text)

		if len(words) == 0 {
			response = app.suggestPlace(command)
		} else if words[0] == "list" {
			response = app.listPlaces(command)
		} else if words[0] == "add" {
			response = app.addPlace(command, words[1:])
		} else {
			response = app.getHelp(command)
		}

		respBody, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		support.ResponseWithJSON(w, respBody, 200)
	}
}

func (app App) suggestPlace(command SlackCommand) SlackResponse {
	place, err := app.places.ProposePlace(command.TeamID)
	if err != nil {
		return errorResponse(err.Error() + " Maybe try adding one using `" + command.Command + " add Awesome Place`")
	}

	message := ":knife_fork_plate: How about *" + place.Name + "*?"
	attachments := []SlackAttachment{attachmentForPlace(place)}

	return SlackResponse{"in_channel", message, attachments}
}

func (app App) listPlaces(command SlackCommand) SlackResponse {
	listURL := app.config.Hostname + "/places"
	return SlackResponse{"ephemeral", "To see the list of places go to " + listURL, []SlackAttachment{}}
}

func (app App) addPlace(command SlackCommand, words []string) SlackResponse {
	if len(words) == 0 {
		return errorResponse("You forgot the place name!")
	}

	placeName := strings.Join(words, " ")

	var place model.Place
	place.Name = placeName
	place.TeamID = command.TeamID

	_, err := app.places.AddPlace(place)

	if err != nil {
		return errorResponse(err.Error())
	}

	return SlackResponse{"ephemeral", "I've added " + placeName, []SlackAttachment{}}
}

func (app App) getHelp(command SlackCommand) SlackResponse {
	message := "To add a place use: `" + command.Command + " add Awesome Place`" + "\n" +
		"To get a suggestion use: `" + command.Command + "`"

	var attachments []SlackAttachment
	// attachments = append(attachments, SlackAttachment{"Go here", "/lunch visit $PLACE_ID"})

	return SlackResponse{"in_channel", message, attachments}
}
