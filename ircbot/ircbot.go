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

type Message struct {
	Username string
	Channel  string
	Command  string
	Text     string
}

func Parse(line string) (Message, bool) {
	if strings.HasPrefix(line, "@") {
		_, rest, ok := strings.Cut(line, " ")
		if !ok {
			return Message{}, false
		}
		line = rest
	}
	if !strings.HasPrefix(line, ":") {
		return Message{}, false
	}
	prefix, rest, ok := strings.Cut(line[1:], " ")
	if !ok {
		return Message{}, false
	}
	verb, params, ok := strings.Cut(rest, " ")
	if !ok || verb != "PRIVMSG" {
		return Message{}, false
	}
	target, text, ok := strings.Cut(params, " :")
	if !ok || !strings.HasPrefix(target, "#") {
		return Message{}, false
	}

	user, _, _ := strings.Cut(prefix, "!")
	msg := Message{
		Username: user,
		Channel:  strings.TrimPrefix(target, "#"),
		Text:     text,
	}
	if strings.HasPrefix(text, "!") {
		cmd, _, _ := strings.Cut(text, " ")
		msg.Command = cmd
		msg.Text = strings.TrimPrefix(text, cmd)
	}
	return msg, true
}

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
