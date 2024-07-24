/*
Gordie is a simple wrapper over the Discord API. It aims to be fully complete but also easy to use.

# Ping pong example:

	func main() {
		client := gordie.Client {
			Token: "your-token-here",
			Intents: 14023
		}

		client.AddHandler("MESSAGE_CREATE", func(e *gordie.Event) {
			if e.D.Content == "!ping" {
				client.SendMessage(e.D.ChannelId, "Pong!")
			}
		})

		client.Start()
	}
*/
package gordie

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

type HeartBeat struct {
	Op int `json:"op"`
	D  int `json:"d"`
}

type Client struct {
	Token         string
	ApplicationId string
	Intents       int

	ws   *websocket.Conn
	http *http.Client

	handlers map[string]func(e *Event)
}

// Add an event handler
func (c *Client) AddHandler(eventType string, handler func(e *Event)) {
	if c.handlers == nil {
		c.handlers = make(map[string]func(e *Event))
	}

	c.handlers[eventType] = handler
}

func (c Client) SendMessage(channelId string, content string) {
	payload := map[string]interface{}{
		"content": content,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("error while marshalling payload:", err)
		return
	}

	payloadBuffer := bytes.NewBuffer(payloadBytes)

	req, err := http.NewRequest("POST", "https://discord.com/api/v10/channels/"+channelId+"/messages", payloadBuffer)
	if err != nil {
		log.Println("error while creating request:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bot "+c.Token)

	_, err = c.http.Do(req)
	if err != nil {
		log.Println("error while making POST request:", err)
		return
	}
}

func (c Client) RegisterGuildCommand(guildId string, command SlashCommand) {
	payload := map[string]interface{}{
		"name":        command.Name,
		"description": command.Description,
		"options":     command.Options,
		"type":        1,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("error while marshalling payload:", err)
		return
	}

	payloadBuffer := bytes.NewBuffer(payloadBytes)

	req, err := http.NewRequest("POST", "https://discord.com/api/v10/applications/"+c.ApplicationId+"/guilds/"+guildId+"/commands", payloadBuffer)
	if err != nil {
		log.Println("error while creating request:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bot "+c.Token)

	_, err = c.http.Do(req)
	if err != nil {
		log.Println("error while making POST request:", err)
		return
	}
}

func (c *Client) Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c.http = &http.Client{}

	gatewayUrl := GetGatewayUrl()

	log.Printf("Connecting to %s", gatewayUrl)

	var err error

	c.ws, _, err = websocket.DefaultDialer.Dial(gatewayUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	defer c.ws.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		var firstMessage = true

		for {
			_, message, err := c.ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var eventRaw EventRaw
			json.Unmarshal(message, &eventRaw)

			event := eventRaw.ToEvent()

			if event.T == "MESSAGE_CREATE" {
				if handler, ok := c.handlers[event.T]; ok {
					handler(&event)
				}
			}

			if event.T == "READY" {
				log.Printf("Logged in as %s#%s", event.User.Username, event.User.Discriminator)
				if handler, ok := c.handlers[event.T]; ok {
					handler(&event)
				}
			}

			if event.T == "INTERACTION_CREATE" {
				if handler, ok := c.handlers[event.T]; ok {
					handler(&event)
				}
			}

			if firstMessage {
				// set up heartbeats
				var helloEventRaw HelloEventRaw
				json.Unmarshal(message, &helloEventRaw)

				hello := helloEventRaw.ToHelloEvent()

				firstMessage = false

				go func() {
					for {
						time.Sleep(time.Millisecond * time.Duration(hello.HeartBeatInterval))

						c.ws.WriteJSON(HeartBeat{
							Op: 1,
						})
					}
				}()

				// identify
				c.ws.WriteJSON(map[string]interface{}{
					"op": 2,
					"d": map[string]interface{}{
						"token":   c.Token,
						"intents": c.Intents,
						"properties": map[string]string{
							"os":      runtime.GOOS,
							"browser": "gordie",
							"device":  "gordie",
						},
					},
				})

			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
