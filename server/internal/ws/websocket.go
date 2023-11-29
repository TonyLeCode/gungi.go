package ws

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"

	"github.com/whitemonarch/gungi-server/server/gungi"
	"github.com/whitemonarch/gungi-server/server/internal/api"
	"github.com/whitemonarch/gungi-server/server/internal/auth"
	db "github.com/whitemonarch/gungi-server/server/internal/db/sqlc"
)

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

type ClientMsg struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ServerMsg struct {
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

func WS(m *melody.Melody, dbConn *api.DBConn) echo.HandlerFunc {
	// PlayConnections := PlayConnections{
	// 	connections: make(map[*melody.Session]bool),
	// 	mu:          sync.RWMutex{},
	// }
	// GameConnections := GameConnections{
	// 	connections: make(map[uuid.UUID]map[*melody.Session]bool),
	// 	mu:          sync.RWMutex{},
	// }
	listeners := Listeners{}
	sessions := Sessions{}

	// log.Println(msgHandlers)

	m.HandleConnect(func(s *melody.Session) {
		log.Println("Connected", s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("Disconnect", s, sessions[s].ID)
		sessions.RemoveUser(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var clientMsg ClientMsg
		err := json.Unmarshal(msg, &clientMsg)
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		switch clientMsg.Type {
		case "auth":
			var token string
			err := json.Unmarshal(clientMsg.Payload, &token)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			claims, err := auth.AuthenticateSupabaseToken(token)
			if err != nil {
				log.Println(err)

				m := ServerMsg{
					Type:    "auth",
					Payload: "failure",
				}

				b, err := json.Marshal(m)
				if err != nil {
					log.Println("Error: ", err)
					return
				}

				err = s.Write(b)
				if err != nil {
					log.Println("Error: ", err)
				}
				return
			}

			// metadata := claims["user_metadata"].(map[string]interface{})
			// username := metadata["username"].(string)
			userid := claims["sub"].(string)
			userUUID, err := uuid.Parse(userid)
			if err != nil {
				log.Println(err)
				return
			}
			sessions.AddUser(s, userUUID)
			// s.Set("username", username)
			// s.Set("id", id)

			m := ServerMsg{
				Type:    "auth",
				Payload: "success",
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			err = s.Write(b)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
		case "joinGame":
			//deletes previous game to prevent memory leak
			// if gameID, ok := s.Get("game"); ok {
			// 	if gameUUID, err := uuid.Parse(gameID.(string)); err == nil {
			// 		log.Println("Removed Game Connection")
			// 		GameConnections.RemoveConnection(gameUUID, s)
			// 	}
			// }

			var gameID string
			err = json.Unmarshal(clientMsg.Payload, &gameID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if gameUUID, err := uuid.Parse(gameID); err == nil {
				unsub := listeners.AddListenerGame(s, gameUUID)
				sessions.ChangeUnsub(s, unsub)
				log.Println("Joined Game: ", gameUUID)
				log.Println("Conns: ", listeners.listeners)
				// GameConnections.AddConnection(gameUUID, s)
				// log.Println("Joined Game: ", gameUUID)
				// log.Println("Conns: ", GameConnections.connections)
				// s.Set("game", gameID)
			}

		case "leaveGame":
			var gameID string
			err = json.Unmarshal(clientMsg.Payload, &gameID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if gameUUID, err := uuid.Parse(gameID); err == nil {
				listeners.RemoveListenerGame(s, gameUUID)
				log.Println("Left Game: ", gameUUID)
				log.Println("Conns: ", listeners.listeners)
				// GameConnections.RemoveConnection(gameUUID, s)
				// log.Println("Conns: ", GameConnections.connections)
			}

		case "makeGameMove":
			ctx := context.Background()
			var makeGameMoveMsg makeGameMoveMsg
			err = json.Unmarshal(clientMsg.Payload, &makeGameMoveMsg)
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

			m := ServerMsg{
				Type:    "game",
				Payload: game,
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			listeners.EmitGameMsg(b, game.ID)
			// var keys []*melody.Session
			// for key := range GameConnections.connections[game.ID] {
			// 	keys = append(keys, key)
			// }

			// m.BroadcastMultiple(marshalledResponse, keys)

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

			m := ServerMsg{
				Type:    "undoRequest",
				Payload: "",
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			listeners.EmitGameMsgFilter(b, gameUUID, func(s *melody.Session) bool {
				return sessions[s].ID == receiverID
			})

			// for key := range GameConnections.connections[gameUUID] {
			// 	if id, exists := key.Get("id"); exists && id.(string) == receiverID.String() {
			// 		ServerMsg := ServerMsg{
			// 			Type:    "undoRequest",
			// 			Payload: "",
			// 		}
			// 		marshalledResponse, err := json.Marshal(ServerMsg)
			// 		if err != nil {
			// 			log.Println("Error: ", err)
			// 			return
			// 		}
			// 		err = key.Write(marshalledResponse)
			// 		if err != nil {
			// 			log.Println("Error: ", err)
			// 		}
			// 	}
			// }

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
			err = json.Unmarshal(clientMsg.Payload, &undoResponse)
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

				m := ServerMsg{
					Type:    "game",
					Payload: game,
				}

				b, err := json.Marshal(m)
				if err != nil {
					log.Println("Error: ", err)
					return
				}

				listeners.EmitGameMsg(b, gameUUID)

				// var keys []*melody.Session
				// for key := range GameConnections.connections[game.ID] {
				// 	keys = append(keys, key)
				// }

				// m.BroadcastMultiple(marshalledResponse, keys)
			}

			sender_id, err := dbConn.ResponseGameUndo(undoResponse, userUUID, gameUUID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			m := ServerMsg{
				Type:    "undoResponse",
				Payload: undoResponse,
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			listeners.EmitGameMsgFilter(b, gameUUID, func(s *melody.Session) bool {
				return sessions[s].ID == sender_id
			})

			// for key := range GameConnections.connections[gameUUID] {
			// 	if id, exists := key.Get("id"); exists && id.(string) == sender_id.String() {
			// 		ServerMsg := ServerMsg{
			// 			Type:    "undoResponse",
			// 			Payload: undoResponse,
			// 		}
			// 		marshalledResponse, err := json.Marshal(ServerMsg)
			// 		if err != nil {
			// 			log.Println("Error: ", err)
			// 			return
			// 		}
			// 		err = key.Write(marshalledResponse)
			// 		if err != nil {
			// 			log.Println("Error: ", err)
			// 		}
			// 	}
			// }

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
			// PlayConnections.AddConnection(s)
			// log.Println("Joined Play")
			// log.Println("Conns: ", PlayConnections.connections)
			unsub := listeners.AddListenerRooms(s)
			sessions.ChangeUnsub(s, unsub)
			log.Println("Joined Play")
			log.Println("Conns: ", listeners.listeners)

			// TODO host id to username
			roomList, err := dbConn.GetRoomList()
			if err != nil {
				log.Println(err)
				return
			}
			roomListMsg := ServerMsg{
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
			listeners.RemoveListenerRooms(s)
			log.Println("Left Play")
			log.Println("Conns: ", listeners.listeners)
			// log.Println("Left Play")
			// PlayConnections.RemoveConnection(s)
			// log.Println("Conns: ", PlayConnections.connections)

		case "createPlayRoom":
			var roomPayload CreateGameRoomRequest
			err := json.Unmarshal(clientMsg.Payload, &roomPayload)
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
			m := ServerMsg{
				Type:    "roomList",
				Payload: roomList,
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
				return
			}

			listeners.EmitRoomMsg(b)

			//TODO keys to array
			// err = m.BroadcastFilter(marshaledRoomListMsg, func(q *melody.Session) bool {
			// 	return PlayConnections.connections[q]
			// })
			// if err != nil {
			// 	log.Println(err)
			// 	return
			// }

		case "acceptPlayRoom":
			username, exists := s.Get("username")
			if !exists {
				return
			}

			var roomid string
			err := json.Unmarshal(clientMsg.Payload, &roomid)
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
			gameAcceptedMsg := ServerMsg{
				Type:    "roomAccepted",
				Payload: gameID,
			}
			marshaledGameAcceptedMsg, err := json.Marshal(gameAcceptedMsg)
			if err != nil {
				log.Println(err)
				return
			}

			//TODO change broadcast
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

			m := ServerMsg{
				Type:    "roomList",
				Payload: roomList,
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
				return
			}

			listeners.EmitRoomMsg(b)
			// var keys []*melody.Session
			// for key := range PlayConnections.connections {
			// 	keys = append(keys, key)
			// }

			// m.BroadcastMultiple(marshaledUpdatedRoomListMsg, keys)

		case "cancelPlayRoom":
			var roomid string
			err := json.Unmarshal(clientMsg.Payload, &roomid)
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

			m := ServerMsg{
				Type:    "roomList",
				Payload: roomList,
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
				return
			}

			listeners.EmitRoomMsg(b)

			// var keys []*melody.Session
			// for key := range PlayConnections.connections {
			// 	keys = append(keys, key)
			// }

			// m.BroadcastMultiple(marshaledUpdatedRoomListMsg, keys)
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
