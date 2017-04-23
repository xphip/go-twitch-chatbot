package ircBot

import (
	"log"
	"fmt"
	"net"
	"bufio"
	"time"
	"net/textproto"
)

type Bot struct {
	Server string
	Port string
	Nick string
	User string
	Channel string
	Pass string
	Pread, pwrite chan string
	Channels []chan bool
	Conn net.Conn
	Buffer *textproto.Reader
}

// Cria novo Bot
func NewBot(server string, port string, nick string, channel string, pass string) *Bot {
	return &Bot{
		Server: server,
		Port: port,
		Nick: nick,
		Channel: channel,
		Pass: pass,
	}
}

// Conecta o bot
func (bot *Bot) Connect() (conn net.Conn, err error) {
	bot.Conn, err = net.Dial("tcp", bot.Server + ":" + bot.Port)
	if err != nil{
		log.Fatal("Unable to connect to IRC server ", err)
		bot.Conn.Close()
		conn = bot.Conn
		return conn, err
	}

	log.Printf("Connected to IRC server %s (%s) \n", bot.Server, bot.Conn.RemoteAddr())
	fmt.Fprintf(bot.Conn, "PASS %s\r\n", bot.Pass)
	// fmt.Fprintf(bot.Conn, "USER %s 8 * :%s\r\n", bot.Nick, bot.Nick)
	fmt.Fprintf(bot.Conn, "NICK %s\r\n", bot.Nick)
	fmt.Fprintf(bot.Conn, "JOIN #%s\r\n", bot.Channel)
	
	conn = bot.Conn
	return conn, err
}

// Cria o buffer
func (bot *Bot) Maker() (err error) {
	reader := bufio.NewReader(bot.Conn)
	bot.Buffer = textproto.NewReader(reader)
	
	return err
}

// Envia comando puro para o IRC
func (bot *Bot) Send(command string) {
	fmt.Fprintf(bot.Conn, "%s\r\n", command)
}

// Envia uma mensagem para um determinado canal
func (bot *Bot) Msg(channel string, msg string) {
	fmt.Fprintf(bot.Conn, "PRIVMSG #%s :%s\r\n", channel, msg)
}

// Entra em um canal
func (bot *Bot) Join(channel string) {
	fmt.Fprintf(bot.Conn, "JOIN #%s\r\n", channel)
}

// Sai de um canal
func (bot *Bot) Leave(channel string) {
	fmt.Fprintf(bot.Conn, "PART #%s\r\n", channel)
}

// Função que dispara uma checagem pra dizer
//	"servidor eu to vivo", a cada 2 minutos
func (bot *Bot) PingLoop() (err error) {
	for {
		time.Sleep(time.Second * 120)
		bot.Send("PING :tmi.twitch.tv")
	}
}