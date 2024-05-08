package gordie

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Gateway struct {
	Url string `json:"url"`
}

type HelloEvent struct {
	Op int `json:"op"`
	D  struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}

type HeartBeat struct {
	Op int `json:"op"`
	D  int `json:"d"`
}

type Client struct {
	Token   string
	Intents int

	ws *websocket.Conn

	handlers map[string]func(e *Event)
}

type Event struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		Content   string `json:"content"`
		GuildId   string `json:"guild_id"`
		ChannelId string `json:"channel_id"`
	} `json:"d"`
}

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

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("error while making POST request:", err)
		return
	}
}

func (c *Client) Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	gatewayUrlRequest, err := http.Get("https://discord.com/api/v10/gateway")
	if err != nil {
		log.Fatalln(err)
		return
	}

	gatewayUrlRequestBody, err := io.ReadAll(gatewayUrlRequest.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}

	var gatewayUrl Gateway
	json.Unmarshal(gatewayUrlRequestBody, &gatewayUrl)

	log.Printf("Connecting to %s", gatewayUrl.Url)

	c.ws, _, err = websocket.DefaultDialer.Dial(gatewayUrl.Url, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer c.ws.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		var firstMessage = true

		for {
			mt, message, err := c.ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s, type: %d", message, mt)

			var event Event
			json.Unmarshal(message, &event)

			if event.T == "MESSAGE_CREATE" {
				if handler, ok := c.handlers[event.T]; ok {
					handler(&event)
				}
			}

			if firstMessage {
				// set up heartbeats
				var hello HelloEvent
				json.Unmarshal(message, &hello)

				firstMessage = false

				go func() {
					for {
						time.Sleep(time.Millisecond * time.Duration(hello.D.HeartbeatInterval))

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
							"os":      "linux",
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
