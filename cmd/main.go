package main

import (
	"botty/internal/discord"
	"botty/internal/openai"
	"flag"
	"fmt"
)

func main() {
	discordToken := flag.String("discord", "", "Discord bot token")
	openAIKey := flag.String("openai", "", "OpenAI API key")
	engine := flag.String("engine", "", "Change default engine")
	channel := flag.String("channel", "", "Limit the bot to a single channel ID")
	flag.Parse()

	if *discordToken != "" && *openAIKey != "" {
		discord.LimitChannelID = *channel
		if *engine != "" {
			discord.GptClient = openai.NewWithEngine(*openAIKey, *engine)
			discord.NewBot(*discordToken)
		} else {
			discord.GptClient = openai.New(*openAIKey)
			discord.NewBot(*discordToken)
		}
	} else {
		fmt.Println("required flags (discord, openai) not provided")
	}
}
