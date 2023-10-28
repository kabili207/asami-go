package listeners

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

type Matcher struct {
	Name        string
	Pattern     string
	Replacement string
	Regex       *regexp.Regexp
}

type LinkFixer struct {
	Matchers []*Matcher
}

func NewLinkFixer() *LinkFixer {
	// TODO: Consider moving to config file?
	matchers := []*Matcher{
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
		{
			Name:        "DeviantArt",
			Pattern:     `^https?://(?:www\.)?deviantart.com/(.*)`,
			Replacement: "https://www.fxdeviantart.com/$1",
		},
		{
			Name:        "FurAffinity",
			Pattern:     `^https?://(?:www\.)?furaffinity.net/(.*)`,
			Replacement: "https://fxfuraffinity.net/$1",
		},
	}
	for _, m := range matchers {
		m.Regex = regexp.MustCompile(m.Pattern)
	}

	return &LinkFixer{
		Matchers: matchers,
	}
}

func (l *LinkFixer) RegisterHandlers(s *discordgo.Session) {

	s.AddHandler(l.messageCreate)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (l *LinkFixer) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	urls := make([]string, 0)

	rx := xurls.Strict()

	origUrls := rx.FindAllString(m.Content, -1)

	for _, u := range origUrls {
		for _, m := range l.Matchers {
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
