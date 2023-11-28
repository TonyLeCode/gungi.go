package api

import (

	// "github.com/whitemonarch/gungi-server/server/db"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	db "github.com/whitemonarch/gungi-server/server/db/sqlc"
)

type UserDataRequest struct {
	ID string `json:"id"`
}

func (dbConn *DBConn) GetOnboarding(c echo.Context) error {
	ctx := context.Background()

	sub := c.Get("sub").(string)
	subid, err := uuid.Parse(sub)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	queries := db.New(dbConn.Conn)

	userData, err := queries.GetOnboarding(ctx, subid)
	if err != nil {
		log.Println(userData)
		log.Println(err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, userData)
}

func (dbConn *DBConn) PutOnboarding(c echo.Context) error {
	ctx := context.Background()

	sub := c.Get("sub").(string)
	subid, err := uuid.Parse(sub)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	queries := db.New(dbConn.Conn)

	err = queries.UpdateOnboarding(ctx, subid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	return c.NoContent(200)
}

type ChangeUsernameParams struct {
	Username string `json:"username"`
}

func (dbConn *DBConn) ChangeUsername(c echo.Context) error {
	ctx := context.Background()

	sub := c.Get("sub").(string)
	subid, err := uuid.Parse(sub)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	newUsername := new(ChangeUsernameParams)
	if err := c.Bind(newUsername); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	queries := db.New(dbConn.Conn)

	changeUsernameParams := db.ChangeUsernameParams{
		Username: newUsername.Username,
		ID:       subid,
	}

	err = queries.ChangeUsername(ctx, changeUsernameParams)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	return c.NoContent(200)
}
