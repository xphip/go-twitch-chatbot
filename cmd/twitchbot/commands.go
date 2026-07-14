package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dpkat/go-twitch-chatbot/ircbot"
)

const (
	commandCooldown = 10 * time.Second
	maxBodySize     = 1 << 20
)

var (
	httpClient = &http.Client{Timeout: 5 * time.Second}
	lastDolar  time.Time
)

func handleCommand(bot *ircbot.Bot, channel, command string) bool {
	switch command {
	case "!dolar":
		if time.Since(lastDolar) < commandCooldown {
			return true
		}
		lastDolar = time.Now()
		dolar, err := fetchDolar()
		if err != nil {
			log.Print("!dolar: ", err)
			return true
		}
		if err := bot.Msg(channel, "[DOLAR hoje] R$ "+dolar); err != nil {
			log.Print("!dolar: ", err)
		}
		return true
	}
	return false
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
