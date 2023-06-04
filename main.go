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
	Matchers       []*Matcher
)

type Matcher struct {
	Name        string
	Pattern     string
	Replacement string
	Regex       *regexp.Regexp
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	// TODO: Consider moving to config file?
	Matchers = []*Matcher{
		{
			Name:        "TikTok",
			Pattern:     `^https?://(?:www\.)?tiktok.com/(.*)`,
			Replacement: "https://vxtiktok.com/$1",
		},
		{
			Name:        "Instagram",
			Pattern:     `^https?://(?:www\.)?instagram.com/(.*)`,
			Replacement: "https://ddinstagram.com/$1",
		},
		{
			Name:        "Pixiv",
			Pattern:     `^https?://(?:www\.)?pixiv.net/(.*)`,
			Replacement: "https://www.phixiv.net/$1",
		},
	}
	for _, m := range Matchers {
		m.Regex = regexp.MustCompile(m.Pattern)
	}

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

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	fmt.Printf("Invite link: https://discord.com/api/oauth2/authorize?client_id=768941376028016651&permissions=%v&scope=bot\n", BotPermissions)

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
	if m.Author.ID == s.State.User.ID {
		return
	}

	urls := make([]string, 0)

	rx := xurls.Strict()

	origUrls := rx.FindAllString(m.Content, -1)

	for _, u := range origUrls {
		for _, m := range Matchers {
			if m.Regex.MatchString(u) {
				urls = append(urls, m.Regex.ReplaceAllString(u, m.Replacement))
			}
		}
	}

	urls = dedupSlice(urls)

	if len(urls) > 0 {
		s.ChannelMessageSend(m.ChannelID, strings.Join(urls, "\n"))

		// Supress the embed from the source message
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: m.ChannelID,
			ID:      m.ID,
			Flags:   m.Flags | discordgo.MessageFlagsSuppressEmbeds,
		})

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func dedupSlice[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
