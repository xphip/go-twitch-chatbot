package ircbot

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"time"
)

var crlfStripper = strings.NewReplacer("\r", "", "\n", "")

type Bot struct {
	addr    string
	nick    string
	channel string
	pass    string
	conn    net.Conn
	reader  *textproto.Reader
}

func New(addr, nick, channel, pass string) *Bot {
	return &Bot{
		addr:    addr,
		nick:    nick,
		channel: channel,
		pass:    pass,
	}
}

func (b *Bot) Connect() error {
	conn, err := tls.Dial("tcp", b.addr, nil)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", b.addr, err)
	}
	b.conn = conn
	b.reader = textproto.NewReader(bufio.NewReader(conn))

	b.Send("PASS " + b.pass)
	b.Send("NICK " + b.nick)
	b.Join(b.channel)
	return nil
}

func (b *Bot) Close() error {
	return b.conn.Close()
}

func (b *Bot) ReadLine() (string, error) {
	return b.reader.ReadLine()
}

func (b *Bot) Send(command string) {
	fmt.Fprintf(b.conn, "%s\r\n", crlfStripper.Replace(command))
}

func (b *Bot) Msg(channel, msg string) {
	b.Send(fmt.Sprintf("PRIVMSG #%s :%s", channel, msg))
}

func (b *Bot) Join(channel string) {
	fmt.Fprintf(b.conn, "JOIN #%s\r\n", channel)
}

func (b *Bot) Leave(channel string) {
	fmt.Fprintf(b.conn, "PART #%s\r\n", channel)
}

func (b *Bot) PingLoop(interval time.Duration) {
	for range time.Tick(interval) {
		b.Send("PING :tmi.twitch.tv")
	}
}
