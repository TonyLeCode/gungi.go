package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
)

type PlayerMove struct {
	Type string `json:"type"`
	Data struct {
		MoveType  int `json:"moveType"`
		FromCoord int `json:"fromCoord"`
		ToCoord   int `json:"toCoord"`
	}
}

type GameRoomType struct {
	Host        string `json:"host"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Rules       string `json:"rules"`
}

type SerializedGameRoom struct {
	GameRoomType
	RoomID uuid.UUID `json:"roomid"`
}

type MsgPayload struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MsgResponse struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

//	type SerializedGameRoomMsg struct {
//		Type    string `json:"type"`
//		Payload string `json:"payload"`
//	}
func InitializeRooms(r *redis.Client) error {
	ctx := context.Background()
	_, err := r.Do(ctx, "JSON.GET", "game_rooms").Result()
	if errors.Is(err, redis.Nil) {
		_, err = r.Do(ctx, "JSON.SET", "game_rooms", "$", "[]").Result()
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
	} else if err != nil {
		log.Println("Error: ", err)
		return err
	}
	return nil
}

func GameRoom(m *melody.Melody, r *redis.Client) echo.HandlerFunc {
	m.HandleConnect(func(s *melody.Session) {
		ctx := context.Background()
		InitializeRooms(r)
		val, err := r.Do(ctx, "JSON.GET", "game_rooms").Result()
		if err != nil {
			log.Println("Error: ", err)
		}

		roomList := MsgResponse{
			Type:    "roomList",
			Payload: val,
		}
		payload, err := json.Marshal(roomList)
		if err != nil {
			log.Println("Error: ", err)
		}

		err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
			return q == s
		})
		if err != nil {
			log.Println("Error: ", err)
		}

		log.Println("Connected")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		ctx := context.Background()

		var unmarshal MsgPayload
		err := json.Unmarshal(msg, &unmarshal)
		if err != nil {
			log.Println("Error: ", err)
		}

		switch unmarshal.Type {
		case "createRoom":
			var roomPayload GameRoomType
			err = json.Unmarshal(unmarshal.Payload, &roomPayload)
			if err != nil {
				log.Println("Error: ", err)
			}
			newRoom := SerializedGameRoom{
				RoomID: uuid.New(),
				GameRoomType: GameRoomType{
					Host:        roomPayload.Host,
					Description: roomPayload.Description,
					Type:        roomPayload.Type,
					Color:       roomPayload.Color,
					Rules:       roomPayload.Rules,
				},
			}
			marshalled, err := json.Marshal(newRoom)
			if err != nil {
				log.Println("Error: ", err)
			}

			InitializeRooms(r)

			val, err := r.Do(ctx, "JSON.ARRAPPEND", "game_rooms", "$", marshalled).Result()
			if err != nil {
				log.Println("Error: ", err)
			}

			val, err = r.Do(ctx, "JSON.GET", "game_rooms").Result()
			if err != nil {
				log.Println("Error: ", err)
			}

			roomList := MsgResponse{
				Type:    "roomList",
				Payload: val,
			}
			payload, err := json.Marshal(roomList)
			if err != nil {
				log.Println("Error: ", err)
			}
			m.Broadcast(payload)
		}

		m.Broadcast(msg)
	})

	return func(c echo.Context) error {
		err := m.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Successfully Connected")
		return nil
	}
}

// func GameRoom(c echo.Context) error {
// 	w := c.Response()
// 	r := c.Request()
// 	upgrader := websocket.NewUpgrader()
// 	upgrader.CheckOrigin = func(r *http.Request) bool {
// 		return true
// 	}

// 	type Msg struct {
// 		Message string `json:"message"`
// 		Data    string `json:"data"`
// 	}
// 	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
// 		msg := Msg{
// 			Message: "this works!",
// 			Data:    string(data),
// 		}
// 		m, err := json.Marshal(msg)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		c.WriteMessage(messageType, m)
// 	})

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	conn.OnClose(func(c *websocket.Conn, err error) {
// 		log.Println("OnClose: ", c.RemoteAddr().String(), err)
// 	})
// 	log.Println("OnOpen: ", conn.RemoteAddr().String())
// 	return nil
// }
