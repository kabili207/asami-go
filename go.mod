module github.com/kabili207/asami-go

go 1.19

require github.com/bwmarrin/discordgo v0.27.1

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	mvdan.cc/xurls/v2 v2.5.0
)

//replace github.com/bwmarrin/discordgo v0.27.1 => github.com/kabili207/discordgo master

replace github.com/bwmarrin/discordgo => github.com/kabili207/discordgo v0.27.1-fix
