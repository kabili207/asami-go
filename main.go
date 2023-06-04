package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

// Variables used for command line parameters
var (
	Token          string
	BotPermissions string = "274878262336"
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
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

	fmt.Printf("Invite link: https://discord.com/api/oauth2/authorize?client_id=768941376028016651&permissions=%v&scope=bot", BotPermissions)

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

	origString := m.Content
	urls := make([]string, 0)

	rx := xurls.Strict()

	origUrls := rx.FindAllString(origString, -1)

	rx_insta := regexp.MustCompile(`^https?://(?:www\.)?instagram.com/(.*)$`)
	rx_tik := regexp.MustCompile(`^https?://(?:www\.)?tiktok.com/(.*)$`)

	for _, u := range origUrls {
		if rx_insta.MatchString(u) {
			urls = append(urls, rx_insta.ReplaceAllString(u, "https://ddinstagram.com/$1"))
		}
		if rx_tik.MatchString(u) {
			urls = append(urls, rx_tik.ReplaceAllString(u, "https://vxtiktok.com/$1"))
		}
	}

	urls = dedupSlice(urls)

	// If the message is "ping" reply with "Pong!"
	if len(urls) > 0 {
		s.ChannelMessageSend(m.ChannelID, strings.Join(urls, "\n"))

		//time.Sleep(4 * time.Second)
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: m.ChannelID,
			ID:      m.ID,
			Flags:   m.Flags | discordgo.MessageFlagsSuppressEmbeds})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func dedupSlice[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
