package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	gosocketio "github.com/graarh/golang-socketio"
	transport "github.com/graarh/golang-socketio/transport"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Channel string `json:"channel"`
	File    []byte `json:"file"`
	Name    string `json:"name"`
	Sender  string `json:"sender"`
}

func roomExists(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	var rooms []string

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {

		println("connected " + strconv.FormatInt(server.AmountOfSids(), 10))
	})

	server.On("/join", func(c *gosocketio.Channel, channel string) string {
		if !roomExists(rooms, channel) {
			rooms = append(rooms, channel)
		}
		c.Join(channel)
		return "joined to " + channel
	})

	server.On("/file", func(c *gosocketio.Channel, m Message) {
		server.BroadcastTo(m.Channel, "/file", m)
	})

	server.On("/room-list", func(c *gosocketio.Channel) string {
		return strings.Join(rooms, ",")
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Println("Starting server...")
	log.Panic(http.ListenAndServe(":3811", serveMux))
}
