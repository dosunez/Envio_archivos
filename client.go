package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	transport "github.com/graarh/golang-socketio/transport"
)

type Message struct {
	Channel string `json:"channel"`
	File    []byte `json:"file"`
	Name    string `json:"name"`
	Sender  string `json:"sender"`
}

var clientName string
var joinedRooms []string

func alreadyJoined(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func createMenu(options []string) {
	for i := 0; i < len(options); i++ {
		option := strings.ReplaceAll(options[i], `"`, "")
		println(strconv.Itoa(i+1) + " - " + option)
	}
}

func showChannels(c *gosocketio.Client) (string, error) {
	result, err := c.Ack("/room-list", "rooms", time.Second*10)
	if err != nil {
		println(err)
		return "", err
	}
	options := strings.Split(result, ",")
	var channelIndex int
	for i := 0; i < len(options); i++ {
		channelName := strings.ReplaceAll(options[i], `"`, "")
		if !alreadyJoined(joinedRooms, channelName) {
			println(strconv.Itoa(i+1) + " - " + channelName)
		}
	}
	fmt.Scanln(&channelIndex)
	selectedOption := strings.ReplaceAll(options[channelIndex-1], `"`, "")
	return selectedOption, nil
}

func showMenu(c *gosocketio.Client) {
	for {
		println("")
		options := []string{"Crear canal", "Conectarse a canal", "Enviar archivo", "Salir"}
		createMenu(options)
		var option int
		fmt.Scanln(&option)

		if option == 4 {
			break
		}

		switch option {
		case 1:

			println("Ingrese el nombre del canal")
			var name string
			fmt.Scanln(&name)
			if alreadyJoined(joinedRooms, name) {
				println("Ya se encuentra suscrito a este canal")
				break
			}
			_, err := c.Ack("/join", name, time.Second*5)
			if err != nil {
				println(err)
				break
			}
			joinedRooms = append(joinedRooms, name)
			println("Se ha creado el canal exitosamente")

			break

		case 2:
			println("Seleccione el canal al que quiere unirse")
			selectedOption, err := showChannels(c)
			if err != nil {
				println(err)
				break
			}
			if alreadyJoined(joinedRooms, selectedOption) {
				println("Ya se encuentra en este canal")
				break
			}
			joinedRooms = append(joinedRooms, selectedOption)
			c.Ack("/join", selectedOption, time.Second*5)
			break
		case 3:
			println("Seleccione el canal a cual quiere enviar el archivo")
			createMenu(joinedRooms)
			var index int
			fmt.Scanln(&index)
			selectedOption := joinedRooms[index-1]
			println("Ingrese la ruta del archivo")
			var fileName string
			fmt.Scanln(&fileName)
			file, _ := ioutil.ReadFile(fileName)
			name := strings.Split(fileName, "/")
			c.Emit("/file", Message{selectedOption, file, name[len(name)-1], clientName})
			println("archivo enviado exitosamente")
			println("")
			break

		default:
			println("seleccion invalida")
		}

	}
}

func main() {
	println("Ingrese su nombre")
	fmt.Scanln(&clientName)

	c, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 3811, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	c.On("/file", func(h *gosocketio.Channel, args Message) {
		if clientName != args.Sender {
			route := "./" + clientName + "/" + args.Name
			os.Mkdir(clientName, 0644)
			println(route)
			ioutil.WriteFile(route, args.File, 0644)
		}
	})
	showMenu(c)
}
