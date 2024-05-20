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
		Token:   vars["TOKEN"],
		Intents: 131071,
	}

	client.AddHandler("MESSAGE_CREATE", func(e *gordie.Event) {
		if e.Content == "!ping" {
			client.SendMessage(e.ChannelId, "Pong!")
		}
	})

	client.Start()
}
