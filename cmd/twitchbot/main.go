package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dpkat/go-twitch-chatbot/ircbot"
)

const server = "irc.chat.twitch.tv:6697"

func main() {
	cfg := loadConfig()

	bot := ircbot.New(server, cfg.Nick, cfg.Channel, cfg.Token)
	if err := bot.Connect(); err != nil {
		log.Fatal(err)
	}
	defer bot.Close()
	log.Printf("connected to %s as %s on #%s", server, cfg.Nick, cfg.Channel)

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

		msg, ok := ircbot.Parse(line)
		if !ok {
			fmt.Println(line)
			continue
		}

		if handleCommand(bot, msg.Channel, msg.Command) {
			continue
		}
		fmt.Printf("[#%s] %s: %s%s\n", msg.Channel, msg.Username, msg.Command, msg.Text)
	}
}
