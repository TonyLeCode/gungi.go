package api

import (

	// "github.com/whitemonarch/gungi-server/server/db"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	db "github.com/whitemonarch/gungi-server/server/internal/db/sqlc"
)

type UserDataRequest struct {
	ID string `json:"id"`
}

func (dbConn *DBConn) GetUsername(c echo.Context) error {
	ctx := context.Background()

	sub := c.Get("sub").(string)
	subid, err := uuid.Parse(sub)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	queries := db.New(dbConn.Conn)

	userData, err := queries.GetUsernameFromId(ctx, subid)
	if err != nil && err != pgx.ErrNoRows {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, userData)
}

type PutUsernameRequest struct {
	Username string `json:"username"`
}

// TODO validation
func (dbConn *DBConn) PutUsername(c echo.Context) error {
	ctx := context.Background()

	sub := c.Get("sub").(string)
	subid, err := uuid.Parse(sub)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	var putUsernameParams PutUsernameRequest
	err = c.Bind(&putUsernameParams)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	createUsernameParams := db.CreateUsernameParams{
		ID:       subid,
		Username: putUsernameParams.Username,
	}

	queries := db.New(dbConn.Conn)

	err = queries.CreateUsername(ctx, createUsernameParams)
	//duplicate key error
	if err != nil {
		if strings.Contains(err.Error(), `ERROR: duplicate key value violates unique constraint "profiles_username_key" (SQLSTATE 23505)`) {
			return c.String(http.StatusConflict, "username already taken")
		} else {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal server error")
		}
	}

	return c.NoContent(200)
}

// func (dbConn *DBConn) PutOnboarding(c echo.Context) error {
// 	ctx := context.Background()

// 	sub := c.Get("sub").(string)
// 	subid, err := uuid.Parse(sub)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "bad request")
// 	}

// 	queries := db.New(dbConn.Conn)

// 	err = queries.UpdateOnboarding(ctx, subid)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "bad request")
// 	}

// 	return c.NoContent(200)
// }

// type ChangeUsernameParams struct {
// 	Username string `json:"username"`
// }

// func (dbConn *DBConn) ChangeUsername(c echo.Context) error {
// 	ctx := context.Background()

// 	sub := c.Get("sub").(string)
// 	subid, err := uuid.Parse(sub)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "bad request")
// 	}

// 	newUsername := new(ChangeUsernameParams)
// 	if err := c.Bind(newUsername); err != nil {
// 		return c.JSON(http.StatusBadRequest, "bad request")
// 	}

// 	queries := db.New(dbConn.Conn)

// 	changeUsernameParams := db.ChangeUsernameParams{
// 		Username: newUsername.Username,
// 		ID:       subid,
// 	}

// 	err = queries.ChangeUsername(ctx, changeUsernameParams)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, "bad request")
// 	}

// 	return c.NoContent(200)
// }
