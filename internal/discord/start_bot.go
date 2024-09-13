package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/PierreOudin/TheEvilBot/internal/twitch"
	"github.com/bwmarrin/discordgo"
)

func StartBot(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)

	var streamChannelID string = ""

	guilds := r.Guilds[0]

	// for _, g := range guilds {
	// 	fmt.Printf("Guild : %v | Name : %v", g.ID, g.Name)
	// }

	channel, err := s.GuildChannels(guilds.ID)

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	for _, c := range channel {
		if strings.ToLower(c.Name) == "stream" {
			streamChannelID = c.ID
		}
	}

	if streamChannelID != "" {
		data := twitch.GetStream("striikerrr_")

		mapD, ok := data["data"].(map[string]interface{})
		var mapData map[string]interface{}
		if ok {
			log.Printf("mapData : %v", mapD)
			mapData = mapD
		}
		var category string = ""
		str, ok := mapData["game_name"].(string)
		if ok {
			category = str
		}
		message := fmt.Sprintf("🚀 @everyone striikerrr_ vient de commencer un stream sur %v! Regardez-le ici: https://www.twitch.tv/striikerrr_", category)
		s.ChannelMessageSend(streamChannelID, message)
		// interval := time.NewTicker(2 * time.Minute)

		// for {
		// 	select {
		// 	case <-interval.C:
		// 		data := twitch.GetStream("striikerrr_")
		// 		s.ChannelMessageSend(streamChannelID, fmt.Sprintf("MessageCreate %v", data))
		// 	}
		// }

	}

	log.Printf("MessageCreate %v", r.Application.Name)
}
