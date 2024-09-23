package discord

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PierreOudin/TheEvilBot/internal/twitch"
	"github.com/PierreOudin/TheEvilBot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

var Streamers []StreamerInfo

type StreamerInfo struct {
	StreamerName string
	DiscordID    int
	StreamInfo   streamInfo
}
type streamInfo struct {
	LastStream time.Time
	IsOnline   bool
	Category   string
}

func init() {
	var err error
	discordToken := utils.GoDotEnvVariable("DISCORD_TOKEN")
	s, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "streamers",
			Description: "List all stream we follow",
		},
		{
			Name:        "add",
			Description: "Add a stream to follow",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "streamer-name",
					Description: "Streamer name. To check if they are online",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "delete",
			Description: "Unfollow a stream",
		},
		{
			Name:        "laststream",
			Description: "Date of the last stream",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"streamers": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Command streamer",
				},
			})
		},
		"add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			streamerName := i.ApplicationCommandData().Options[0].StringValue()
			var message string
			added, err := BotAddStreamers(streamerName)
			if added {
				message = fmt.Sprintf("Le streamer %v a été ajouté à la liste", streamerName)
			} else {
				if err != nil {
					message = fmt.Sprint(err)
				} else {
					message = fmt.Sprintf("Le streamer %v est déjà présent dans la liste", streamerName)
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message,
				},
			})
		},
		"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Command delete",
				},
			})
		},
		"laststream": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Command laststream",
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func InitDiscordBot() *discordgo.Session {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(StartBot)

	err := s.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	for _, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		if cmd != nil {
			log.Printf("Create command %v", cmd.Name)
		}
	}

	return s
}

func BotAddStreamers(streamer string) (bool, error) {
	twitchExist := twitch.StreamExist(streamer)

	if !twitchExist {
		return false, errors.New("le streamer n'existe pas")
	}

	if Streamers != nil {
		var alreadyExist bool = false
		for _, s := range Streamers {
			if s.StreamerName == streamer {
				alreadyExist = true
			}
		}
		if !alreadyExist {
			Streamers = append(Streamers, StreamerInfo{StreamerName: streamer})
			return true, nil
		}
		return false, nil
	} else {
		Streamers = []StreamerInfo{
			{
				StreamerName: streamer,
			},
		}
		return true, nil
	}
}
