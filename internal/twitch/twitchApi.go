package twitch

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PierreOudin/TheEvilBot/internal/utils"
)

var clientId string
var clientSecret string

const (
	TWITCH_AUTH_URL string = "https://id.twitch.tv/oauth2/token"
)

func init() {
	clientId = utils.GoDotEnvVariable("TWITCH_CLIENTID")
	clientSecret = utils.GoDotEnvVariable("TWITCH_CLIENTSECRET")
}

func GetTwitchToken() {

	data := []byte(fmt.Sprintf(`{"client_id": "%v", "client_secret": "%v", "grant_type": "client_credentials"}`, clientId, clientSecret))

	client := &http.Client{}

	req, err := http.NewRequest("POST", TWITCH_AUTH_URL, bytes.NewBuffer(data))

	if err != nil {
		log.Fatalf("Http Request error : %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	//req.Header.Add("Authorization", "Bearer YOUR_ACCESS_TOKEN")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error while requesting : %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(body))

	//resp, err := http.Post(TWITCH_AUTH_URL)
}
