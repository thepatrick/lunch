package slackbot

// SlackCommand is a command sent to slack
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
	TriggerID   string `param:"trigger_id"`
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
	User       slackUser               `json:"user"`
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

type slackUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// type slackChannel {
// 	ID string `json:"id"`
// 	Name string `json:"name"`
// }
