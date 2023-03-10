package main

import (
	"net/http"

	"github.com/TonyLeCode/gungi.go/gungi"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

var rooms = make(map[string][]*websocket.Conn)

func game(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	return nil
}

func main() {
	board := gungi.Board{}
	board.InitializeBoard()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world")
	})
	e.GET("/websocket", game)

	e.Logger.Fatal(e.Start("localhost:5080"))
}
