package listeners

import "github.com/bwmarrin/discordgo"

type Listener interface {
	RegisterHandlers(s *discordgo.Session) []interface{}
}
