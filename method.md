# Lunch system

1. Create lunch places by POSTing them to /places/
2. Propose a lunch place by GETing /places/propose
3. If lunch place accpetable, POST /places/:id/visit
4. If lunch place not desirable, POST /places/:id/skip, then GET /places/propose again
   (will blacklist this venue for one day)

# Docs I need

mongo: https://godoc.org/gopkg.in/mgo.v2
goji docs: https://godoc.org/goji.io#Mux.Use
Slack:
* https://api.slack.com/slash-commands
* https://api.slack.com/docs/message-buttons