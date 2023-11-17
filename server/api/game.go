package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strings"

	db "github.com/TonyLeCode/gungi.go/server/db/sqlc"
	"github.com/TonyLeCode/gungi.go/server/gungi"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
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
	db.GetGameRow
	MoveList map[int][]int `json:"moveList"`
}

func (dbConn *DBConn) GetGame(id string) (GameWithMoves, error) {
	ctx := context.Background()

	uuid, err := uuid.Parse(id)
	if err != nil {
		return GameWithMoves{}, err
	}

	queries := db.New(dbConn.Conn)

	game, err := queries.GetGame(ctx, uuid)
	if err != nil {
		return GameWithMoves{}, err
	}

	//TODO properly get ruleset
	newBoard := gungi.CreateBoard("revised")
	err = newBoard.SetBoardFromFen(game.CurrentState)
	if err != nil {
		return GameWithMoves{}, err
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
	gameWithMoves := GameWithMoves{
		GetGameRow: game,
		MoveList:   correctedLegalMoves,
	}

	return gameWithMoves, nil
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
		GetGameRow: game,
		MoveList:   correctedLegalMoves,
	}

	return c.JSON(http.StatusOK, gameWithMoves)
}

func (dbConn *DBConn) CreateGame(currentState string, ruleset string, gameType string, username1 string, username2 string) (string, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	gameQuery := db.CreateGameParams{
		CurrentState: currentState,
		Ruleset:      ruleset,
		Type:         gameType,
		Username:     username1,
		Username_2:   username2,
	}

	gameID, err := queries.CreateGame(ctx, gameQuery)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return gameID.String(), nil
}

func (dbConn *DBConn) CreateGameFromRoom(room db.DeleteRoomRow, acceptingUsername string) (string, error) {
	ctx := context.Background()
	var username1 string
	var username2 string
	var currentState string
	switch room.Color {
	case "white":
		username1 = room.Host
		username2 = acceptingUsername
	case "black":
		username1 = acceptingUsername
		username2 = room.Host
	case "random":
		r := rand.Intn(2)
		if r == 0 {
			username1 = room.Host
			username2 = acceptingUsername
		} else {
			username1 = acceptingUsername
			username2 = room.Host
		}
	default:
		r := rand.Intn(2)
		if r == 0 {
			username1 = room.Host
			username2 = acceptingUsername
		} else {
			username1 = acceptingUsername
			username2 = room.Host
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
		Username:     username1,
		Username_2:   username2,
	}
	queries := db.New(dbConn.Conn)
	gameID, err := queries.CreateGame(ctx, queriesParams)
	if err != nil {
		return "", err
	}
	return gameID.String(), nil
}

func (dbConn *DBConn) CreateRoom(username string, description string, rules string, roomType string, color string) error {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	err := queries.CreateRoom(ctx, db.CreateRoomParams{
		Username:    username,
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

func (dbConn *DBConn) DeleteRoom(id uuid.UUID) (db.DeleteRoomRow, error) {
	ctx := context.Background()

	queries := db.New(dbConn.Conn)

	room, err := queries.DeleteRoom(ctx, id)
	if err != nil {
		return db.DeleteRoomRow{}, err
	}

	return room, nil
}
