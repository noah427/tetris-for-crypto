package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/websocket"
)

type ActionMessage struct {
	Move string
}

type Client struct {
	Board      *Board
	Connection *websocket.Conn
	Close      chan int
}

func (C *Client) sendBoard() {
	C.Connection.WriteJSON(C.Board)
}

func (C *Client) recieveAction() {
	for {
		time.Sleep(1 * time.Second)
		C.Board.tick()
		C.sendBoard()
	}
}

type Room struct {
	sync.Mutex
	Name    string
	Clients map[*websocket.Conn]Client
}

type RoomManager struct {
	sync.Mutex
	Rooms map[string]Room
}

func routeRoomConnections(c *websocket.Conn) {
	roomManager.Lock()
	room := roomManager.Rooms[c.Params("room")]

	if room.Name == "" {
		fmt.Println("this closed")
		c.Close()
	}

	room.Lock()

	client := Client{
		Board: &Board{
			Grid: initiateGrid(),
		},
		Connection: c,
	}

	defer func() {
		roomManager.Unlock()
		room.Unlock()
		close := <-client.Close
		client.Connection.Close()

		fmt.Printf("Close status %d", close)
	}()

	room.Clients[c] = client

	client.sendBoard()

	go client.recieveAction()
}
