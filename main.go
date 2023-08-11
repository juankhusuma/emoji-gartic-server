package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	ws "github.com/gofiber/websocket/v2"
)

var hub = Hub{
	clients:   make(map[*Client]bool),
	broadcast: make(chan []byte),
	register:  make(chan *Client),
}

func main() {
	go hub.Run()
	app := fiber.New()
	app.Use("/ws", func(c *fiber.Ctx) error {
		if ws.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", ws.New(func(c *ws.Conn) {
		client := &Client{hub: &hub, conn: c, send: make(chan []byte, 256)}
		hub.register <- client
		for {
			if _, msg, err := c.ReadMessage(); err == nil {
				fmt.Println(string(msg))
				for client := range hub.clients {
					fmt.Println(client)
					err := client.conn.WriteMessage(ws.TextMessage, []byte("pong "+string(msg)))
					if err != nil {
						fmt.Println(err)
						break
					}
				}
			} else {
				fmt.Println(err)
				break
			}
		}
		defer c.Close()
	}))

	app.Listen(":8000")
}
