package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type ActionMessage struct {
	Misc string
	Key  string
}

type Client struct {
	Board      *Board
	Connection *websocket.Conn
	Close      chan int
}

func (C *Client) sendBoard() {
	if C.Connection == nil {
		return
	}
	C.Connection.WriteJSON(C.Board)
}

func (C *Client) gameloop() {
	for {
		if C.Connection == nil {
			return
		}

		time.Sleep(1 * time.Second)
		C.Board.Lock()
		C.Board.tick()
		C.sendBoard()
		C.Board.Unlock()
	}
}
func (C *Client) communicationLoop() {
	for {
		_, reader, err := C.Connection.NextReader()

		if err != nil {
			C.Connection.Close()
			C.Connection = nil
			return
		}

		p := make([]byte, 10000)
		n, _ := reader.Read(p)
		var action ActionMessage
		json.Unmarshal(p[:n], &action)

		if C.Board.fallingPiece == nil {
			continue
		}

		switch action.Key {
		case "KeyD":
			C.Board.move(1)
			C.Board.drawFalling(C.Board.fallingPiece.PositionX-1, C.Board.fallingPiece.PositionY, make([][]int, 0))
			C.sendBoard()
			break
		case "KeyA":
			C.Board.move(0)
			C.Board.drawFalling(C.Board.fallingPiece.PositionX+1, C.Board.fallingPiece.PositionY, make([][]int, 0))
			C.sendBoard()
			break
		case "KeyZ":
			preX := C.Board.fallingPiece.PositionX
			preY := C.Board.fallingPiece.PositionY
			preGrid := C.Board.fallingPiece.Grid
			C.Board.flip(0)
			C.Board.drawFalling(preX, preY, preGrid)
			C.sendBoard()
			break
		case "KeyX":
			preX := C.Board.fallingPiece.PositionX
			preY := C.Board.fallingPiece.PositionY
			preGrid := C.Board.fallingPiece.Grid
			C.Board.flip(1)
			C.Board.drawFalling(preX, preY, preGrid)
			C.sendBoard()
			break
		case "KeyS":
			C.Board.drop()
			C.Board.drawFalling(-1, -1, [][]int{})
			C.sendBoard()
			break
		}

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
			Grid: initiateGrid(20, 10),
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

	go client.communicationLoop()
	go client.gameloop()

}
