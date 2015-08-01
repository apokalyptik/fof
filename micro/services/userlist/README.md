# Microservice: User List

This microservice is a simple cachign proxy for the Slack user list.
You simply supply it with a slack authentication token via the command
line switch `-st="..."` and git it an http address:port to listen on via
the command lins switch `-listen="..."`

It will fetch the user list from slack, and serve it at 
`http://address:port/users.json`

It is recommended that you only listen on localhost (`127.0.0.1`) because
there is no authentication mechanism

The user list will attempt an update once per 15 minutes, changes will
reflect immediately in users.json
