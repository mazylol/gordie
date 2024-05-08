package main

import (
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

func main() {
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

	c, _, err := websocket.DefaultDialer.Dial(gatewayUrl.Url, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		var firstMessage = true

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s, type: %d", message, mt)

			if firstMessage {
				var hello HelloEvent
				json.Unmarshal(message, &hello)

				firstMessage = false

				go func() {
					for {
						time.Sleep(time.Millisecond * time.Duration(hello.D.HeartbeatInterval))

						c.WriteJSON(HeartBeat{
							Op: 1,
						})
					}
				}()

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

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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