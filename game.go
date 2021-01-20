package main

import (
	"fmt"
	"math/rand"
)

var piecesList Pieces

type Pieces struct {
	Pieces [][][]int
}

type Piece struct {
	Grid      [][]int
	PositionX int
	PositionY int
	Collided  bool
}

type Board struct {
	fallingPiece *Piece
	Grid         [][]int
}

func (B *Board) spawnPiece() {
	piece := Piece{
		Grid:      piecesList.Pieces[rand.Intn(len(piecesList.Pieces))],
		PositionX: 5,
		PositionY: 0,
	}

	B.fallingPiece = &piece

}

func (B *Board) tick() {

	if B.fallingPiece == nil {
		B.spawnPiece()
	} else if B.fallingPiece.PositionY == 20-len(B.fallingPiece.Grid) {
		B.placePiece()
	} else if B.fallingPiece.Collided {
		B.placePiece()
	} else {
		B.fallingPiece.PositionY++
		B.collision()
		B.drawFalling(B.fallingPiece.PositionX, B.fallingPiece.PositionY-1)
	}
}

func (B *Board) drawFalling(preX int, preY int) {
	if preX != -1 && preY != -1 {
		for i, y := range B.fallingPiece.Grid {
			for j, x := range y {
				if x != 0 {
					B.Grid[preY+i][preX+j] = 0
				}
			}
		}
	}

	for i, y := range B.fallingPiece.Grid {
		for j, x := range y {
			if x != 0 {
				if B.fallingPiece.PositionX+j < 10 && B.fallingPiece.PositionY+i < 20 {
					B.Grid[B.fallingPiece.PositionY+i][B.fallingPiece.PositionX+j] = x
				}
			}
		}
	}
}

func (B *Board) collision() {
	for i, y := range B.fallingPiece.Grid {
		for j, x := range y {
			if x != 0 {
				if B.fallingPiece.PositionX+j < 10 && B.fallingPiece.PositionY+i+1 < 20 {
					if B.Grid[B.fallingPiece.PositionY+i+1][B.fallingPiece.PositionX+j] != 0 {

						B.fallingPiece.Collided = true
					}
				}
			}
		}
	}
}

func (B *Board) placePiece() {
	B.drawFalling(-1, -1)
	B.fallingPiece = nil
}

func (B *Board) move(direction int) { // 0 left 1 right
	if direction == 0 {
		if B.fallingPiece.PositionX == 0 {
			return
		}
		B.fallingPiece.PositionX--
	} else {
		if B.fallingPiece.PositionX == 10 {
			return
		}
		B.fallingPiece.PositionX++
	}
}

func (B *Board) prettyPrint() string {
	var output string
	for _, x := range B.Grid {
		for j, y := range x {
			if j == len(x)-1 {
				output += fmt.Sprintf("|%d|\n", y)
				output += "---------------------\n"
			} else {
				output += fmt.Sprintf("|%d", y)
			}
		}

	}
	return output
}
