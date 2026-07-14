package ircbot

import (
	"bufio"
	"net"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		line string
		want Message
		ok   bool
	}{
		{
			name: "command",
			line: ":k_u_p!k_u_p@k_u_p.tmi.twitch.tv PRIVMSG #canal :!dolar agora",
			want: Message{Username: "k_u_p", Channel: "canal", Command: "!dolar", Text: " agora"},
			ok:   true,
		},
		{
			name: "chat message",
			line: ":user123!user123@user123.tmi.twitch.tv PRIVMSG #canal :oi galera",
			want: Message{Username: "user123", Channel: "canal", Text: "oi galera"},
			ok:   true,
		},
		{
			name: "with IRCv3 tags",
			line: "@badges=broadcaster/1;color=#FF0000 :dono!dono@dono.tmi.twitch.tv PRIVMSG #canal :!dolar",
			want: Message{Username: "dono", Channel: "canal", Command: "!dolar"},
			ok:   true,
		},
		{
			name: "server notice",
			line: ":tmi.twitch.tv 001 nick :Welcome, GLHF!",
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Parse(tt.line)
			if ok != tt.ok || got != tt.want {
				t.Errorf("Parse(%q) = %+v, %v; want %+v, %v", tt.line, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestMsgStripsCRLF(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	b := &Bot{conn: client}
	go func() {
		if err := b.Msg("canal", "preço\r\nQUIT :injetado"); err != nil {
			t.Error(err)
		}
	}()

	got, err := bufio.NewReader(server).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	want := "PRIVMSG #canal :preçoQUIT :injetado\r\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
