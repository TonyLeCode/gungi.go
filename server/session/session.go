package session

import "github.com/TonyLeCode/gungi.go/server/gungi"

// var rooms = make(map[string][]*websocket.Conn)

type Users struct {
	Id    string
	Color int
}

type Session struct {
	GameId     string
	BoardState gungi.Board
	History    []string
	Users      Users
}

type Sessions map[string]Session
