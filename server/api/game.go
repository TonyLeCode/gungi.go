package api

import (
	"context"
	"net/http"

	db "github.com/TonyLeCode/gungi.go/server/db/sqlc"
	"github.com/TonyLeCode/gungi.go/server/gungi"
	"github.com/google/uuid"

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

func (dbConn *DBConn) GetRoomList(c echo.Context) error {

	return nil
}

func (dbConn *DBConn) GetGame(c echo.Context) error {
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

	return c.JSON(http.StatusOK, game)
}
