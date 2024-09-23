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
		data, err := twitch.GetStream("striikerrr_")

		if err != nil {
			log.Fatalf("Error while getting stream info : %v", err)
		}

		var category string = data.Data[0].GameName
		var discordStreamName string = strings.ReplaceAll(data.Data[0].UserLogin, "_", "\\_")
		message := fmt.Sprintf("🚀 @everyone %v vient de commencer un stream sur %v! Regardez-le ici: https://www.twitch.tv/%v", discordStreamName, category, data.Data[0].UserLogin)
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
