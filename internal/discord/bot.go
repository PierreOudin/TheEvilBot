package discord

import (
	"log"

	discordcommands "github.com/PierreOudin/TheEvilBot/internal/discord/discord_commands"
	"github.com/PierreOudin/TheEvilBot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

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
		"add": discordcommands.AddStreamers,
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
