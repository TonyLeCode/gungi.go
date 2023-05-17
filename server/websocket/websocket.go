package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

type PlayerMove struct {
	Type string `json:"type"`
	Data struct {
		MoveType  int `json:"moveType"`
		FromCoord int `json:"fromCoord"`
		ToCoord   int `json:"toCoord"`
	}
}

func GameRoom(c echo.Context) error {
	w := c.Response()
	r := c.Request()
	upgrader := websocket.NewUpgrader()
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	type Msg struct {
		Message string `json:"message"`
		Data    string `json:"data"`
	}
	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		msg := Msg{
			Message: "this works!",
			Data:    string(data),
		}
		m, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
		c.WriteMessage(messageType, m)
	})

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	conn.OnClose(func(c *websocket.Conn, err error) {
		log.Println("OnClose: ", c.RemoteAddr().String(), err)
	})
	log.Println("OnOpen: ", conn.RemoteAddr().String())
	return nil
}
