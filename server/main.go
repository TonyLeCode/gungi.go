package main

import (
	"log"
	"net/http"

	"github.com/TonyLeCode/gungi.go/server/api"
	"github.com/TonyLeCode/gungi.go/server/middleware"
	"github.com/TonyLeCode/gungi.go/server/utils"
	"github.com/TonyLeCode/gungi.go/server/websocket"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
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

	dbs := api.DBConn{}
	dbs.PostgresConnect(config.DB_SOURCE)
	defer dbs.PostgresDB.Close()
	dbs.RedisConnect(config.REDIS_CONN_STRING)
	defer dbs.RedisClient.Close()

	websocket.InitializeRooms(dbs.RedisClient)

	// postgresDB, err := api.NewConnection(config.DB_SOURCE)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })

	// dbs := &db.DBs{
	// 	PostgresDB:  postgresDB,
	// 	RedisClient: redisClient,
	// }

	e := echo.New()
	m := melody.New()
	m.Config.MaxMessageSize = 1024

	// e.Use(middleware.VerifySupabaseTokenMiddleware)
	verify := e.Group("", middleware.VerifySupabaseTokenMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world")
	})

	// verify.GET("/websocket", game)
	verify.GET("/getongoinggamelist", dbs.GetOngoingGameList)

	e.GET("/getgame/:id", dbs.GetGame)

	e.GET("/room", websocket.GameRoom(m, &dbs))

	// e.POST("/user/register", )

	e.Logger.Fatal(e.Start("localhost:5080"))
}
