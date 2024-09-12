package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/PierreOudin/TheEvilBot/internal/discord"
	"github.com/PierreOudin/TheEvilBot/internal/twitch"
)

func main() {
	// discordToken := goDotEnvVariable("DISCORD_TOKEN")

	// dg, err := discordgo.New("Bot " + discordToken)

	// if err != nil {
	// 	log.Fatalln("Error opening discord session")
	// }

	s := discord.InitDiscordBot()

	twitch.GetTwitchToken()

	twitch.GetStreams("striikerrr_")

	// if err != nil {
	// 	log.Fatalf("Error : %v", err)
	// }

	// fmt.Println(token)

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
