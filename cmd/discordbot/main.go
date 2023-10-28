package discordbot

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kabili207/asamigo/internal/commands"
	"github.com/kabili207/asamigo/internal/listeners"
	"github.com/zekrotja/ken"
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
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	linkFixer := listeners.NewLinkFixer()
	linkFixer.RegisterHandlers(session)

	k, err := ken.New(session)

	k.RegisterCommands(
		new(commands.Moths),
	)

	defer k.Unregister()

	// In this example, we only care about receiving message events.
	session.Identify.Intents = discordgo.IntentsGuildMessages

	err = session.Open()
	defer session.Close()
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
}
