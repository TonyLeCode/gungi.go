package api

import (

	// "github.com/TonyLeCode/gungi.go/server/db"
	"context"
	"log"
	"net/http"

	db "github.com/TonyLeCode/gungi.go/server/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func (dbConn *DBConn) PostOnboarding(c echo.Context) error {
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
