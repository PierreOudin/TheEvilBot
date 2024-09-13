package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PierreOudin/TheEvilBot/internal/utils"
)

var clientId string
var clientSecret string

var token string
var refreshToken string

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
	}

	var dataBody map[string]interface{}

	err = json.Unmarshal(body, &dataBody)

	if err != nil {
		log.Fatalf("Error while decoding body : %v", err)
	}

	str, ok := dataBody["access_token"].(string)
	if ok {
		token = str
	}

	str, ok = dataBody["refresh_token"].(string)
	if ok {
		refreshToken = str
	}

	//resp, err := http.Post(TWITCH_AUTH_URL)
}

func GetStream(streamer string) map[string]interface{} {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams", nil)

	if err != nil {
		log.Fatalf("Error : %v", err)
		return nil
	}

	q := req.URL.Query()
	q.Add("user_login", streamer)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Client-ID", clientId)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error while requesting : %v", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	var dataBody map[string]interface{}

	err = json.Unmarshal(body, &dataBody)

	if err != nil {
		log.Fatalf("Error while decoding body : %v", err)
		return nil
	}

	return dataBody
}
