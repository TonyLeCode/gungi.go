package api

import (
	"context"
	"log"
	"net/http"

	db "github.com/TonyLeCode/gungi.go/server/db/sqlc"
	"github.com/TonyLeCode/gungi.go/server/gungi"
	"github.com/google/uuid"
	"github.com/tabbed/pqtype"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type SerializedGame struct {
	Fen              string        `json:"fen"`
	History          string        `json:"history"`
	InCheck          bool          `json:"inCheck"`
	PawnDrop         []int         `json:"pawnDrop"`
	InvalidPawnFiles []int         `json:"invalidPawnFiles"`
	MoveList         map[int][]int `json:"moveList"`
}

func (dbConn *DBConn) GetOngoingGame(c echo.Context) error {
	board := gungi.Board{}
	serializedGame := SerializedGame{}

	//get fen from database
	board.SetBoardFromFen("")
	// get history from database
	// serializedGame.History =

	inCheck, inDoubleCheck, _, enemyXRaySquares, attackPiece, filteredMoves := board.GenerateLegalMoves()

	serializedGame.InCheck = inCheck
	serializedGame.MoveList = filteredMoves

	if inCheck && !inDoubleCheck {
		inbetween := []int{}
		if len(enemyXRaySquares.Path) != 0 {
			for _, move := range enemyXRaySquares.Path {
				if move.InBetween {
					inbetween = append(inbetween, move.Coordinate)
				}
			}
			inbetween = append(inbetween, enemyXRaySquares.Coordinate)
			serializedGame.PawnDrop = inbetween
		} else if attackPiece != -1 {
			serializedGame.PawnDrop[0] = attackPiece
		}
	}

	return nil
}

func (dbConn *DBConn) GetOngoingGameList(c echo.Context) error {
	ctx := context.Background()

	sub := c.Get("sub").(string)
	subid, err := uuid.Parse(sub)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	queries := db.New(dbConn.PostgresDB)

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

	queries := db.New(dbConn.PostgresDB)

	game, err := queries.GetGame(ctx, uuid)
	if err != nil {
		return GameWithMoves{}, err
	}

	newBoard := gungi.Board{}
	err = newBoard.SetBoardFromFen(game.CurrentState)
	if err != nil {
		return GameWithMoves{}, err
	}
	newBoard.SetMoveHistory(game.History.String)

	_, _, _, _, _, legalMoves := newBoard.GenerateLegalMoves()
	correctedLegalMoves := make(map[int][]int)
	for key, element := range legalMoves {
		correctedKey := gungi.SquareToIndex(key)
		for _, index := range element {
			correctedElement := gungi.SquareToIndex(index)
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

	queries := db.New(dbConn.PostgresDB)

	game, err := queries.GetGame(ctx, uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	newBoard := gungi.Board{}
	err = newBoard.SetBoardFromFen(game.CurrentState)
	if err != nil {
		return err
	}
	newBoard.SetMoveHistory(game.History.String)

	_, _, _, _, _, legalMoves := newBoard.GenerateLegalMoves()
	correctedLegalMoves := make(map[int][]int)
	for key, element := range legalMoves {
		correctedKey := gungi.SquareToIndex(key)
		for _, index := range element {
			correctedElement := gungi.SquareToIndex(index)
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

	queries := db.New(dbConn.PostgresDB)

	var nullRawUser1 pqtype.NullRawMessage
	var nullRawUser2 pqtype.NullRawMessage

	if username1 != "" {
		nullRawUser1.Valid = true
		nullRawUser1.RawMessage = []byte(username1)
	} else {
		nullRawUser1.Valid = false
	}

	if username2 != "" {
		nullRawUser2.Valid = true
		nullRawUser2.RawMessage = []byte(username2)
	} else {
		nullRawUser2.Valid = false
	}

	userID_1, err := queries.GetIdFromUsername(ctx, nullRawUser1)
	if err != nil {
		log.Println(err)
		return "", err
	}
	userID_2, err := queries.GetIdFromUsername(ctx, nullRawUser2)
	if err != nil {
		log.Println(err)
		return "", err
	}

	gameQuery := db.CreateGameParams{
		CurrentState: currentState,
		Ruleset:      ruleset,
		Type:         gameType,
	}

	gameID, err := queries.CreateGame(ctx, gameQuery)
	if err != nil {
		log.Println(err)
		return "", err
	}

	junctionQuery1 := db.GameJunctionParams{
		UserID: userID_1,
		GameID: gameID,
		Color:  "w",
	}

	junctionQuery2 := db.GameJunctionParams{
		UserID: userID_2,
		GameID: gameID,
		Color:  "b",
	}

	err = queries.GameJunction(ctx, junctionQuery1)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = queries.GameJunction(ctx, junctionQuery2)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return gameID.String(), nil
}
