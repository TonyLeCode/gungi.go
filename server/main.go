package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TonyLeCode/gungi.go/server/api"
	"github.com/TonyLeCode/gungi.go/server/middleware"
	"github.com/TonyLeCode/gungi.go/server/utils"
	"github.com/TonyLeCode/gungi.go/server/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp"
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

	e.GET("/room", websocket.GameRoom)

	// e.POST("/user/register", )

	// e.Logger.Fatal(e.Start("localhost:5080"))
	svr := nbhttp.NewServer(nbhttp.Config{
		Network: "tcp",
		Addrs:   []string{"localhost:8080"},
	}, e, nil)
	err = svr.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v\n", err)
	}

	log.Println("serving [labstack/echo] on [nbio]")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	svr.Shutdown(ctx)
}
