package main

import (
	"log"
	"net/http"

	"github.com/TonyLeCode/gungi.go/server/api"
	"github.com/TonyLeCode/gungi.go/server/middleware"
	"github.com/TonyLeCode/gungi.go/server/utils"
	"github.com/labstack/echo/v4"
	// "nhooyr.io/websocket"
)

// func game(c echo.Context) error {
// 	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	defer ws.Close()

// 	for {
// 		_, message, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		log.Println("Message received: ", message)
// 		if true {
// 			break
// 		}
// 		// err = ws.ReadJSON()
// 	}

// 	return nil
// }

func main() {
	// board := gungi.Board{}
	// board.InitializeBoard()
	// board.PrintBoard()
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatalln("Cannot load config")
	}

	db, err := api.NewConnection(config.DB_SOURCE)
	if err != nil {
		log.Fatalln(err)
	}

	e := echo.New()

	// e.Use(middleware.VerifySupabaseTokenMiddleware)
	verify := e.Group("", middleware.VerifySupabaseTokenMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world")
	})

	// verify.GET("/websocket", game)
	verify.GET("/getongoinggamelist", db.GetOngoingGameList)

	e.GET("/getgame/:id", db.GetGame)

	// e.POST("/user/register", )

	e.Logger.Fatal(e.Start("localhost:5080"))
}
