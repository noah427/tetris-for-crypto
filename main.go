package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	roomManager *RoomManager
)

func main() {

	loadPiecesJSON()

	roomManager = &RoomManager{
		Rooms: map[string]Room{},
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Static("/", "./public/home")
	app.Static("/images", "./public/images")
	app.Static("/js", "./public/js")
	app.Static("/css", "./public/css")

	app.Get("/createRoom", createRoom)
	app.Get("/rooms/:room", joinRoom)
	app.Get("/rooms/:room/ws", websocket.New(routeRoomConnections))

	app.Listen(":3000")
}
