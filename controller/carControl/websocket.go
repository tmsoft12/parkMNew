package carcontrol

import (
	modelscar "park/models/modelsCar"

	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan modelscar.Car_Model)

func Ws(c *websocket.Conn) {
	defer c.Close()
	clients[c] = true

	for {
		var car modelscar.Car_Model
		if err := c.ReadJSON(&car); err != nil {
			delete(clients, c)
			break
		}
		broadcast <- car
	}
}

func HandleMessages() {
	for {
		car := <-broadcast
		for client := range clients {
			if err := client.WriteJSON(car); err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
