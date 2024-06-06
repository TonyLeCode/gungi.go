package ws

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/olahol/melody"
	"github.com/whitemonarch/gungi-server/server/gungi"
	"github.com/whitemonarch/gungi-server/server/internal/api"
	"github.com/whitemonarch/gungi-server/server/internal/auth"
	db "github.com/whitemonarch/gungi-server/server/internal/db/sqlc"
)

type SessionHandlers struct {
	Listeners *Listeners
	Sessions  *Sessions
	ClientMsg ClientMsg
	Session   *melody.Session
	DBConn    *api.DBConn
}

type MsgHandlerFunc = func(s *melody.Session)

type MsgHandlersMap = map[string]MsgHandlerFunc

func (ss *SessionHandlers) handleAuth(s *melody.Session) {
	var token string
	err := json.Unmarshal(ss.ClientMsg.Payload, &token)
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
	ss.Sessions.AddUser(s, userUUID)

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
}

func (ss *SessionHandlers) handleJoinGame(s *melody.Session) {
	// log.Println("-----------------------------------------")
	log.Println("handleJoinGame")
	var gameID string
	err := json.Unmarshal(ss.ClientMsg.Payload, &gameID)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	if gameUUID, err := uuid.Parse(gameID); err == nil {
		unsub := ss.Listeners.AddListenerGame(s, gameUUID)
		ss.Sessions.ChangeGame(s, gameUUID)
		ss.Sessions.ChangeUnsub(s, unsub)
		log.Println("Joined Game")
		log.Println("Connections for Game", gameUUID, ": ")
		// for _, v := range ss.Listeners.listeners["game-"+gameUUID.String()] {
		// 	log.Println("-user: ", (*ss.Sessions)[v])
		// }
		// log.Println("-----------------------------------------")
	}
}
func (ss *SessionHandlers) handleLeaveGame(s *melody.Session) {
	log.Println("-----------------------------------------")
	log.Println("handleLeaveGame")
	var gameID string
	err := json.Unmarshal(ss.ClientMsg.Payload, &gameID)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	if gameUUID, err := uuid.Parse(gameID); err == nil {
		ss.Sessions.ChangeGame(s, uuid.Nil)
		ss.Listeners.RemoveListenerGame(s, gameUUID)
		ss.Sessions.Unsub(s)
		log.Println("Left Game: ", gameUUID)
		log.Println("Conns: ", ss.Listeners.listeners)
	}
	log.Println("-----------------------------------------")
}
func (ss *SessionHandlers) handleMakeGameMove(s *melody.Session) {
	ctx := context.Background()
	x := (*ss.Sessions)[s]
	log.Println(x)

	if (*ss.Sessions)[s].GameID == uuid.Nil || (*ss.Sessions)[s].ID == uuid.Nil {
		return
	}

	var makeGameMoveMsg makeGameMoveMsg
	err := json.Unmarshal(ss.ClientMsg.Payload, &makeGameMoveMsg)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	game, err := getGame(ss.DBConn, (*ss.Sessions)[s].GameID)
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

	if board.GetTurnColor() == 0 && game.User1 != (*ss.Sessions)[s].ID {
		log.Println("Error: Wrong color")
		return
	} else if board.GetTurnColor() == 1 && game.User2 != (*ss.Sessions)[s].ID {
		log.Println("Error: Wrong color")
		return
	}

	err = board.MakeMove(makeGameMoveMsg.FromPiece, board.ConvertInputCoord(makeGameMoveMsg.FromCoord), makeGameMoveMsg.MoveType, board.ConvertInputCoord(makeGameMoveMsg.ToCoord))
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	// board.PrintBoard()
	log.Println("Made move, new fen: ", board.BoardToFen())

	game.CurrentState = board.BoardToFen()
	game.History = board.SerializeHistory()

	makeMoveParams := db.MakeMoveParams{
		ID:           game.ID,
		CurrentState: game.CurrentState,
		History:      game.History,
	}
	queries := db.New(ss.DBConn.Conn)
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
	//TODO not good order because db will get updated by MakeMove before checking game result
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
			result.String = "draw"
			gameEndString = "Draw!"
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

	ss.Listeners.EmitGameMsg(b, game.ID)

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
		ss.Listeners.EmitGameMsg(b, game.ID)
	}
}
func (ss *SessionHandlers) handleGameResign(s *melody.Session) {
	ctx := context.Background()
	gameID := (*ss.Sessions)[s].GameID
	if gameID == uuid.Nil || (*ss.Sessions)[s].ID == uuid.Nil {
		return
	}
	userID := (*ss.Sessions)[s].ID

	resignParams := db.ResignGameParams{
		ID:    gameID,
		User1: userID,
	}

	queries := db.New(ss.DBConn.Conn)
	game, err := queries.ResignGame(ctx, resignParams)
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

	ss.Listeners.EmitGameMsg(b, game.ID)

	m = ServerMsg{
		Type:    "gameResign",
		Payload: game.Result,
	}

	b, err = json.Marshal(m)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	ss.Listeners.EmitGameMsgFilter(b, game.ID, func(s2 *melody.Session) bool {
		return s2 != s
	})
}
func (ss *SessionHandlers) handleRequestGameUndo(s *melody.Session) {
	gameID := (*ss.Sessions)[s].GameID
	if gameID == uuid.Nil || (*ss.Sessions)[s].ID == uuid.Nil {
		return
	}

	receiverID, err := ss.DBConn.RequestGameUndo(gameID, (*ss.Sessions)[s].ID)
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
	ss.Listeners.EmitGameMsgFilter(b, gameID, func(s *melody.Session) bool {
		return (*ss.Sessions)[s].ID == receiverID
	})
}
func (ss *SessionHandlers) handleResponseGameUndo(s *melody.Session) {
	ctx := context.Background()
	if (*ss.Sessions)[s].GameID == uuid.Nil || (*ss.Sessions)[s].ID == uuid.Nil {
		return
	}

	var undoResponse string
	err := json.Unmarshal(ss.ClientMsg.Payload, &undoResponse)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	if undoResponse != "accept" && undoResponse != "reject" {
		return
	}

	if undoResponse == "accept" {
		//TODO undo logic here
		game, err := getGame(ss.DBConn, (*ss.Sessions)[s].GameID)
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
		queries := db.New(ss.DBConn.Conn)
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

		ss.Listeners.EmitGameMsg(b, (*ss.Sessions)[s].GameID)
	}

	sender_id, err := ss.DBConn.ResponseGameUndo(undoResponse, (*ss.Sessions)[s].ID, (*ss.Sessions)[s].GameID)
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

	ss.Listeners.EmitGameMsgFilter(b, (*ss.Sessions)[s].GameID, func(s *melody.Session) bool {
		return (*ss.Sessions)[s].ID == sender_id
	})
}
func (ss *SessionHandlers) handleCompleteGameUndo(s *melody.Session) {
	if (*ss.Sessions)[s].GameID == uuid.Nil || (*ss.Sessions)[s].ID == uuid.Nil {
		return
	}

	err := ss.DBConn.CompleteGameUndo((*ss.Sessions)[s].ID, (*ss.Sessions)[s].GameID)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
}
func (ss *SessionHandlers) handleJoinPlay(s *melody.Session) {
	unsub := ss.Listeners.AddListenerRooms(s)
	ss.Sessions.ChangeUnsub(s, unsub)
	log.Println("Joined Play")
	log.Println("Conns: ", ss.Listeners.listeners)

	roomList, err := ss.DBConn.GetRoomList()
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
}
func (ss *SessionHandlers) handleLeavePlay(s *melody.Session) {
	ss.Listeners.RemoveListenerRooms(s)
	ss.Sessions.Unsub(s)
	log.Println("Left Play")
	log.Println("Conns: ", ss.Listeners.listeners)
}
func (ss *SessionHandlers) handleCreatePlayRoom(s *melody.Session) {
	if (*ss.Sessions)[s].ID == uuid.Nil {
		log.Println("Error: no user")
		return
	}
	var roomPayload CreateGameRoomRequest
	err := json.Unmarshal(ss.ClientMsg.Payload, &roomPayload)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	err = ss.DBConn.CreateRoom((*ss.Sessions)[s].ID, roomPayload.Description, roomPayload.Rules, roomPayload.Type, roomPayload.Color)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	roomList, err := ss.DBConn.GetRoomList()
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

	ss.Listeners.EmitRoomMsg(b)
}
func (ss *SessionHandlers) handleAcceptPlayRoom(s *melody.Session) {
	if (*ss.Sessions)[s].ID == uuid.Nil {
		return
	}
	var roomid string
	err := json.Unmarshal(ss.ClientMsg.Payload, &roomid)
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
	room, err := ss.DBConn.DeleteRoom(roomUUID)
	if err != nil {
		log.Println(err)
		return
	}

	gameID, err := ss.DBConn.CreateGameFromRoom(room, (*ss.Sessions)[s].ID)
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
	ss.Listeners.EmitRoomMsgFilter(b, func(s *melody.Session) bool {
		return (*ss.Sessions)[s].ID == room.HostID
	})

	roomList, err := ss.DBConn.GetRoomList()
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

	ss.Listeners.EmitRoomMsg(b)
}
func (ss *SessionHandlers) handleCancelPlayRoom(s *melody.Session) {
	if (*ss.Sessions)[s].ID == uuid.Nil {
		return
	}
	var roomid string
	err := json.Unmarshal(ss.ClientMsg.Payload, &roomid)
	if err != nil {
		log.Println(err)
		return
	}

	roomUUID, err := uuid.Parse(roomid)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = ss.DBConn.DeleteRoomSafe(roomUUID, (*ss.Sessions)[s].ID)
	if err != nil {
		log.Println(err)
		return
	}

	roomList, err := ss.DBConn.GetRoomList()
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

	ss.Listeners.EmitRoomMsg(b)
}
