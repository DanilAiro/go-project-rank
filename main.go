package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Printf("Error creating Discord session: %s\n", err.Error())
		os.Exit(2)
	}

	err = discord.Open()
	if err != nil {
		fmt.Printf("Error opening Discord connection: %s\n", err.Error())
		os.Exit(3)
	}

	channels, err := findRankChannels(discord)
	if err == nil {
		fmt.Println(channels[0].Name, 1)
		discord.ChannelMessageSend(channels[0].ID, "Я работаю")
	}

	defer discord.Close()
}

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "-t <discord_bot_token>")
	flag.Parse()

	if Token == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func findRankChannels(discord *discordgo.Session) (c []*discordgo.Channel, err error) {
	var channels []*discordgo.Channel

	if len(discord.State.Guilds) == 0 {
		return channels, errors.New("no guilds available")
	}

	for _, guild := range discord.State.Guilds {
		c, _ := discord.GuildChannels(guild.ID)
		for _, channel := range c {
			if channel.Name == "rank" || channel.Name == "ранг" {
				channels = append(channels, channel)
				break
			}
		}
	}

	if len(channels) == 0 {
		return channels, errors.New("no channels available")
	}

	return channels, nil
}