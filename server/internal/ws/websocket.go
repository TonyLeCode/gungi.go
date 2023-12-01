package ws

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	FromPiece int `json:"fromPiece"`
	FromCoord int `json:"fromCoord"`
	MoveType  int `json:"moveType"`
	ToCoord   int `json:"toCoord"`
}

func WS(m *melody.Melody, dbConn *api.DBConn) echo.HandlerFunc {
	listeners := Listeners{}
	sessions := Sessions{}

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

		if clientMsg.Type != "auth" && sessions[s] == nil {
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

			userid := claims["sub"].(string)
			userUUID, err := uuid.Parse(userid)
			if err != nil {
				log.Println(err)
				return
			}
			sessions.AddUser(s, userUUID)

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
			var gameID string
			err = json.Unmarshal(clientMsg.Payload, &gameID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if gameUUID, err := uuid.Parse(gameID); err == nil {
				unsub := listeners.AddListenerGame(s, gameUUID)
				sessions.ChangeGame(s, gameUUID)
				sessions.ChangeUnsub(s, unsub)
				log.Println("Joined Game: ", gameUUID)
				log.Println("Conns: ", listeners.listeners)
			}

		case "leaveGame":
			var gameID string
			err = json.Unmarshal(clientMsg.Payload, &gameID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if gameUUID, err := uuid.Parse(gameID); err == nil {
				sessions.ChangeGame(s, uuid.Nil)
				listeners.RemoveListenerGame(s, gameUUID)
				log.Println("Left Game: ", gameUUID)
				log.Println("Conns: ", listeners.listeners)
			}

		case "makeGameMove":
			ctx := context.Background()

			if sessions[s].GameID == uuid.Nil {
				return
			}

			var makeGameMoveMsg makeGameMoveMsg
			err = json.Unmarshal(clientMsg.Payload, &makeGameMoveMsg)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			game, err := getGame(dbConn, sessions[s].GameID)
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

			board.SetHistory(strings.Fields(game.History))
			// board.PrintBoard()

			if board.GetTurnColor() == 0 && game.User1 != sessions[s].ID {
				log.Println("Error: Wrong color")
				return
			} else if board.GetTurnColor() == 1 && game.User2 != sessions[s].ID {
				log.Println("Error: Wrong color")
				return
			}

			err = board.MakeMove(makeGameMoveMsg.FromPiece, board.ConvertInputCoord(makeGameMoveMsg.FromCoord), makeGameMoveMsg.MoveType, board.ConvertInputCoord(makeGameMoveMsg.ToCoord))
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			// board.PrintBoard()
			log.Println("fen: ", board.BoardToFen())

			game.CurrentState = board.BoardToFen()
			game.History = board.SerializeHistory()

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

			checkStatus, legalMoves := board.GetLegalMoves()
			correctedLegalMoves := make(map[int][]int)
			for key, element := range legalMoves {
				correctedKey := board.ConvertOutputCoord(key)
				for _, index := range element {
					correctedElement := board.ConvertOutputCoord(index)
					correctedLegalMoves[correctedKey] = append(correctedLegalMoves[correctedKey], correctedElement)
				}
			}
			//TODO not good order
			game.MoveList = correctedLegalMoves

			var gameEndString string

			if checkStatus == "checkmated" || checkStatus == "stalemate" {
				result := pgtype.Text{
					String: "",
					Valid:  true,
				}

				if board.GetTurnColor() == 0 && checkStatus == "checkmated" {
					result.String = "b"
					gameEndString = "Black Wins by checkmate!"
				} else if board.GetTurnColor() == 1 && checkStatus == "checkmated" {
					result.String = "w"
					gameEndString = "White Wins by checkmate!"
				} else {
					result.String = "stalemate"
					gameEndString = "Stalemate!"
				}

				changeGameParams := db.ChangeGameResultParams{
					ID:        game.ID,
					Completed: true,
					Result:    result,
				}

				err = queries.ChangeGameResult(ctx, changeGameParams)
				if err != nil {
					log.Println("Error: ", err)
					return
				}

				game.Completed = true
				game.Result = result
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
			// log.Println(correctedLegalMoves)

			listeners.EmitGameMsg(b, game.ID)

			m = ServerMsg{
				Type:    "gameEnd",
				Payload: gameEndString,
			}
			b, err = json.Marshal(m)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			if gameEndString != "" {
				listeners.EmitGameMsg(b, game.ID)
			}

		//TODO finish cases
		case "gameResign":

		case "requestGameUndo":
			gameID := sessions[s].GameID
			if gameID == uuid.Nil {
				return
			}

			receiverID, err := dbConn.RequestGameUndo(gameID, sessions[s].ID)
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
			listeners.EmitGameMsgFilter(b, gameID, func(s *melody.Session) bool {
				return sessions[s].ID == receiverID
			})

		case "responseGameUndo":
			ctx := context.Background()
			if sessions[s].GameID == uuid.Nil {
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
				game, err := getGame(dbConn, sessions[s].GameID)
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

				board.SetHistory(strings.Fields(game.History))
				// board.PrintBoard()
				board.UndoMove()
				// board.PrintBoard()
				game.CurrentState = board.BoardToFen()
				game.History = board.SerializeHistory()

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

				listeners.EmitGameMsg(b, sessions[s].GameID)
			}

			sender_id, err := dbConn.ResponseGameUndo(undoResponse, sessions[s].ID, sessions[s].GameID)
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

			listeners.EmitGameMsgFilter(b, sessions[s].GameID, func(s *melody.Session) bool {
				return sessions[s].ID == sender_id
			})

		case "completeGameUndo":
			if sessions[s].GameID == uuid.Nil {
				return
			}

			err = dbConn.CompleteGameUndo(sessions[s].ID, sessions[s].GameID)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

		case "joinPlay":
			unsub := listeners.AddListenerRooms(s)
			sessions.ChangeUnsub(s, unsub)
			log.Println("Joined Play")
			log.Println("Conns: ", listeners.listeners)

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
			s.Write(b)

		case "leavePlay":
			listeners.RemoveListenerRooms(s)
			log.Println("Left Play")
			log.Println("Conns: ", listeners.listeners)

		case "createPlayRoom":
			var roomPayload CreateGameRoomRequest
			err := json.Unmarshal(clientMsg.Payload, &roomPayload)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			err = dbConn.CreateRoom(sessions[s].ID, roomPayload.Description, roomPayload.Rules, roomPayload.Type, roomPayload.Color)
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

		case "acceptPlayRoom":

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

			gameID, err := dbConn.CreateGameFromRoom(room, sessions[s].ID)
			if err != nil {
				log.Println(err)
				return
			}
			m := ServerMsg{
				Type:    "roomAccepted",
				Payload: gameID,
			}
			b, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
				return
			}

			s.Write(b)
			listeners.EmitRoomMsgFilter(b, func(s *melody.Session) bool {
				return sessions[s].ID == room.HostID
			})

			roomList, err := dbConn.GetRoomList()
			if err != nil {
				log.Println(err)
				return
			}

			m = ServerMsg{
				Type:    "roomList",
				Payload: roomList,
			}
			b, err = json.Marshal(m)
			if err != nil {
				log.Println(err)
				return
			}

			listeners.EmitRoomMsg(b)

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

			_, err = dbConn.DeleteRoomSafe(roomUUID, sessions[s].ID)
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

func getGame(dbs *api.DBConn, gameID uuid.UUID) (api.GameWithUndo, error) {
	game, err := dbs.GetGame(gameID)
	if err != nil {
		log.Println("Error: ", err)
		return api.GameWithUndo{}, err
	}
	return game, nil
}
