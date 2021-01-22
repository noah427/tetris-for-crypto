package main

import (
	"fmt"
	"math/rand"
	"sync"
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
	Inverted  bool
	Rotated   bool
	Direction int
}

type Board struct {
	sync.Mutex
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
		B.drawFalling(B.fallingPiece.PositionX, B.fallingPiece.PositionY-1, make([][]int, 0))
	}
}

func (B *Board) drawFalling(preX int, preY int, preGrid [][]int) {
	if preX != -1 && preY != -1 {
		if len(preGrid) != 0 {
			B.prettyPrint(preGrid)
			for i, y := range preGrid {
				for j, x := range y {
					if x != 0 {
						if preY+i > 19 || preX+j > 9 {
							fmt.Println(B.prettyPrint(preGrid))
						} else {
							B.Grid[preY+i][preX+j] = 0
						}

					}
				}
			}
		} else {
			for i, y := range B.fallingPiece.Grid {
				for j, x := range y {
					if x != 0 {
						B.Grid[preY+i][preX+j] = 0
					}
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

func (B *Board) collision() bool {
	// bottom of piece
	pBottom := B.fallingPiece.PositionY + len(B.fallingPiece.Grid)
	// check if it hit the bottom
	if pBottom == 20 {
		B.fallingPiece.Collided = true
		return true
	}

	for i, y := range B.fallingPiece.Grid {
		for j, x := range y {
			// for every block in the grid
			coord := []int{B.fallingPiece.PositionX + j, B.fallingPiece.PositionY + i}

			if x != 0 { // if the piece isn't empty
				if B.Grid[coord[1]+1][coord[0]] != 0 { // if the piece directly under it isn't empty
					return true
					B.fallingPiece.Collided = true
				}
			}
		}
	}

	return false
}

func (B *Board) placePiece() {
	B.drawFalling(-1, -1, make([][]int, 0))
	B.fallingPiece = nil
}

func (B *Board) move(direction int) { // 0 left 1 right

	if direction == 0 {
		if B.fallingPiece.PositionX == 0 {
			return
		}

		B.fallingPiece.PositionX--
	} else {
		if B.fallingPiece.PositionX+len(B.fallingPiece.Grid[0]) == 10 {
			return
		}
		B.fallingPiece.PositionX++
	}
}

func (B *Board) flip(d int) { // 0 cc 1 c
	if d == 0 {
		switch B.fallingPiece.Direction {
		case 0:
			B.rotate()
			B.fallingPiece.Direction++
			break
		case 1:
			B.rotate()
			B.invert()
			B.fallingPiece.Direction++
			break
		case 2:
			B.invert()
			B.fallingPiece.Direction++
			break
		case 3:
			B.rotate()
			B.invert()
			B.fallingPiece.Direction = 0
			break

		}

	} else {
		switch B.fallingPiece.Direction {
		case 0:
			B.rotate()
			B.invert()
			B.fallingPiece.Direction = 3
			break
		case 1:
			B.rotate()
			B.fallingPiece.Direction--
			break
		case 2:
			B.invert()
			B.rotate()
			B.fallingPiece.Direction--
			break
		case 3:
			B.rotate()
			B.fallingPiece.Direction--
			break

		}
	}
}

func (B *Board) invert() {

	if B.fallingPiece.Rotated {
		for i, _ := range B.Grid {
			swap := B.Grid[i][0]
			B.Grid[i][0] = B.Grid[i][1]
			B.Grid[i][1] = swap
		}
	} else {
		swap := B.fallingPiece.Grid[0]
		B.fallingPiece.Grid[0] = B.fallingPiece.Grid[1]
		B.fallingPiece.Grid[1] = swap
	}

	if B.fallingPiece.Inverted {
		B.fallingPiece.Inverted = false
	} else {
		B.fallingPiece.Inverted = true
	}
}

func (B *Board) drop() {
	for {
		B.fallingPiece.PositionY++

		if B.collision() {
			B.placePiece()
			return
		}
	}
}

func (B *Board) rotate() {
	newGrid := initiateGrid(len(B.fallingPiece.Grid[0]), len(B.fallingPiece.Grid))
	for i, y := range B.fallingPiece.Grid {
		for j, x := range y {
			if x != 0 {
				newGrid[j][i] = x
			}
		}
	}
	B.fallingPiece.Grid = newGrid
	if B.fallingPiece.Rotated {
		B.fallingPiece.Rotated = false
	} else {
		B.fallingPiece.Rotated = true
	}
}

func (B *Board) prettyPrint(opt [][]int) string {
	var output string

	if len(opt) != 0 {
		for _, x := range opt {
			for j, y := range x {
				if j == len(x)-1 {
					output += fmt.Sprintf("|%d|\n", y)
					output += "---------------------\n"
				} else {
					output += fmt.Sprintf("|%d", y)
				}
			}

		}
	} else {
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
	}

	return output
}
