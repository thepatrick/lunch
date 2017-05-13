package model

// LunchConfig encapsulate all of the config needed to run this Slack App
type LunchConfig struct {
	ClientID     string
	ClientSecret string
	MongoURL     string
	DatabaseName string
	Hostname     string
	Port         int
}
