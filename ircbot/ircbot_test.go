package ircbot

import (
	"bufio"
	"net"
	"testing"
)

func TestMsgStripsCRLF(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	b := &Bot{conn: client}
	go b.Msg("canal", "preço\r\nQUIT :injetado")

	got, err := bufio.NewReader(server).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	want := "PRIVMSG #canal :preçoQUIT :injetado\r\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
