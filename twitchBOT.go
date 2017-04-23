package main

import (
	// Includes padrões pro Bot
	"fmt"
	"os"
	"log"
	"regexp"
	"bufio"
	"strings"
	"./ircBot"
	// Includes dos comandos criados
	"io/ioutil"
	"net/http"
	"time"
)

func main(){

	// "Try catch" do Go
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()

	// Olhar na função do Package "ircBot"
	// Cria o bot
	bot := ircBot.NewBot(
		"irc.chat.twitch.tv", // Host (default)
		"6667", // Porta (default)
		"nick-here", // Nick do Bot
		"channel-here", // Canal para se auto conectar
		"oauth:token-here", // Token do Bot
	)

	// Conecta o Bot
	conn, _ := bot.Connect()
	defer conn.Close()

	// Cria o buffer
	bot.Maker()

	// Inicia a função "PingLoop()" em uma nova thread
	go bot.PingLoop()

	// Expressão Regular default para tratamento da resposta do servidor
	reg, _ := regexp.Compile(":([a-z0-9_]*)!([a-z0-9_]*)@([a-z0-9.-]*) ([A-Z]*) #([a-z0-9_]*) :(![a-z0-9_]*)?(.*)")
	for {
		// Recebe os dados do servidor
		line, err := bot.Buffer.ReadLine()
		if err != nil {
			log.Fatal("Unable to receive data from IRC server ", err)
			os.Exit(1)
		}

		// Função para segunda checagem do "Eu to vivo"
		isPing, _ := regexp.MatchString("PING", line)
		if isPing == true {
			data := strings.Split(line, "PING ")
			bot.Send("PONG "+data[1]+"")
			time.Sleep(50 * time.Millisecond)
			continue
		}

		// Aplica a Expressão Regular
		t := reg.FindStringSubmatch(line)
		// Se a Expressão Regular retornar 8 argumentos
		if len(t) >= 8 {
			username := t[1]
			channel := t[5]
			command := t[6]
			text := t[7]

			// Comandos
			if command == "!gold" {
				request, _ := http.Get("https://wowtoken.info/snapshot.json")
				r_body, _ := ioutil.ReadAll(request.Body)
				request.Body.Close()
				body := string(r_body[:])
				reg2, _ := regexp.Compile("{\"NA\":{\"timestamp\":([0-9]*),\"raw\":{\"buy\":([0-9]*),\"")
				json := reg2.FindStringSubmatch(body)
				gold_ := json[2]
				gold := gold_[:len(gold_)-3] + "," + gold_[len(gold_)-3:]
				// fmt.Printf("%s\n", json[2])
				bot.Msg(channel, "[Bolsa de Azeroth] informa:")
				time.Sleep(1500 * time.Millisecond)
				bot.Msg(channel, "[NA] Cotação do OURO: " + gold + "g")
				continue
			}
			if command == "!dolar" {
				// bot.Send("PRIVMSG #"+ t[4] +" : Obtendo cotação, aguarde...")
				res, _ := http.Get("http://api.dolarhoje.com")
				body, _ := ioutil.ReadAll(res.Body)
				time.Sleep(1000 * time.Millisecond)
				bot.Msg(channel, "[DOLAR hoje] R$ " + string(body))
				continue
			}

			// Imprime a mensagem formatada recebida de um canal
			// 	"[#canal] Username: comando/mensagem"
			fmt.Printf("[#%s] %s: %s%s\n", channel, username, command, text)
		} else {
			// Impreme outras mensagens recebidas do servidor que
			//	não sejam referentes à mensagens de usuários
			fmt.Printf("%s\n", line)
		}

	}

	// Bye bye
	os.Exit(0)
}