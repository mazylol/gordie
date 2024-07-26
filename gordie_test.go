package gordie_test

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/mazylol/gordie"
)

func loadDotEnv() map[string]string {
	readFile, err := os.Open(".env")

	if err != nil {
		log.Fatalln(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	varMap := make(map[string]string)

	for _, line := range fileLines {
		splitted := strings.Split(line, "=")

		varMap[splitted[0]] = splitted[1]
	}

	return varMap
}

func TestConnect(t *testing.T) {
	vars := loadDotEnv()

	client := gordie.Client{
		Token:         vars["TOKEN"],
		ApplicationId: vars["APPLICATION_ID"],
		Intents:       131071,
	}

	client.AddHandler("READY", func(e *gordie.Event) {
		client.RegisterGuildCommand(vars["GUILD_ID"], gordie.SlashCommand{
			Name:        "pingo",
			Description: "Ping pong from gordie lib!",
			Options:     nil,
		})

		log.Println("Commands registered!")
	})

	client.AddHandler("MESSAGE_CREATE", func(e *gordie.Event) {
		if e.Content == "!ping" {
			client.SendMessage(e.ChannelId, "Pong!")
		}
	})

	client.AddHandler("INTERACTION_CREATE", func(e *gordie.Event) {
		if e.Data.Name == "pingo" {
			log.Println("Interaction received!")
			client.SendInteractionResponse("Pong!", e.Data.Id, e.Token)

		}
	})

	client.Start()
}
