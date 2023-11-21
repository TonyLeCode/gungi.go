package api

import (

	// "github.com/TonyLeCode/gungi.go/server/db"
	"context"
	"log"
	"net/http"

	db "github.com/TonyLeCode/gungi.go/server/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type UserDataRequest struct {
	ID string `json:"id"`
}

func (dbConn *DBConn) GetUserData(c echo.Context) error {
	ctx := context.Background()
	request := new(UserDataRequest)
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	id, err := uuid.Parse(request.ID)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	log.Println(id)

	queries := db.New(dbConn.Conn)

	userData, err := queries.GetUserData(ctx, id)
	if err != nil {
		log.Println(userData)
		log.Println(err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, userData)
}
