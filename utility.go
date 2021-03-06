package main

import (
	"math/rand"

	"encoding/json"
	"io/ioutil"

	"strings"
	"time"
)

func loadPiecesJSON() {
	data, _ := ioutil.ReadFile("pieces.json")
	json.Unmarshal(data, &piecesList)
}

func generateRoomName() string {
	rand.Seed(time.Now().Unix())
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var output strings.Builder
	for i := 0; i < 6; i++ {
		char := charset[rand.Intn(len(charset))]
		output.WriteString(string(char))
	}

	return output.String()
}

func initiateGrid(y,x int) [][]int {
	a := make([][]int, y)
	for i := range a {
		a[i] = make([]int, x)
	}
	return a
}
