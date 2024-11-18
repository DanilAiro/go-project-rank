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

	command := getCommands()

	cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, channels[0].GuildID, command)
	if err != nil {
		fmt.Printf("Error adding command: %s\n", err.Error())
	}

	// Block until we get ctrl-c
	fmt.Println("Bot running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Clean up
	fmt.Println("Bot exiting")
	err = discord.ApplicationCommandDelete(discord.State.User.ID, channels[0].GuildID, cmd.ID)
	if err != nil {
		fmt.Printf("Error removing command: %s\n", err.Error())
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

func getCommands() *discordgo.ApplicationCommand {
	command := &discordgo.ApplicationCommand{
		Name:        "test",
		Description: "A test command with subcommand-groups and subcommands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "test-a",
				Description: "Test-a sub-command group",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "test-a-a",
						Description: "Test-a-a sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "test-a-b",
						Description: "Test-a-b sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
			},
			{
				Name:        "test-b",
				Description: "Test-b sub-command group",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "test-b-a",
						Description: "Test-b-a sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "test-b-b",
						Description: "Test-b-b sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
			},
			{
				Name:        "test-c",
				Description: "Test-c sub-command",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}

	return command
}

func (dh *discordHandler) ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateListeningStatus("/test")
}

func (dh *discordHandler) command(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {

	case "test":
		options := i.ApplicationCommandData().Options
		response := ""

		switch options[0].Name {

		case "test-c":
			response = "Test C Command"

		case "test-a":
			options := options[0].Options
			switch options[0].Name {

			case "test-a-a":
				response = "Test A A Command"

			case "test-a-b":
				response = "Test A B Command"

			default:
				response = "Error!"
			}

		case "test-b":
			options := options[0].Options
			switch options[0].Name {

			case "test-b-a":
				response = "Test B A Command"

			case "test-b-b":
				response = "Test B B Command"

			default:
				response = "Error!"
			}

		default:
			response = "Error!"
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	}
}