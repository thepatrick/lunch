package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/goji/param"

	mgo "gopkg.in/mgo.v2"

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
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	MarkdownIn []string `json:"mrkdwn_in"`
}

func handleSlack(s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

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
			response = slackSuggestPlace(command, session)
		} else if words[0] == "skip" {
			response = slackSkipPlace(command, words[1:], session)
		} else if words[0] == "ok" {
			response = slackOKPlace(command, words[1:], session)
		} else if words[0] == "list" {
			response = slackListPlaces(command, session)
		} else if words[0] == "add" {
			response = slackAddPlace(command, words[1:], session)
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
	message := "You said: \"" + command.Text + "\" to " + command.Command

	var attachments []SlackAttachment
	// attachments = append(attachments, SlackAttachment{"Go here", "/lunch visit $PLACE_ID"})

	return SlackResponse{"in_channel", message, attachments}
}

func slackAttachmentForPlace(command string, place Place) SlackAttachment {
	lunchURL := place.ID.Hex()

	log.Println("lunchURL", lunchURL)

	loveMessage := "Good for this one? `" + command + " ok " + lunchURL + "` "
	skipMessage := "Don't like this? `" + command + " skip " + lunchURL + "`"

	return SlackAttachment{"", loveMessage + "\n" + skipMessage, []string{"text"}}
}

func slackSuggestPlace(command SlackCommand, s *mgo.Session) SlackResponse {
	place, err := proposePlace(s)
	if err != nil {
		slackErrorResponse(err.Error())
	}

	message := "How about " + place.Name
	attachments := []SlackAttachment{slackAttachmentForPlace(command.Command, place)}

	return SlackResponse{"in_channel", message, attachments}
}

func slackSkipPlace(command SlackCommand, words []string, s *mgo.Session) SlackResponse {
	if len(words) == 0 {
		return slackErrorResponse("You forgot the place ID!")
	}
	if len(words) > 1 {
		return slackErrorResponse("You said too much! I can only skip one place at a time.")
	}

	err := skipPlace(words[0], s)

	if err != nil {
		return slackErrorResponse(err.Error())
	}

	place, err := proposePlace(s)
	if err != nil {
		slackErrorResponse(err.Error())
	}

	message := "Ok, how about " + place.Name + "?"
	attachments := []SlackAttachment{slackAttachmentForPlace(command.Command, place)}

	return SlackResponse{"in_channel", message, attachments}
}

func slackOKPlace(command SlackCommand, words []string, s *mgo.Session) SlackResponse {
	if len(words) == 0 {
		return slackErrorResponse("You forgot the place ID!")
	}
	if len(words) > 1 {
		return slackErrorResponse("Oops, you can only go to one place at a time (for now?)")
	}

	err := visitPlace(words[0], s)

	if err != nil {
		return slackErrorResponse(err.Error())
	}

	message := "Ok, I'll remember."
	return SlackResponse{"in_channel", message, []SlackAttachment{}}
}

func slackListPlaces(command SlackCommand, s *mgo.Session) SlackResponse {
	places, err := allPlaces(s)
	if err != nil {
		return slackErrorResponse(err.Error())
	}

	if len(places) == 0 {
		return slackErrorResponse("No places yet, try `/" + command.Command + " add Place Name` to add your first one!")
	}

	listURL := "http://localhost:8080/places?teamID=" + command.TeamID + "&channelID=" + command.ChannelID + "&userID=" + command.UserID

	return SlackResponse{"ephemeral", "To see the list of places go to [http://localhost:8080/places](" + listURL + ")", []SlackAttachment{}}
}

func slackAddPlace(command SlackCommand, words []string, s *mgo.Session) SlackResponse {
	if len(words) == 0 {
		return slackErrorResponse("You forgot the place name!")
	}

	placeName := strings.Join(words, " ")

	var place Place
	place.Name = placeName

	_, err := addPlace(place, s)

	if err != nil {
		return slackErrorResponse(err.Error())
	}

	return SlackResponse{"ephemeral", "I've added " + placeName, []SlackAttachment{}}
}

// {
//     "text": "Would you like to play a game?",
//     "attachments": [
//         {
//             "text": "Choose a game to play",
//             "fallback": "You are unable to choose a game",
//             "callback_id": "wopr_game",
//             "color": "#3AA3E3",
//             "attachment_type": "default",
//             "actions": [
//                 {
//                     "name": "chess",
//                     "text": "Chess",
//                     "type": "button",
//                     "value": "chess"
//                 },
//                 {
//                     "name": "maze",
//                     "text": "Falken's Maze",
//                     "type": "button",
//                     "value": "maze"
//                 },
//                 {
//                     "name": "war",
//                     "text": "Thermonuclear War",
//                     "style": "danger",
//                     "type": "button",
//                     "value": "war",
//                     "confirm": {
//                         "title": "Are you sure?",
//                         "text": "Wouldn't you prefer a good game of chess?",
//                         "ok_text": "Yes",
//                         "dismiss_text": "No"
//                     }
//                 }
//             ]
//         }
//     ]
// }
