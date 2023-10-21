module github.com/kabili207/asamigo

go 1.19

require github.com/bwmarrin/discordgo v0.27.1

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/zekrotja/dgrs v0.5.7 // indirect
	github.com/zekrotja/safepool v1.1.0 // indirect
	golang.org/x/net v0.10.0 // indirect
)

require (
	github.com/PuerkitoBio/goquery v1.8.1
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/zekrotja/ken v0.20.0
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	mvdan.cc/xurls/v2 v2.5.0
)

//replace github.com/bwmarrin/discordgo v0.27.1 => github.com/kabili207/discordgo master

replace github.com/bwmarrin/discordgo => github.com/kabili207/discordgo v0.27.1-fix
