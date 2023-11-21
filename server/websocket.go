package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/TonyLeCode/gungi.go/server/api"
	"github.com/TonyLeCode/gungi.go/server/auth"
	db "github.com/TonyLeCode/gungi.go/server/db/sqlc"
	"github.com/TonyLeCode/gungi.go/server/gungi"
	"github.com/TonyLeCode/gungi.go/server/gungi/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
)

// type UserConn struct {
// 	Username string
// 	Conn     *melody.Session
// }

type PlayConnections struct {
	connections map[*melody.Session]bool
	mu          sync.RWMutex
}
type GameConnections struct {
	connections map[uuid.UUID]map[*melody.Session]bool
	mu          sync.RWMutex
}

func (pc *PlayConnections) AddConnection(s *melody.Session) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.connections[s] = true
}
func (pc *PlayConnections) RemoveConnection(s *melody.Session) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	delete(pc.connections, s)
}

func (gc *GameConnections) AddConnection(gameID uuid.UUID, s *melody.Session) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	if _, ok := gc.connections[gameID]; !ok {
		gc.connections[gameID] = make(map[*melody.Session]bool)
	}

	gc.connections[gameID][s] = true
}
func (gc *GameConnections) RemoveConnection(gameID uuid.UUID, s *melody.Session) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	delete(gc.connections[gameID], s)
}

// type ClientList struct {
// 	clients []*melody.Session
// 	mu      sync.RWMutex
// }

// This won't scale, find solution in future
func RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

type client struct {
	username string
	route    string
	gameID   string
}

type clientList struct {
	clients map[*melody.Session]client
	mutex   sync.Mutex
}

func (c *clientList) AddClient(s *melody.Session, username string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.clients[s] = client{username: username}
}

func (c *clientList) RemoveClient(s *melody.Session) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.clients, s)
}

func (c *clientList) ChangeRoute(s *melody.Session, route string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, ok := c.clients[s]; ok {
		entry.route = route
		c.clients[s] = entry
	} else {
		return errors.New("unable to change route")
	}
	return nil
}

type Room struct {
	clients map[*melody.Session]bool
	mutex   sync.Mutex
}

func (r *Room) AddClient(s *melody.Session) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.clients[s] = true
}

func (r *Room) RemoveClient(s *melody.Session) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.clients, s)
}

type PlayerMove struct {
	Type string `json:"type"`
	Data struct {
		MoveType  int `json:"moveType"`
		FromCoord int `json:"fromCoord"`
		ToCoord   int `json:"toCoord"`
	}
}

type GameRoomType struct {
	Host        string `json:"host"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Rules       string `json:"rules"`
}

type CreateGameRoomRequest struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Rules       string `json:"rules"`
}

type SerializedGameRoom struct {
	GameRoomType
	RoomID uuid.UUID `json:"roomid"`
}

type MsgPayload struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MsgResponse struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type makeGameMoveMsg struct {
	GameID    string `json:"gameID"`
	FromPiece int    `json:"fromPiece"`
	FromCoord int    `json:"fromCoord"`
	MoveType  int    `json:"moveType"`
	ToCoord   int    `json:"toCoord"`
}

// func RedisRoomExists(r *redis.Client, roomKey string) error {
// 	ctx := context.Background()
// 	exists, err := r.Exists(ctx, roomKey).Result()
// 	if err != nil {
// 		return err
// 	}

// 	if exists == 0 {
// 		_, err = r.Do(ctx, "JSON.SET", roomKey, "$", "[]").Result()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func ws2(m *melody.Melody, dbConn *api.DBConn) echo.HandlerFunc {
	PlayConnections := PlayConnections{
		connections: make(map[*melody.Session]bool),
		mu:          sync.RWMutex{},
	}
	GameConnections := GameConnections{
		connections: make(map[uuid.UUID]map[*melody.Session]bool),
		mu:          sync.RWMutex{},
	}

	m.HandleConnect(func(s *melody.Session) {
		log.Println("Connected2", s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		if gameID, ok := s.Get("Game"); ok {
			if gameUUID, err := uuid.Parse(gameID.(string)); err == nil {
				log.Println("Removed Game Connection")
				GameConnections.RemoveConnection(gameUUID, s)
			}
		} else if _, ok := PlayConnections.connections[s]; ok {
			log.Println("Removed Play Connection")
			PlayConnections.RemoveConnection(s)
		}
		log.Println("Disconnect")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		// ctx := context.Background()

		var unmarshal MsgPayload
		err := json.Unmarshal(msg, &unmarshal)
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		switch unmarshal.Type {
		case "auth":
			log.Println("auth1")
			var token string
			err := json.Unmarshal(unmarshal.Payload, &token)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			claims, err := auth.AuthenticateSupabaseToken(token)
			if err != nil {
				log.Println(err)
				// 0 means failed
				authResponse := MsgResponse{
					Type:    "auth",
					Payload: "0",
				}
				payload, err := json.Marshal(authResponse)
				if err != nil {
					log.Println("Error: ", err)
				} else {
					err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
						return q == s
					})
					if err != nil {
						log.Println("Error: ", err)
					}
				}
				return
			}

			metadata := claims["user_metadata"].(map[string]interface{})
			username := metadata["username"].(string)
			s.Set("username", username)

			// 1 means authenticated
			authResponse := MsgResponse{
				Type:    "auth",
				Payload: "1",
			}
			payload, err := json.Marshal(authResponse)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			log.Println("auth2")
			err = s.Write(payload)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			return

		case "joinGame":
			//deletes previous game to prevent memory leak
			if gameID, ok := s.Get("Game"); ok {
				if gameUUID, err := uuid.Parse(gameID.(string)); err == nil {
					log.Println("Removed Game Connection")
					GameConnections.RemoveConnection(gameUUID, s)
				}
			}

			var gameID string
			err = json.Unmarshal(unmarshal.Payload, &gameID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if gameUUID, err := uuid.Parse(gameID); err == nil {
				GameConnections.AddConnection(gameUUID, s)
				log.Println("Joined Game: ", gameUUID)
				log.Println("Conns: ", GameConnections.connections)
				s.Set("Game", gameID)
			}

		case "leaveGame":
			var gameID string
			err = json.Unmarshal(unmarshal.Payload, &gameID)
			log.Println("gameID", gameID)
			log.Println("gameID", unmarshal.Payload)
			log.Println("gameID", unmarshal)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if gameUUID, err := uuid.Parse(gameID); err == nil {
				log.Println("Left Game: ", gameUUID)
				GameConnections.RemoveConnection(gameUUID, s)
				log.Println("Conns: ", GameConnections.connections)
			}

		case "makeGameMove":
			ctx := context.Background()
			var makeGameMoveMsg makeGameMoveMsg
			err = json.Unmarshal(unmarshal.Payload, &makeGameMoveMsg)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			game, err := getGame(dbConn, makeGameMoveMsg.GameID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			board := gungi.CreateBoard(game.Ruleset)
			err = board.SetBoardFromFen(game.CurrentState)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			board.SetHistory(strings.Fields(game.History.String))
			board.PrintBoard()

			err = board.MakeMove(makeGameMoveMsg.FromPiece, board.ConvertInputCoord(makeGameMoveMsg.FromCoord), makeGameMoveMsg.MoveType, board.ConvertInputCoord(makeGameMoveMsg.ToCoord))
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			board.PrintBoard()
			log.Println("fen: ", board.BoardToFen())

			game.CurrentState = board.BoardToFen()
			game.History.String = board.SerializeHistory()
			game.History.Valid = true

			makeMoveParams := db.MakeMoveParams{
				ID:           game.ID,
				CurrentState: game.CurrentState,
				History:      game.History,
			}
			queries := db.New(dbConn.Conn)
			err = queries.MakeMove(ctx, makeMoveParams)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			_, legalMoves := board.GetLegalMoves()
			correctedLegalMoves := make(map[int][]int)
			for key, element := range legalMoves {
				correctedKey := board.ConvertOutputCoord(key)
				for _, index := range element {
					correctedElement := board.ConvertOutputCoord(index)
					correctedLegalMoves[correctedKey] = append(correctedLegalMoves[correctedKey], correctedElement)
				}
			}
			game.MoveList = correctedLegalMoves
			log.Println(correctedLegalMoves)

			msgResponse := MsgResponse{
				Type:    "game",
				Payload: game,
			}
			marshalledResponse, err := json.Marshal(msgResponse)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			var keys []*melody.Session
			for key := range GameConnections.connections[game.ID] {
				keys = append(keys, key)
			}

			m.BroadcastMultiple(marshalledResponse, keys)

		//TODO finish cases
		case "gameResign":
		case "requestGameUndo":
		case "responseGameUndo":

		case "joinPlay":
			PlayConnections.AddConnection(s)
			log.Println("Joined Play")
			log.Println("Conns: ", PlayConnections.connections)
			// TODO host id to username
			roomList, err := dbConn.GetRoomList()
			if err != nil {
				log.Println(err)
				return
			}
			roomListMsg := MsgResponse{
				Type:    "roomList",
				Payload: roomList,
			}
			marshaledRoomListMsg, err := json.Marshal(roomListMsg)
			if err != nil {
				log.Println(err)
				return
			}
			s.Write(marshaledRoomListMsg)

		case "leavePlay":
			log.Println("Left Play")
			PlayConnections.RemoveConnection(s)
			log.Println("Conns: ", PlayConnections.connections)

		case "createPlayRoom":
			var roomPayload CreateGameRoomRequest
			err := json.Unmarshal(unmarshal.Payload, &roomPayload)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			username, exists := s.Get("username")
			if !exists {
				return
			}

			err = dbConn.CreateRoom(username.(string), roomPayload.Description, roomPayload.Rules, roomPayload.Type, roomPayload.Color)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			roomList, err := dbConn.GetRoomList()
			if err != nil {
				log.Println(err)
				return
			}
			roomListMsg := MsgResponse{
				Type:    "roomList",
				Payload: roomList,
			}
			marshaledRoomListMsg, err := json.Marshal(roomListMsg)
			if err != nil {
				log.Println(err)
				return
			}

			err = m.BroadcastFilter(marshaledRoomListMsg, func(q *melody.Session) bool {
				return PlayConnections.connections[q]
			})
			if err != nil {
				log.Println(err)
				return
			}

		case "acceptPlayRoom":
			username, exists := s.Get("username")
			if !exists {
				return
			}

			var roomid string
			err := json.Unmarshal(unmarshal.Payload, &roomid)
			if err != nil {
				log.Println(err)
				return
			}

			roomUUID, err := uuid.Parse(roomid)
			if err != nil {
				log.Println(err)
				return
			}

			//TODO db transaction
			room, err := dbConn.DeleteRoom(roomUUID)
			if err != nil {
				log.Println(err)
				return
			}

			gameID, err := dbConn.CreateGameFromRoom(room, username.(string))
			if err != nil {
				log.Println(err)
				return
			}
			gameAcceptedMsg := MsgResponse{
				Type:    "roomAccepted",
				Payload: gameID,
			}
			marshaledGameAcceptedMsg, err := json.Marshal(gameAcceptedMsg)
			if err != nil {
				log.Println(err)
				return
			}

			m.BroadcastFilter(marshaledGameAcceptedMsg, func(q *melody.Session) bool {
				qUsername, exists := q.Get("username")
				if !exists {
					return false
				}
				return qUsername == room.Host || qUsername == username
			})

			roomList, err := dbConn.GetRoomList()
			if err != nil {
				log.Println(err)
				return
			}
			updatedRoomListMsg := MsgResponse{
				Type:    "roomList",
				Payload: roomList,
			}
			marshaledUpdatedRoomListMsg, err := json.Marshal(updatedRoomListMsg)
			if err != nil {
				log.Println(err)
				return
			}

			var keys []*melody.Session
			for key := range PlayConnections.connections {
				keys = append(keys, key)
			}

			m.BroadcastMultiple(marshaledUpdatedRoomListMsg, keys)

		case "cancelPlayRoom":
			var roomid string
			err := json.Unmarshal(unmarshal.Payload, &roomid)
			if err != nil {
				log.Println(err)
				return
			}

			roomUUID, err := uuid.Parse(roomid)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = dbConn.DeleteRoom(roomUUID)
			if err != nil {
				log.Println(err)
				return
			}

			roomList, err := dbConn.GetRoomList()
			if err != nil {
				log.Println(err)
				return
			}
			updatedRoomListMsg := MsgResponse{
				Type:    "roomList",
				Payload: roomList,
			}
			marshaledUpdatedRoomListMsg, err := json.Marshal(updatedRoomListMsg)
			if err != nil {
				log.Println(err)
				return
			}

			var keys []*melody.Session
			for key := range PlayConnections.connections {
				keys = append(keys, key)
			}

			m.BroadcastMultiple(marshaledUpdatedRoomListMsg, keys)
		}
	})

	return func(c echo.Context) error {
		err := m.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			return err
		}
		return nil
	}
}

func getGame(dbs *api.DBConn, roomID string) (api.GameWithMoves, error) {
	game, err := dbs.GetGame(roomID)
	if err != nil {
		log.Println("Error: ", err)
		return api.GameWithMoves{}, err
	}

	// msgResponse := MsgResponse{
	// 	Type:    "game",
	// 	Payload: game,
	// }
	return game, nil
}

type Move struct {
	FromPiece int `json:"fromPiece"`
	FromCoord int `json:"fromCoord"`
	MoveType  int `json:"moveType"`
	ToCoord   int `json:"toCoord"`
}

func handleGameMessages(msg MsgPayload, m *melody.Melody, s *melody.Session, dbs *api.DBConn, ctx context.Context, clientList *clientList, gameRooms *map[string]*Room, mutex *sync.Mutex) error {
	switch msg.Type {
	case "makeMove":
		gameID := clientList.clients[s].gameID
		if gameID == "" {
			return errors.New("could not find gameID")
		}

		game, err := getGame(dbs, gameID)
		if err != nil {
			return err
		}

		var move Move
		err = json.Unmarshal(msg.Payload, &move)
		if err != nil {
			return err
		}

		newBoard := gungi.CreateBoard("revised") //TODO change hardcode
		err = newBoard.SetBoardFromFen(game.CurrentState)
		if err != nil {
			return err
		}
		newBoard.SetHistory(strings.Fields(game.History.String))
		newBoard.PrintBoard()
		// log.Println(move)
		// log.Println("fen: ", newBoard.BoardToFen())
		// log.Println(newBoard.TurnColor)
		// log.Println(utils.GetColor(move.FromPiece))
		err = newBoard.MakeMove(move.FromPiece, utils.IndexToSquare(move.FromCoord), move.MoveType, utils.IndexToSquare(move.ToCoord))
		if err != nil {
			return err
		}
		newBoard.PrintBoard()
		log.Println("fen: ", newBoard.BoardToFen())

		game.CurrentState = newBoard.BoardToFen()
		game.History.String = newBoard.SerializeHistory()
		game.History.Valid = true

		makeMoveParams := db.MakeMoveParams{
			ID:           game.ID,
			CurrentState: game.CurrentState,
			History:      game.History,
		}
		queries := db.New(dbs.Conn)
		err = queries.MakeMove(ctx, makeMoveParams)
		if err != nil {
			return err
		}

		_, legalMoves := newBoard.GetLegalMoves()
		correctedLegalMoves := make(map[int][]int)
		for key, element := range legalMoves {
			correctedKey := utils.SquareToIndex(key)
			for _, index := range element {
				correctedElement := utils.SquareToIndex(index)
				correctedLegalMoves[correctedKey] = append(correctedLegalMoves[correctedKey], correctedElement)
			}
		}
		game.MoveList = correctedLegalMoves
		log.Println(correctedLegalMoves)

		msgResponse := MsgResponse{
			Type:    "game",
			Payload: game,
		}
		payload, err := json.Marshal(msgResponse)
		if err != nil {
			return err
		}
		for session := range (*gameRooms)[gameID].clients {
			session.Write(payload)
		}
	case "requestUndo":
		gameID := clientList.clients[s].gameID
		if gameID == "" {
			return errors.New("could not find gameID")
		}

		// gameUUID, err := uuid.Parse(gameID)
		// if err != nil {
		// 	return err
		// }

		// queryParams := db.CreateUndoParams{
		// 	GameID: gameUUID,
		// 	Color:  "w",
		// }

		// queries := db.New(dbs.PostgresDB)
		// undoID, err := queries.CreateUndo(ctx, queryParams)
		// if err != nil {
		// 	return err
		// }

		// msgResponse := MsgResponse{
		// 	Type:    "requestUndo",
		// 	Payload: undoID.String(),
		// }

		// payload, err := json.Marshal(msgResponse)
		// if err != nil {
		// 	return err
		// }

		// err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
		// 	// clientList.clients[q].username
		// 	// gameRooms[q]
		// 	return q == s //TODO only send to opponent
		// })
	case "responseUndo":
	case "resign":
	}
	return nil
}

// func RedisGetGameRooms(r *redis.Client) ([]SerializedGameRoom, error) {
// 	ctx := context.Background()
// 	err := RedisRoomExists(r, "game_rooms")
// 	if err != nil {
// 		log.Println("Error: ", err)
// 		return []SerializedGameRoom{}, err
// 	}

// 	val, err := r.Do(ctx, "JSON.GET", "game_rooms").Result()
// 	if err != nil {
// 		log.Println("Error: ", err)
// 		return []SerializedGameRoom{}, err
// 	}

// 	var roomList []SerializedGameRoom
// 	err = json.Unmarshal([]byte(val.(string)), &roomList)
// 	if err != nil {
// 		log.Println("Error: ", err)
// 		return []SerializedGameRoom{}, err
// 	}

// 	return roomList, nil
// }
