package slackbot

func errorResponse(response string) SlackResponse {
	var attachments []SlackAttachment
	return SlackResponse{"in_channel", response, attachments}
}
