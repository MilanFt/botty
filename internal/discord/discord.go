package discord

import (
	"botty/internal/openai"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/bwmarrin/discordgo"
)

var (
	botName  string
	identity string
	users    []string
	logs     []string
)

var GptClient gpt3.Client
var LimitChannelID string

// Initialize a discord bot with the provided token.
// Boilerplate code taken from https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
func NewBot(token string) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	_, err = dg.UserUpdate("Botty", BottyAvatarBase64)
	if err != nil {
		fmt.Println("error updating user,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If we provided a channel ID, only respond to messages in that channel
	if LimitChannelID != "" && m.ChannelID != LimitChannelID {
		return
	}

	// Tell Discord to show the typing indicator.
	s.ChannelTyping(m.ChannelID)

	var cmd, arg string
	var found bool

	// If the message begins with '!' we parse the provided command and its argument.
	if strings.HasPrefix(m.Content, "!") {
		cmd, arg, found = strings.Cut(m.Content, " ")
		if !found {
			s.ChannelMessageSend(m.ChannelID, "Argument for command not found")
			return
		}
	}

	if found {
		switch cmd {
		// Set the bot's name, this is required so that when we request
		// OpenAI's API we can refer to the bot with a name instead of
		// the default 'AI:'.
		case "!setname":
			botName = arg
			s.ChannelMessageSend(m.ChannelID, "New bot name: "+botName)
		// We need to provide context to the conversation, writing a comprehensive
		// identity is crucial, because OpenAI's model performance heavily depends on it.
		case "!setidentity":
			identity = arg
			s.ChannelMessageSend(m.ChannelID, "New identity: "+identity)
		// We can add up to 4 users to the conversation, this is the most currently supported
		// because of OpenAI's 'Stop sequences' limitation.
		case "!adduser":
			if len(users) < 4 {
				users = append(users, arg)
				s.ChannelMessageSend(m.ChannelID, "New conversation partner added: "+arg)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Conversation partner limit (4) reached")
			}
		// Remove a user from the conversation, their messages will not be stored anymore.
		case "!removeuser":
			var userIndex int = -1
			for i, v := range users {
				if v == arg {
					userIndex = i
				}
			}
			if userIndex != -1 {
				users = remove(users, userIndex)
				s.ChannelMessageSend(m.ChannelID, "Conversation partner removed: "+arg)

			} else {
				s.ChannelMessageSend(m.ChannelID, "Conversation partner not found: "+arg)
			}
		// Clear all parameters to completely reset the conversation.
		case "!clear":
			botName = ""
			identity = ""
			logs = []string{}
			users = []string{}
			s.ChannelMessageSend(m.ChannelID, "All parameters (name, identity, chat log, users) are cleared")
		// Only clear the logs. The bot's identity and the conversation partners will remain.
		case "!clearlogs":
			logs = []string{}
			s.ChannelMessageSend(m.ChannelID, "Cleared chat logs")
		// Clear the bot's identity, this also resets the chat logs.
		case "!clearidentity":
			identity = ""
			logs = []string{}
			s.ChannelMessageSend(m.ChannelID, "Cleared identity and former chat logs")
		default:
			s.ChannelMessageSend(m.ChannelID, "Unknown command: "+cmd)
		}
		return
	}

	var foundUser bool
	for _, v := range users {
		if m.Author.Username == v {
			foundUser = true
		}
	}

	if !foundUser {
		s.ChannelMessageSend(m.ChannelID, "You are not set as a conversation partner, please use the provided command below to join the conversation: ```!adduser [user name]```")
		return
	}

	if botName == "" {
		s.ChannelMessageSend(m.ChannelID, "The bot has no name, please use the provided command below: ```!setname [name]```")
		return
	}

	if identity == "" {
		s.ChannelMessageSend(m.ChannelID, "The bot has no identity, please use the provided command below: ```!setidentity [story]```")
		return
	}

	if len(users) == 0 {
		s.ChannelMessageSend(m.ChannelID, "The bot has conversation partners set, please use the provided command below: ```!adduser [user name]```")
		return
	}

	logString := fmt.Sprintf("%s: %s", m.Author.Username, m.Content)
	logs = append(logs, logString)

	prompt := openai.CreatePrompt(
		botName,
		identity,
		logs,
	)

	resp, err := openai.Complete(GptClient, prompt, users)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "[Failed to generate response]")
		return
	}

	botLogString := fmt.Sprintf("%s: %s", botName, resp)
	logs = append(logs, botLogString)

	s.ChannelMessageSend(m.ChannelID, resp)
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
