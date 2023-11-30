package api

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/whitemonarch/gungi-server/server/gungi"
	db "github.com/whitemonarch/gungi-server/server/internal/db/sqlc"

	"github.com/labstack/echo/v4"
)

type SerializedGame struct {
	Fen         string        `json:"fen"`
	History     string        `json:"history"`
	CheckStatus string        `json:"checkStatus"`
	Ruleset     string        `json:"ruleset"`
	MoveList    map[int][]int `json:"moveList"`
}

func (dbConn *DBConn) GetOngoingGame(c echo.Context) error {
	board := gungi.CreateBoard("revised")
	serializedGame := SerializedGame{}

	//get fen from database
	board.SetBoardFromFen("")
	// get history from database
	// serializedGame.History =

	checkStatus, moves := board.GetLegalMoves()

	serializedGame.CheckStatus = checkStatus
	serializedGame.MoveList = moves

	// if inCheck && !inDoubleCheck {
	// 	inbetween := []int{}
	// 	if len(enemyXRaySquares.Path) != 0 {
	// 		for _, move := range enemyXRaySquares.Path {
	// 			if move.InBetween {
	// 				inbetween = append(inbetween, move.Coordinate)
	// 			}
	// 		}
	// 		inbetween = append(inbetween, enemyXRaySquares.Coordinate)
	// 		serializedGame.PawnDrop = inbetween
	// 	} else if attackPiece != -1 {
	// 		serializedGame.PawnDrop[0] = attackPiece
	// 	}
	// }

	return nil
}

func (dbConn *DBConn) GetOngoingGameList(c echo.Context) error {
	ctx := context.Background()

	sub := c.Get("sub").(string)
	subid, err := uuid.Parse(sub)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	queries := db.New(dbConn.Conn)

	games, err := queries.GetOngoingGames(ctx, subid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, games)
}

type GameWithMoves struct {
	ID           uuid.UUID          `json:"id"`
	Fen          pgtype.Text        `json:"fen"`
	History      pgtype.Text        `json:"history"`
	Completed    bool               `json:"completed"`
	DateStarted  pgtype.Timestamptz `json:"date_started"`
	DateFinished pgtype.Timestamptz `json:"date_finished"`
	CurrentState string             `json:"current_state"`
	Ruleset      string             `json:"ruleset"`
	Type         string             `json:"type"`
	Player1      string             `json:"player1"`
	Player2      string             `json:"player2"`
	MoveList     map[int][]int      `json:"moveList"`
}

type GameWithUndo struct {
	ID           uuid.UUID          `json:"id"`
	Fen          pgtype.Text        `json:"fen"`
	History      pgtype.Text        `json:"history"`
	Completed    bool               `json:"completed"`
	DateStarted  pgtype.Timestamptz `json:"date_started"`
	DateFinished pgtype.Timestamptz `json:"date_finished"`
	CurrentState string             `json:"current_state"`
	Ruleset      string             `json:"ruleset"`
	Type         string             `json:"type"`
	User1        uuid.UUID          `json:"user1"`
	User2        uuid.UUID          `json:"user2"`
	Player1      string             `json:"player1"`
	Player2      string             `json:"player2"`
	MoveList     map[int][]int      `json:"moveList"`
}

func (dbConn *DBConn) GetGame(id uuid.UUID) (GameWithUndo, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	game, err := queries.GetGame(ctx, id)
	if err != nil {
		return GameWithUndo{}, err
	}

	//TODO properly get ruleset
	newBoard := gungi.CreateBoard("revised")
	err = newBoard.SetBoardFromFen(game.CurrentState)
	if err != nil {
		return GameWithUndo{}, err
	}
	newBoard.SetHistory(strings.Split(game.History.String, " "))

	_, legalMoves := newBoard.GetLegalMoves()
	correctedLegalMoves := make(map[int][]int)
	for key, element := range legalMoves {
		correctedKey := newBoard.ConvertOutputCoord(key)
		for _, index := range element {
			correctedElement := newBoard.ConvertOutputCoord(index)
			correctedLegalMoves[correctedKey] = append(correctedLegalMoves[correctedKey], correctedElement)
		}
	}
	gameWithUndo := GameWithUndo{
		ID:           game.ID,
		Fen:          game.Fen,
		History:      game.History,
		Completed:    game.Completed,
		DateStarted:  game.DateStarted,
		DateFinished: game.DateFinished,
		CurrentState: game.CurrentState,
		Ruleset:      game.Ruleset,
		Type:         game.Type,
		User1:        game.User1,
		User2:        game.User2,
		Player1:      game.Player1,
		Player2:      game.Player2,
		MoveList:     correctedLegalMoves,
	}

	return gameWithUndo, nil
}

func (dbConn *DBConn) GetGameRoute(c echo.Context) error {
	ctx := context.Background()

	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	queries := db.New(dbConn.Conn)

	game, err := queries.GetGame(ctx, uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	newBoard := gungi.CreateBoard("revised")
	err = newBoard.SetBoardFromFen(game.CurrentState)
	if err != nil {
		return err
	}
	newBoard.SetHistory(strings.Split(game.History.String, " "))

	_, legalMoves := newBoard.GetLegalMoves()
	correctedLegalMoves := make(map[int][]int)
	for key, element := range legalMoves {
		correctedKey := newBoard.ConvertOutputCoord(key)
		for _, index := range element {
			correctedElement := newBoard.ConvertOutputCoord(index)
			correctedLegalMoves[correctedKey] = append(correctedLegalMoves[correctedKey], correctedElement)
		}
	}
	log.Println(correctedLegalMoves)
	log.Println("fen: ", newBoard.BoardToFen())
	gameWithMoves := GameWithMoves{
		ID:           game.ID,
		Fen:          game.Fen,
		History:      game.History,
		Completed:    game.Completed,
		DateStarted:  game.DateStarted,
		DateFinished: game.DateFinished,
		CurrentState: game.CurrentState,
		Ruleset:      game.Ruleset,
		Type:         game.Type,
		Player1:      game.Player1,
		Player2:      game.Player2,
		MoveList:     correctedLegalMoves,
	}

	return c.JSON(http.StatusOK, gameWithMoves)
}

type GameWithUndoRowAndMoves struct {
	ID           uuid.UUID          `json:"id"`
	Fen          pgtype.Text        `json:"fen"`
	History      pgtype.Text        `json:"history"`
	Completed    bool               `json:"completed"`
	DateStarted  pgtype.Timestamptz `json:"date_started"`
	DateFinished pgtype.Timestamptz `json:"date_finished"`
	CurrentState string             `json:"current_state"`
	Ruleset      string             `json:"ruleset"`
	Type         string             `json:"type"`
	Player1      string             `json:"player1"`
	Player2      string             `json:"player2"`
	MoveList     map[int][]int      `json:"moveList"`
	UndoRequests []UndoRequests     `json:"undo_requests"`
}

type UndoRequests struct {
	SenderUsername   string `json:"sender_username"`
	ReceiverUsername string `json:"receiver_username"`
	Status           string `json:"status"`
}

func (dbConn *DBConn) GetGameWithUndoRoute(c echo.Context) error {
	ctx := context.Background()

	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	queries := db.New(dbConn.Conn)

	game, err := queries.GetGameWithUndo(ctx, uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	newBoard := gungi.CreateBoard("revised")
	err = newBoard.SetBoardFromFen(game.CurrentState)
	if err != nil {
		return err
	}
	newBoard.SetHistory(strings.Split(game.History.String, " "))

	_, legalMoves := newBoard.GetLegalMoves()
	correctedLegalMoves := make(map[int][]int)
	for key, element := range legalMoves {
		correctedKey := newBoard.ConvertOutputCoord(key)
		for _, index := range element {
			correctedElement := newBoard.ConvertOutputCoord(index)
			correctedLegalMoves[correctedKey] = append(correctedLegalMoves[correctedKey], correctedElement)
		}
	}
	log.Println(correctedLegalMoves)
	log.Println("fen: ", newBoard.BoardToFen())

	log.Println(string(game.UndoRequests))

	var undoRequests []UndoRequests = []UndoRequests{}
	// if game.UndoRequests == nil {
	// 	undoRequests = []UndoRequests{}
	// } else {

	// }
	err = json.Unmarshal(game.UndoRequests, &undoRequests)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	gameWithMoves := GameWithUndoRowAndMoves{
		ID:           game.ID,
		Fen:          game.Fen,
		History:      game.History,
		Completed:    game.Completed,
		DateStarted:  game.DateStarted,
		DateFinished: game.DateFinished,
		CurrentState: game.CurrentState,
		Ruleset:      game.Ruleset,
		Type:         game.Type,
		Player1:      game.Player1,
		Player2:      game.Player2,
		MoveList:     correctedLegalMoves,
		UndoRequests: undoRequests,
	}

	return c.JSON(http.StatusOK, gameWithMoves)
}

func (dbConn *DBConn) CreateGame(currentState string, ruleset string, gameType string, user1 uuid.UUID, user2 uuid.UUID) (string, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	gameQuery := db.CreateGameParams{
		CurrentState: currentState,
		Ruleset:      ruleset,
		Type:         gameType,
		User1:        user1,
		User2:        user2,
	}

	gameID, err := queries.CreateGame(ctx, gameQuery)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return gameID.String(), nil
}

func (dbConn *DBConn) CreateGameFromRoom(room db.DeleteRoomRow, acceptingUser uuid.UUID) (string, error) {
	ctx := context.Background()
	var user1 uuid.UUID
	var user2 uuid.UUID
	var currentState string
	switch room.Color {
	case "white":
		user1 = room.HostID
		user2 = acceptingUser
	case "black":
		user1 = acceptingUser
		user2 = room.HostID
	case "random":
		r := rand.Intn(2)
		if r == 0 {
			user1 = room.HostID
			user2 = acceptingUser
		} else {
			user1 = acceptingUser
			user2 = room.HostID
		}
	default:
		r := rand.Intn(2)
		if r == 0 {
			user1 = room.HostID
			user2 = acceptingUser
		} else {
			user1 = acceptingUser
			user2 = room.HostID
		}
	}
	switch room.Rules {
	case "revised":
		currentState = "9/9/9/9/9/9/9/9/9 9446222122211/9446222122211 w 00"
	case "universal-music":
		currentState = "9/9/9/9/9/9/9/9/9 9446222122211/9446222122211 w 00"
	case "default":
		currentState = "9/9/9/9/9/9/9/9/9 9446222122211/9446222122211 w 00"
	default:
		currentState = "9/9/9/9/9/9/9/9/9 9446222122211/9446222122211 w 00"
	}
	queriesParams := db.CreateGameParams{
		CurrentState: currentState,
		Ruleset:      room.Rules,
		Type:         room.Type,
		User1:        user1,
		User2:        user2,
	}
	queries := db.New(dbConn.Conn)
	gameID, err := queries.CreateGame(ctx, queriesParams)
	if err != nil {
		return "", err
	}
	return gameID.String(), nil
}

func (dbConn *DBConn) CreateRoom(hostID uuid.UUID, description string, rules string, roomType string, color string) error {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	err := queries.CreateRoom(ctx, db.CreateRoomParams{
		HostID:      hostID,
		Description: description,
		Rules:       rules,
		Type:        roomType,
		Color:       color,
	})
	if err != nil {
		return err
	}

	return nil
}

func (dbConn *DBConn) GetRoomList() ([]db.GetRoomListRow, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	roomList, err := queries.GetRoomList(ctx)
	if err != nil {
		return nil, err
	}

	return roomList, nil
}

func (dbConn *DBConn) DeleteRoomSafe(roomID uuid.UUID, hostID uuid.UUID) (db.DeleteRoomSafeRow, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)
	params := db.DeleteRoomSafeParams{
		ID:     roomID,
		HostID: hostID,
	}

	room, err := queries.DeleteRoomSafe(ctx, params)
	if err != nil {
		return db.DeleteRoomSafeRow{}, err
	}

	return room, nil
}

func (dbConn *DBConn) DeleteRoom(roomID uuid.UUID) (db.DeleteRoomRow, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	room, err := queries.DeleteRoom(ctx, roomID)
	if err != nil {
		return db.DeleteRoomRow{}, err
	}

	return room, nil
}

func (dbConn *DBConn) RequestGameUndo(gameID uuid.UUID, senderID uuid.UUID) (uuid.UUID, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	createUndoParams := db.CreateUndoParams{
		GameID:   gameID,
		SenderID: senderID,
	}

	receiver_id, err := queries.CreateUndo(ctx, createUndoParams)
	if err != nil {
		return uuid.UUID{}, err
	}

	return receiver_id, nil
}

func (dbConn *DBConn) ResponseGameUndo(status string, receiverID uuid.UUID, gameID uuid.UUID) (uuid.UUID, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	ChangeUndoParams := db.ChangeUndoParams{
		Status:     status,
		ReceiverID: receiverID,
		GameID:     gameID,
	}

	sender_id, err := queries.ChangeUndo(ctx, ChangeUndoParams)
	if err != nil {
		return uuid.UUID{}, err
	}

	return sender_id, nil
}

func (dbConn *DBConn) CompleteGameUndo(senderID uuid.UUID, gameID uuid.UUID) error {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	ChangeUndoParams := db.RemoveUndoParams{
		SenderID: senderID,
		GameID:   gameID,
	}

	err := queries.RemoveUndo(ctx, ChangeUndoParams)
	if err != nil {
		return err
	}

	return nil
}
