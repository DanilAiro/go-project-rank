package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	BotToken string `split_words:"true" required:"true"`
	// GuildID  string `split_words:"true" required:"true"`
}

type discordHandler struct {
	config Config
}

func discordMain(token string) {
	var config Config

	config.BotToken = token

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("Error creating Discord session: %s\n", err.Error())
		os.Exit(2)
	}

	// Register callbacks
	dh := &discordHandler{}
	discord.AddHandler(dh.ready)
	discord.AddHandler(dh.command)

	err = discord.Open()
	if err != nil {
		fmt.Printf("Error opening Discord connection: %s\n", err.Error())
		os.Exit(3)
	}

	channels, err := findRankChannels(discord)
	if err != nil {
		fmt.Printf("Error getting rank channels: %s\n", err.Error())
		os.Exit(4)
	}

	commands := getCommands()
	var listenCommands []*discordgo.ApplicationCommand

	for _, command := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, channels[0].GuildID, command)
		if err != nil {
			fmt.Printf("Error adding command: %s\n", err.Error())
		}
		listenCommands = append(listenCommands, cmd)
	}

	// Block until we get ctrl-c
	fmt.Println("Bot running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Clean up
	fmt.Println("Bot exiting")
	for _, command := range listenCommands {
		err = discord.ApplicationCommandDelete(discord.State.User.ID, channels[0].GuildID, command.ID)
		if err != nil {
			fmt.Printf("Error removing command: %s\n", err.Error())
		}
	}

	discord.Close()
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

func getCommands() []*discordgo.ApplicationCommand {
	var commands []*discordgo.ApplicationCommand
	command := &discordgo.ApplicationCommand{
		Name:        "link",
		Description: "link discord user with valorant user",
		Options:     []*discordgo.ApplicationCommandOption{
		},
	}

	commands = append(commands, command)

	command = &discordgo.ApplicationCommand{
		Name:        "unlink",
		Description: "unlink discord user with valorant user",
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	commands = append(commands, command)

	command = &discordgo.ApplicationCommand{
		Name:        "bug",
		Description: "report bug",
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	commands = append(commands, command)

	return commands
}

func (dh *discordHandler) ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateListeningStatus("/link")
	s.UpdateListeningStatus("/unlink")
	s.UpdateListeningStatus("/bug")
}

func (dh *discordHandler) command(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {

	case "link":
		response := ""
		
		message, err := s.ChannelMessageSend(i.ChannelID, "Введите данные в формате - nickname tag")
		if err != nil {
			fmt.Printf("Error getting rank channels: %s\n", err.Error())
			os.Exit(5)
		}

		fmt.Println(message.Content)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	}
}
