package main

import (
	"log"
	"net/http"

	"github.com/TonyLeCode/gungi.go/server/gungi"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

func game(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Message received: ", message)
		if true {
			break
		}
		// err = ws.ReadJSON()
	}

	return nil
}

func main() {
	board := gungi.Board{}
	board.InitializeBoard()
	board.PrintBoard()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world")
	})

	e.GET("/websocket", game)

	// e.POST("/user/register", )

	e.Logger.Fatal(e.Start("localhost:5080"))
}
