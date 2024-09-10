package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/PierreOudin/TheEvilBot/internal/discord"
)

func main() {
	// discordToken := goDotEnvVariable("DISCORD_TOKEN")

	// dg, err := discordgo.New("Bot " + discordToken)

	// if err != nil {
	// 	log.Fatalln("Error opening discord session")
	// }

	s := discord.InitDiscordBot()

	//twitch.GetTwitchToken()

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
