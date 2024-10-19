package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type Flags struct {
	Host     string
	Token    string
	LogLevel string
}

func parseFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.Host, "host", "192.168.66.1:9999", "ShellCrash address (Required)")
	flag.StringVar(&flags.Token, "token", "", "Crash token")
	flag.StringVar(&flags.LogLevel, "log-level", "info", "LogLevel: debug|info|warning|error. Default: info.")
	flag.Parse()
	return flags
}

func WsLogs(args *Flags) <-chan string {
	url := fmt.Sprintf("ws://%s/logs?token=%s&level=%s", args.Host, args.Token, args.LogLevel)
	log.Printf("connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	// defer conn.Close()
	conn.SetReadLimit(10 * 1024 * 1024)

	ch := make(chan string, 100)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			ch <- string(message)
		}
	}()
	return ch
}

func WsConnections(args *Flags) <-chan string {
	url := fmt.Sprintf("ws://%s/connections?token=%s", args.Host, args.Token)
	log.Printf("connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	// defer conn.Close()
	conn.SetReadLimit(10 * 1024 * 1024)

	ch := make(chan string, 100)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			ch <- string(message)
		}
	}()
	return ch
}

// ws://192.168.66.1:9999/traffic?token=

func main() {
	args := parseFlags()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	chLog := WsLogs(&args)
	chConn := WsConnections(&args)

	for {
		select {
		case s := <-chLog:
			log.Printf("recv: %s", s)
		case s := <-chConn:
			log.Printf("recv: %s", s)
		case <-interrupt:
			log.Println("interrupt")
			return
		}
	}
}
