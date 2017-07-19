package slackbot

import (
	"log"

	"github.com/thepatrick/lunch/model"
)

func attachmentForPlace(place model.Place) SlackAttachment {
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
