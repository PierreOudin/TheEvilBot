package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/PierreOudin/TheEvilBot/internal/utils"
)

type Token struct {
	token           string
	expiration_date time.Time
}

type twitchTokenResponse struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Token_type   string `json:"token_type"`
}

var clientId string
var clientSecret string
var tokenData Token

const (
	TWITCH_AUTH_URL string = "https://id.twitch.tv/oauth2/token"
)

func init() {
	clientId = utils.GoDotEnvVariable("TWITCH_CLIENTID")
	clientSecret = utils.GoDotEnvVariable("TWITCH_CLIENTSECRET")
}

func getTwitchToken() {

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

	var tokenResponse twitchTokenResponse

	err = json.Unmarshal(body, &tokenResponse)

	if err != nil {
		log.Fatalf("Error while decoding body : %v", err)
	}

	tokenData.token = tokenResponse.Access_token
	tokenData.expiration_date = time.Now().Local().Add(time.Second * time.Duration(tokenResponse.Expires_in))
}

func validateTwitchToken() {
	if tokenData.token == "" {
		getTwitchToken()
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)

	if err != nil {
		log.Fatalf("Error : %v", err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tokenData.token))

	fmt.Printf(tokenData.token)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error while requesting : %v", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var dataBody map[string]interface{}

	err = json.Unmarshal(body, &dataBody)

	fmt.Printf("dofy : %v", dataBody)

	if err != nil {
		log.Fatalf("Error while decoding body : %v", err)
		return
	}

	nbr, ok := dataBody["expire_in"].(int)

	if ok {
		if nbr < 60 {
			getTwitchToken()
		}
	}
}

func GetStream(streamer string) map[string]interface{} {
	validateTwitchToken()

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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tokenData.token))

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

	fmt.Printf("dofy : %v", dataBody)

	if err != nil {
		log.Fatalf("Error while decoding body : %v", err)
		return nil
	}

	return dataBody
}
