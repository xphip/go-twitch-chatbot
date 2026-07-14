package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dpkat/go-twitch-chatbot/ircbot"
)

const server = "irc.chat.twitch.tv:6697"

var msgRegex = regexp.MustCompile(`:([a-z0-9_]*)!([a-z0-9_]*)@([a-z0-9.-]*) ([A-Z]*) #([a-z0-9_]*) :(![a-z0-9_]*)?(.*)`)

func main() {
	nick := mustEnv("TWITCH_NICK")
	channel := mustEnv("TWITCH_CHANNEL")
	token := mustEnv("TWITCH_TOKEN")

	bot := ircbot.New(server, nick, channel, token)
	if err := bot.Connect(); err != nil {
		log.Fatal(err)
	}
	defer bot.Close()
	log.Printf("connected to %s as %s on #%s", server, nick, channel)

	go bot.PingLoop(2 * time.Minute)

	for {
		line, err := bot.ReadLine()
		if err != nil {
			log.Fatal("read from IRC server: ", err)
		}

		if strings.HasPrefix(line, "PING ") {
			if err := bot.Send("PONG " + strings.TrimPrefix(line, "PING ")); err != nil {
				log.Print("send PONG: ", err)
			}
			continue
		}

		m := msgRegex.FindStringSubmatch(line)
		if m == nil {
			fmt.Println(line)
			continue
		}
		username, channel, command, text := m[1], m[5], m[6], m[7]

		if handleCommand(bot, channel, command) {
			continue
		}
		fmt.Printf("[#%s] %s: %s%s\n", channel, username, command, text)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}
