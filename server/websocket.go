package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/whitemonarch/gungi-server/server/api"
	"github.com/whitemonarch/gungi-server/server/auth"
	db "github.com/whitemonarch/gungi-server/server/db/sqlc"
	"github.com/whitemonarch/gungi-server/server/gungi"
)

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
	if len(gc.connections[gameID]) == 0 {
		delete(gc.connections, gameID)
	}
}

// This won't scale, find solution in future
func RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
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
		if gameID, ok := s.Get("game"); ok {
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
			id := claims["sub"].(string)
			s.Set("username", username)
			s.Set("id", id)

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
			err = s.Write(payload)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			return

		case "joinGame":
			//deletes previous game to prevent memory leak
			if gameID, ok := s.Get("game"); ok {
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
				s.Set("game", gameID)
			}

		case "leaveGame":
			var gameID string
			err = json.Unmarshal(unmarshal.Payload, &gameID)
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

			username, exists := s.Get("username")
			if !exists {
				log.Println("Error: Username does not exist")
				return
			}

			if board.GetTurnColor() == 0 && game.Player1 != username {
				log.Println("Error: Wrong color")
				return
			} else if board.GetTurnColor() == 1 && game.Player2 != username {
				log.Println("Error: Wrong color")
				return
			}

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
			gameID, exists := s.Get("game")
			if !exists {
				log.Println("Error: gameID does not exist")
				return
			}
			gameUUID, err := uuid.Parse(gameID.(string))
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			username, exists := s.Get("username")
			if !exists {
				log.Println("Error: Username does not exist")
				return
			}

			receiverID, err := dbConn.RequestGameUndo(gameUUID, username.(string))
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			for key := range GameConnections.connections[gameUUID] {
				if id, exists := key.Get("id"); exists && id.(string) == receiverID.String() {
					msgResponse := MsgResponse{
						Type:    "undoRequest",
						Payload: "",
					}
					marshalledResponse, err := json.Marshal(msgResponse)
					if err != nil {
						log.Println("Error: ", err)
						return
					}
					err = key.Write(marshalledResponse)
					if err != nil {
						log.Println("Error: ", err)
					}
				}
			}

		case "responseGameUndo":
			ctx := context.Background()
			gameID, exists := s.Get("game")
			if !exists {
				log.Println("Error: gameID does not exist")
				return
			}
			gameUUID, err := uuid.Parse(gameID.(string))
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			userID, exists := s.Get("id")
			if !exists {
				log.Println("Error: userID does not exist")
				return
			}
			userUUID, err := uuid.Parse(userID.(string))
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			var undoResponse string
			err = json.Unmarshal(unmarshal.Payload, &undoResponse)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if undoResponse != "accept" && undoResponse != "reject" {
				return
			}

			if undoResponse == "accept" {
				//TODO undo logic here
				game, err := getGame(dbConn, gameID.(string))
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
				board.UndoMove()
				board.PrintBoard()
				game.CurrentState = board.BoardToFen()
				game.History.String = board.SerializeHistory()
				game.History.Valid = true

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
				//TODO, update database
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
			}

			sender_id, err := dbConn.ResponseGameUndo(undoResponse, userUUID, gameUUID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			for key := range GameConnections.connections[gameUUID] {
				if id, exists := key.Get("id"); exists && id.(string) == sender_id.String() {
					msgResponse := MsgResponse{
						Type:    "undoResponse",
						Payload: undoResponse,
					}
					marshalledResponse, err := json.Marshal(msgResponse)
					if err != nil {
						log.Println("Error: ", err)
						return
					}
					err = key.Write(marshalledResponse)
					if err != nil {
						log.Println("Error: ", err)
					}
				}
			}

		case "completeGameUndo":
			gameID, exists := s.Get("game")
			if !exists {
				log.Println("Error: gameID does not exist")
				return
			}
			gameUUID, err := uuid.Parse(gameID.(string))
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			userID, exists := s.Get("id")
			if !exists {
				log.Println("Error: userID does not exist")
				return
			}
			userUUID, err := uuid.Parse(userID.(string))
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			err = dbConn.CompleteGameUndo(userUUID, gameUUID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

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

func getGame(dbs *api.DBConn, gameID string) (api.GameWithMoves, error) {
	game, err := dbs.GetGame(gameID)
	if err != nil {
		log.Println("Error: ", err)
		return api.GameWithMoves{}, err
	}
	return game, nil
}
