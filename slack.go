package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	goji "goji.io"
	"goji.io/pat"

	"github.com/goji/param"

	"github.com/thepatrick/lunch/support"
)

// SlackCommand is an incoming slack command
type SlackCommand struct {
	Token       string `param:"token"`
	TeamID      string `param:"team_id"`
	TeamDomain  string `param:"team_domain"`
	ChannelID   string `param:"channel_id"`
	ChannelName string `param:"channel_name"`
	UserID      string `param:"user_id"`
	UserName    string `param:"user_name"`
	Command     string `param:"command"`
	Text        string `param:"text"`
	ResponseURL string `param:"response_url"`
}

// SlackResponse is a response to a slack command
type SlackResponse struct {
	ResponseType string            `json:"response_type"`
	Text         string            `json:"text"`
	Attachments  []SlackAttachment `json:"attachments"`
}

// SlackAttachment is an attachment in a response to a slack command
type SlackAttachment struct {
	Title          string                  `json:"title"`
	Text           string                  `json:"text"`
	MarkdownIn     []string                `json:"mrkdwn_in"`
	Fallback       string                  `json:"fallback"`
	Color          string                  `json:"color"`
	AttachmentType string                  `json:"attachment_type"`
	Actions        []SlackAttachmentAction `json:"actions"`
	CallbackID     string                  `json:"callback_id"`
}

// SlackAttachmentAction is an action button on a slack message
type SlackAttachmentAction struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Style string `json:"style"`
	Type  string `json:"type"`
	Value string `json:"value"`
	// Confirm SlackAttachmentConfirm `json:"confirm"`
}

// SlackAttachmentConfirm is a confirmation prompt on a SlackAttachmentAction
type SlackAttachmentConfirm struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	OkText      string `json:"ok_text"`
	DismissText string `json:"dismiss_text"`
}

type slackActionPayload struct {
	Actions    []SlackAttachmentAction `json:"actions"`
	CallbackID string                  `json:"callback_id"`
	Team       slackTeam               `json:"team"`
	// UserID slackUser `json:"user"`
	// ActionTs string `json:"action_ts"`
	// MessageTs string `json:"message_ts"`
	// AttachmentID string `json:"attachment_id"`
	// Token string `json:"token"`
	// OriginalMessage slackMessage `json:"original_message"`
	// .. "text":"How about Somewhere",
	// .. "bot_id":"B3M7BQ0A0",
	// ... the message we sent
	//..  "ts":"1483532908.000012"
	// ResponseURL string `json:"response_url"`
}

type slackTeam struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

// type slackUser {
// 	ID string `json:"id"`
// 	Name string `json:"name"`
// }
// type slackChannel {
// 	ID string `json:"id"`
// 	Name string `json:"name"`
// }

func newSlackMux(config LunchConfig, places *Places) *goji.Mux {
	mux := goji.SubMux()

	mux.Handle(pat.New("/action"), handleSlackAction(config, places))
	mux.Handle(pat.New("/command"), handleSlack(config, places))

	mux.Use(support.Logging)
	return mux
}

func handleSlackAction(config LunchConfig, places *Places) http.HandlerFunc {
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
			response = slackSkipPlace(payload.Team.ID, payload.CallbackID, places)
		} else if action == "ok" {
			response = slackOKPlace(payload.Team.ID, payload.CallbackID, places)
		} else {
			response = slackGetHelp(SlackCommand{})
		}

		respBody, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		support.ResponseWithJSON(w, respBody, 200)
	}
}

func handleSlack(config LunchConfig, places *Places) http.HandlerFunc {
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
			response = slackSuggestPlace(command, places)
		} else if words[0] == "list" {
			response = slackListPlaces(config, command, places)
		} else if words[0] == "add" {
			response = slackAddPlace(command, words[1:], places)
		} else {
			response = slackGetHelp(command)
		}

		respBody, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		support.ResponseWithJSON(w, respBody, 200)
	}
}

func slackErrorResponse(response string) SlackResponse {
	var attachments []SlackAttachment
	return SlackResponse{"in_channel", response, attachments}
}

func slackGetHelp(command SlackCommand) SlackResponse {
	message := "To add a place use: `" + command.Command + " add Awesome Place`" + "\n" +
		"To get a suggestion use: `" + command.Command + "`"

	var attachments []SlackAttachment
	// attachments = append(attachments, SlackAttachment{"Go here", "/lunch visit $PLACE_ID"})

	return SlackResponse{"in_channel", message, attachments}
}

func slackAttachmentForPlace(place Place) SlackAttachment {
	lunchURL := place.ID.Hex()

	log.Println("lunchURL", lunchURL)

	attachment := SlackAttachment{Title: "", Text: ""}
	attachment.MarkdownIn = []string{"text"}

	attachment.CallbackID = place.ID.Hex()
	attachment.Fallback = "You are unable to play this game, sorry."
	attachment.Color = "#3AA3E3"
	attachment.AttachmentType = "default"

	okAction := SlackAttachmentAction{Name: "ok", Style: "primary", Text: "Sounds good", Type: "button", Value: "ok"}
	skipAction := SlackAttachmentAction{Name: "skip", Text: "Not today", Type: "button", Value: "skip"}

	attachment.Actions = []SlackAttachmentAction{okAction, skipAction}

	return attachment
}

func slackSuggestPlace(command SlackCommand, places *Places) SlackResponse {
	place, err := places.proposePlace(command.TeamID)
	if err != nil {
		return slackErrorResponse(err.Error() + " Maybe try adding one using `" + command.Command + " add Awesome Place`")
	}

	message := "How about " + place.Name
	attachments := []SlackAttachment{slackAttachmentForPlace(place)}

	return SlackResponse{"in_channel", message, attachments}
}

func slackSkipPlace(teamID string, placeID string, places *Places) SlackResponse {
	err := places.skipPlace(teamID, placeID)

	if err != nil {
		return slackErrorResponse(err.Error())
	}

	place, err := places.proposePlace(teamID)
	if err != nil {
		return slackErrorResponse(err.Error())
	}

	message := "Ok, how about " + place.Name + "?"
	attachments := []SlackAttachment{slackAttachmentForPlace(place)}

	return SlackResponse{"in_channel", message, attachments}
}

func slackOKPlace(teamID string, placeID string, places *Places) SlackResponse {
	place, err := places.visitPlace(teamID, placeID)

	if err != nil {
		return slackErrorResponse(err.Error())
	}

	message := "Ok, we're going to " + place.Name + " today!"
	return SlackResponse{"in_channel", message, []SlackAttachment{}}
}

func slackListPlaces(config LunchConfig, command SlackCommand, places *Places) SlackResponse {
	// places, err := allPlaces(s)
	// if err != nil {
	// 	return slackErrorResponse(err.Error())
	// }

	// if len(places) == 0 {
	// 	return slackErrorResponse("No places yet, try `/" + command.Command + " add Place Name` to add your first one!")
	// }

	listURL := config.Hostname + "/places"

	return SlackResponse{"ephemeral", "To see the list of places go to " + listURL, []SlackAttachment{}}
}

func slackAddPlace(command SlackCommand, words []string, places *Places) SlackResponse {
	if len(words) == 0 {
		return slackErrorResponse("You forgot the place name!")
	}

	placeName := strings.Join(words, " ")

	var place Place
	place.Name = placeName
	place.TeamID = command.TeamID

	_, err := places.addPlace(place)

	if err != nil {
		return slackErrorResponse(err.Error())
	}

	return SlackResponse{"ephemeral", "I've added " + placeName, []SlackAttachment{}}
}
