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

	for _, cmd := range []string{"PASS " + b.pass, "NICK " + b.nick} {
		if err := b.Send(cmd); err != nil {
			return err
		}
	}
	return b.Join(b.channel)
}

func (b *Bot) Close() error {
	return b.conn.Close()
}

func (b *Bot) ReadLine() (string, error) {
	return b.reader.ReadLine()
}

func (b *Bot) Send(command string) error {
	_, err := fmt.Fprintf(b.conn, "%s\r\n", crlfStripper.Replace(command))
	return err
}

func (b *Bot) Msg(channel, msg string) error {
	return b.Send(fmt.Sprintf("PRIVMSG #%s :%s", channel, msg))
}

func (b *Bot) Join(channel string) error {
	return b.Send("JOIN #" + channel)
}

func (b *Bot) Leave(channel string) error {
	return b.Send("PART #" + channel)
}

func (b *Bot) PingLoop(interval time.Duration) {
	for range time.Tick(interval) {
		if err := b.Send("PING :tmi.twitch.tv"); err != nil {
			return
		}
	}
}
