package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "log"

    "github.com/bwmarrin/discordgo"
    gecko "github.com/superoo7/go-gecko/v3"
)

// Variables used for command line parameters
var (
    Token string
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
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

    cg := gecko.NewClient(nil)

    ids := []string{"bitcoin", "ethereum","dogecoin"}
    vc := []string{"usd", "inr"}
    sp, err := cg.SimplePrice(ids, vc)
    if err != nil {
        log.Fatal(err)
    }
    bitcoin := (*sp)["bitcoin"]
    eth := (*sp)["ethereum"]
    doge := (*sp)["dogecoin"]
    reply := fmt.Sprintf("Bitcoin is worth %f usd (inr %f)", bitcoin["usd"], bitcoin["inr"])
    reply2 := fmt.Sprintf("Ethereum is worth %f usd (inr %f)", eth["usd"], eth["inr"])
    reply3 := fmt.Sprintf("Doge is worth %f usd (inr %f)", doge["usd"], doge["inr"])

    if m.Author.ID == s.State.User.ID {
        return
    }
    if m.Content == "btc" {
        s.ChannelMessageSend(m.ChannelID, reply )
    }

    if m.Content == "eth" {
        s.ChannelMessageSend(m.ChannelID, reply2)
    }
    if m.Content == "doge" {
        s.ChannelMessageSend(m.ChannelID, reply3)
    }
}
