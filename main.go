package main

/*
This is the first component of the trading bot

-> The messenging system

Bot Components:
-> Messenging System (Delivers messages to inform stakeholder of trades and information from the trading system)
-> Trading System (Using a certain trading strategy, this component will execute trades)

Once this code is operational, it will be legacy code. This bot needs to work
and not be perfect. After completion, launch this to aws.

*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Discord-bot Configuration
var (
	Prefix = "!"
	Colour = 0x009688
	Name   = "Araphor0"

	//API Token (required)
	Token = "ODE1NTg1MjM2MTg0MDcyMjAy.YDui8Q.QPPGrgu0BhBfJ12gIJd3kC2QDh8"
)

// This function is called when loading the bot
func init() {
	flag.StringVar(&Token, "t", Token, "Bot Token")
	flag.Parse()
}

// Bot startup/closing is handled here
func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Unable to create discord session!")
		return
	}

	// Setup handlers.
	dg.AddHandler(messageCreate)
	dg.AddHandler(ready)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Unable to open connection!")
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("[ CONNECTED ] Araphor0 ")
}

// This is called on every message received
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}

	// Do not respond to self
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Help embed
	help := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{Name: Name + " Commands"},
		Color:  Colour,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   Prefix + "help",
				Value:  "Display a list of commands.",
				Inline: true,
			},
			{
				Name:   Prefix + "hello",
				Value:  "Displays a greeting message",
				Inline: true,
			},
			{
				Name:   Prefix + "rand",
				Value:  "Generates a random number between 1 and 10",
				Inline: true,
			},
			{
				Name:   Prefix + "test",
				Value:  "Tests the latency of a single message",
				Inline: true,
			},
			{
				Name:   Prefix + "yeezy",
				Value:  "Displays a quote from Ye.",
				Inline: true,
			},
			{
				Name:   "AraphorVPS",
				Value:  "Access AraphorVPS using ~help in the vps-management channel",
				Inline: true,
			},
			{
				Name:   "matthew",
				Value:  "does this work?",
				Inline: true,
			},
		},
	}
	if strings.HasPrefix(m.Content, Prefix+"help") {
		s.ChannelMessageSendEmbed(m.ChannelID, help)
	}

	// Generate a greeting message
	if strings.HasPrefix(m.Content, Prefix+"hello") {
		s.ChannelMessageSend(m.ChannelID, "Hello!")
	}

	// Generates a random number between 1 and 10
	if strings.HasPrefix(m.Content, Prefix+"rand") {
		rand.Seed(time.Now().UnixNano())

		//This generates a random number between 1 and 10 while converting
		//the integer to a string
		randomInteger := strconv.Itoa(1 + rand.Intn(11-1))
		s.ChannelMessageSend(m.ChannelID, randomInteger)
	}

	if strings.HasPrefix(m.Content, Prefix+"test") {
		start := time.Now()
		s.ChannelMessageSend(m.ChannelID, "*test*: ")
		duration := time.Since(start)
		s.ChannelMessageSend(m.ChannelID, duration.String())
	}

	if strings.HasPrefix(m.Content, Prefix+"yeezy") {
		resp, err := client.Get("https://api.kanye.rest") //api.kanye.rest is full of kanye quotes
		if resp != nil {
			defer resp.Body.Close()
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Yeezy failed!")
		}
		//Convert the body to type string
		s.ChannelMessageSend(m.ChannelID, string(body))
	}

	if strings.HasPrefix(m.Content, Prefix+"matthew") {
		s.ChannelMessageSend(m.ChannelID, "Computer System Fastlane")
	}
}
