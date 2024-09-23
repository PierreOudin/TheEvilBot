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

type twitchTokenResponse struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Token_type   string `json:"token_type"`
}

type TwitchStreamerResponse struct {
	Data []twitchStreamerData `json:"data"`
}

type twitchStreamerData struct {
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	StreamID     string    `json:"id"`
	IsMature     bool      `json:"is_mature"`
	Language     string    `json:"language"`
	StartedAt    time.Time `json:"started_at"`
	TagIDs       []string  `json:"tag_ids"`
	Tags         []string  `json:"tags"`
	ThumbnailUrl string    `json:"thumbnail_url"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	ViewerCount  int       `json:"viewer_count"`
}

type twitchExistResponse struct {
	Data []streamExistResponse `json:"data"`
}

type streamExistResponse struct {
	BroadcasterType string    `json:"broadcaster_type"`
	CreatedAt       time.Time `json:"created_at"`
	Description     string    `json:"description"`
	DisplayName     string    `json:"display_name"`
	StreamerId      string    `json:"id"`
	Login           string    `json:"login"`
	OfflineImageUrl string    `json:"offline_image_url"`
	ProfileImageUrl string    `json:"profile_image_url"`
	Type            string    `json:"type"`
	ViewCount       int       `json:"view_count"`
}

var clientId string
var clientSecret string
var twitchToken string

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

	twitchToken = tokenResponse.Access_token
}

func validateTwitchToken() {
	if twitchToken == "" {
		getTwitchToken()
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)

	if err != nil {
		log.Fatalf("Error : %v", err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", twitchToken))

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

func GetStream(streamer string) (TwitchStreamerResponse, error) {
	validateTwitchToken()

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams", nil)

	if err != nil {
		log.Fatalf("Error : %v", err)
		return TwitchStreamerResponse{}, err
	}

	q := req.URL.Query()
	q.Add("user_login", streamer)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Client-ID", clientId)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", twitchToken))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error while requesting : %v", err)
		return TwitchStreamerResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return TwitchStreamerResponse{}, err
	}

	var twitchResponse TwitchStreamerResponse

	err = json.Unmarshal(body, &twitchResponse)

	if err != nil {
		log.Fatalf("Error while decoding body : %v", err)
		return TwitchStreamerResponse{}, err
	}

	return twitchResponse, nil
}

func StreamExist(streamer string) bool {
	validateTwitchToken()

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)

	if err != nil {
		log.Fatalf("Error : %v", err)
		return false
	}

	q := req.URL.Query()
	q.Add("login", streamer)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Client-ID", clientId)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", twitchToken))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error while requesting : %v", err)
		return false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	var dataBody twitchExistResponse

	err = json.Unmarshal(body, &dataBody)

	if err != nil {
		log.Fatalf("Error while decoding body : %v", err)
		return false
	}

	return len(dataBody.Data) > 0

}
