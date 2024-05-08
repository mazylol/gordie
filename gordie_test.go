package gordie_test

import (
	"testing"

	"github.com/mazylol/gordie"
)

func TestConnect(t *testing.T) {
	client := gordie.Client{
		Token:   "33335",
		Intents: 14023,
	}

	client.AddHandler("MESSAGE_CREATE", func(e *gordie.Event) {
		if e.D.Content == "!ping" {
			client.SendMessage(e.D.ChannelId, "Pong!")
		}
	})

	client.Start()
}
