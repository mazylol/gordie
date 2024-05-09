package gordie

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Gateway struct {
	Url string `json:"url"`
}

func GetGatewayUrl() string {
	gatewayUrlRequest, err := http.Get("https://discord.com/api/v10/gateway")
	if err != nil {
		log.Fatalln(err)
	}

	gatewayUrlRequestBody, err := io.ReadAll(gatewayUrlRequest.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var gatewayUrl Gateway
	json.Unmarshal(gatewayUrlRequestBody, &gatewayUrl)

	return gatewayUrl.Url
}
