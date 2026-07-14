package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dpkat/go-twitch-chatbot/ircbot"
)

const (
	server          = "irc.chat.twitch.tv:6697"
	commandCooldown = 10 * time.Second
	maxBodySize     = 1 << 20
)

var (
	msgRegex   = regexp.MustCompile(`:([a-z0-9_]*)!([a-z0-9_]*)@([a-z0-9.-]*) ([A-Z]*) #([a-z0-9_]*) :(![a-z0-9_]*)?(.*)`)
	httpClient = &http.Client{Timeout: 5 * time.Second}
)

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

	var lastDolar time.Time

	for {
		line, err := bot.ReadLine()
		if err != nil {
			log.Fatal("read from IRC server: ", err)
		}

		if strings.HasPrefix(line, "PING ") {
			bot.Send("PONG " + strings.TrimPrefix(line, "PING "))
			continue
		}

		m := msgRegex.FindStringSubmatch(line)
		if m == nil {
			fmt.Println(line)
			continue
		}
		username, channel, command, text := m[1], m[5], m[6], m[7]

		switch command {
		case "!dolar":
			if time.Since(lastDolar) < commandCooldown {
				continue
			}
			lastDolar = time.Now()
			dolar, err := fetchDolar()
			if err != nil {
				log.Print("!dolar: ", err)
				continue
			}
			bot.Msg(channel, "[DOLAR hoje] R$ "+dolar)
		default:
			fmt.Printf("[#%s] %s: %s%s\n", channel, username, command, text)
		}
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}

func fetchDolar() (string, error) {
	res, err := httpClient.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var quote struct {
		USDBRL struct {
			Bid string `json:"bid"`
		} `json:"USDBRL"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, maxBodySize)).Decode(&quote); err != nil {
		return "", err
	}
	return quote.USDBRL.Bid, nil
}
