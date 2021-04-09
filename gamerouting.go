package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func createRoom(c *fiber.Ctx) error {
	roomManager.Lock()
	defer roomManager.Unlock()

	name := generateRoomName()
	room := Room{
		Name:    name,
		Clients: make(map[*websocket.Conn]Client),
	}
	roomManager.Rooms[name] = room
	c.Redirect(fmt.Sprintf("../rooms/%s", name))
	return nil
}

func joinRoom(c *fiber.Ctx) error {
	room := roomManager.Rooms[c.Params("room")]
	fmt.Println(room)
	c.SendFile("./public/html/tetris.html")
	return nil
}
