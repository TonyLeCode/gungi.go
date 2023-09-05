package session

// var rooms = make(map[string][]*websocket.Conn)

type Users struct {
	Id    string
	Color int
}

type Session struct {
	GameId     string
	BoardState string
	History    string
	Users      Users
}

type Sessions map[string]Session
