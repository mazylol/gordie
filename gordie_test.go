package gordie_test

import (
	"testing"

	"github.com/mazylol/gordie"
)

func TestConnect(t *testing.T) {
	client := gordie.Client{
		Token:   "33335",
		Intents: 33281,
	}

	client.Start()
}
