package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/TonyLeCode/gungi.go/server/api"
	"github.com/TonyLeCode/gungi.go/server/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
)

type client struct {
	socket   *melody.Session
	username string
	route    string
}

type Room struct {
	clients map[client]bool
	mutex   sync.Mutex
}

func (r *Room) AddClient(c client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.clients[c] = true
}

func (r *Room) RemoveClient(c client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.clients, c)
}

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

type MsgAuth struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

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

func ws(m *melody.Melody, dbs *api.DBConn) echo.HandlerFunc {
	m.HandleConnect(func(s *melody.Session) {
		log.Println("Connected")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		// log.Println(s.Keys)
		log.Println("Disconnect")
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
			// TODO username and validate
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

			InitializeRooms(dbs.RedisClient)

			_, err = dbs.RedisClient.Do(ctx, "JSON.ARRAPPEND", "game_rooms", "$", marshalled).Result()
			if err != nil {
				log.Println("Error: ", err)
			}

			val, err := dbs.RedisClient.Do(ctx, "JSON.GET", "game_rooms").Result()
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
			} else {
				m.Broadcast(payload)
			}

		case "auth":

			var authMsg MsgAuth
			err := json.Unmarshal(msg, &authMsg)
			if err != nil {
				log.Println("Error: ", err)
			}

			claims, err := auth.AuthenticateSupabaseToken(authMsg.Payload)
			if err != nil {
				log.Println(err)
				authResponse := MsgResponse{
					Type:    "auth",
					Payload: "0",
				}
				payload, err := json.Marshal(authResponse)
				if err != nil {
					log.Println("Error: ", err)
				} else {
					log.Println("payload: ", payload)
					err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
						return q == s
					})
					if err != nil {
						log.Println("Error: ", err)
					}
				}
			} else {
				metadata := claims["user_metadata"].(map[string]interface{})
				username := metadata["username"].(string)
				s.Set("username", username)

				ctx := context.Background()
				InitializeRooms(dbs.RedisClient)
				val, err := dbs.RedisClient.Do(ctx, "JSON.GET", "game_rooms").Result()
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
					username1, _ := q.Get("username")
					username2, _ := s.Get("username")
					return username1 == username2
				})
				if err != nil {
					log.Println("Error: ", err)
				}

				authResponse := MsgResponse{
					Type:    "auth",
					Payload: "1",
				}
				payload, err = json.Marshal(authResponse)
				if err != nil {
					log.Println("Error: ", err)
				} else {
					err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
						username1, _ := q.Get("username")
						username2, _ := s.Get("username")
						return username1 == username2
					})
					if err != nil {
						log.Println("Error: ", err)
					}
				}
			}

		case "cancel":
			//TODO verify
			var roomid string
			err = json.Unmarshal(unmarshal.Payload, &roomid)
			if err != nil {
				log.Println("Error: ", err)
			}

			val, err := dbs.RedisClient.Do(ctx, "JSON.GET", "game_rooms").Result()
			if err != nil {
				log.Println("Error: ", err)
			}

			var roomList []SerializedGameRoom
			err = json.Unmarshal([]byte(val.(string)), &roomList)
			if err != nil {
				log.Println("Error: ", err)
			}

			var filteredRoomList []SerializedGameRoom
			for _, room := range roomList {
				if room.RoomID.String() != roomid {
					filteredRoomList = append(filteredRoomList, room)
				}
			}

			marshalled, err := json.Marshal(filteredRoomList)
			if err != nil {
				log.Println("Error: ", err)
			}

			_, err = dbs.RedisClient.Do(ctx, "JSON.SET", "game_rooms", "$", marshalled).Result()
			if err != nil {
				log.Println("Error: ", err)
			} else {
				newRoomList := MsgResponse{
					Type:    "roomList",
					Payload: string(marshalled),
				}
				payload, err := json.Marshal(newRoomList)
				if err != nil {
					log.Println("Error: ", err)
				} else {
					m.Broadcast(payload)
				}
			}

		case "accept":
			var roomid string
			err = json.Unmarshal(unmarshal.Payload, &roomid)
			if err != nil {
				log.Println("Error: ", err)
			}

			val, err := dbs.RedisClient.Do(ctx, "JSON.GET", "game_rooms").Result()
			if err != nil {
				log.Println("Error: ", err)
			}

			var roomList []SerializedGameRoom
			err = json.Unmarshal([]byte(val.(string)), &roomList)
			if err != nil {
				log.Println("Error: ", err)
			}

			var acceptedRoom SerializedGameRoom
			var filteredRoomList []SerializedGameRoom
			for _, room := range roomList {
				if room.RoomID.String() != roomid {
					filteredRoomList = append(filteredRoomList, room)
				} else {
					acceptedRoom = room
				}
			}

			marshalled, err := json.Marshal(filteredRoomList)
			if err != nil {
				log.Println("Error: ", err)
			}

			_, err = dbs.RedisClient.Do(ctx, "JSON.SET", "game_rooms", "$", marshalled).Result()
			if err != nil {
				log.Println("Error: ", err)
			} else {
				newRoomList := MsgResponse{
					Type:    "roomList",
					Payload: string(marshalled),
				}
				payload, err := json.Marshal(newRoomList)
				if err != nil {
					log.Println("Error: ", err)
				} else {
					m.Broadcast(payload)
				}
			}

			// Create new game
			var user1 string
			var user2 string
			sessionUsername, _ := s.Get("username")
			if acceptedRoom.Color == "w" {
				user1 = acceptedRoom.Host
				user2 = sessionUsername.(string)
			} else if acceptedRoom.Color == "b" {
				user1 = sessionUsername.(string)
				user2 = acceptedRoom.Host
			} else {
				rand.Seed(time.Now().UnixNano())
				randBool := rand.Intn(2) == 0
				if randBool {
					user1 = acceptedRoom.Host
					user2 = sessionUsername.(string)
				} else {
					user1 = sessionUsername.(string)
					user2 = acceptedRoom.Host
				}
			}
			gameID, err := dbs.CreateGame("9/9/9/9/9/9/9/9/9 9446222122211/9446222122211 w", acceptedRoom.Rules, acceptedRoom.Type, user1, user2)
			if err != nil {
				log.Println("Error: ", err)
			} else {
				acceptResponse := MsgResponse{
					Type:    "accepted",
					Payload: gameID,
				}
				payload, err := json.Marshal(acceptResponse)
				if err != nil {
					log.Println("Error: ", err)
				}

				err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
					username, _ := q.Get("username")
					return username == user1 || username == user2
				})
				if err != nil {
					log.Println("Error: ", err)
				}
			}

		}

		// m.Broadcast(msg)
	})

	return func(c echo.Context) error {
		err := m.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			return err
		}
		return nil
	}
}

// func handleAuth(msg []byte) error {
// 	var authMsg MsgAuth
// 	err := json.Unmarshal(msg, &authMsg)
// 	if err != nil {
// 		log.Println("Error: ", err)
// 	}

// 	claims, err := auth.AuthenticateSupabaseToken(authMsg.Payload)
// 	if err != nil {
// 		log.Println(err)
// 		authResponse := MsgResponse{
// 			Type:    "auth",
// 			Payload: "0",
// 		}
// 		payload, err := json.Marshal(authResponse)
// 		if err != nil {
// 			log.Println("Error: ", err)
// 		} else {
// 			log.Println("payload: ", payload)
// 			err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
// 				return q == s
// 			})
// 			if err != nil {
// 				log.Println("Error: ", err)
// 			}
// 		}
// 	} else {
// 		metadata := claims["user_metadata"].(map[string]interface{})
// 		username := metadata["username"].(string)
// 		s.Set("username", username)

// 		ctx := context.Background()
// 		InitializeRooms(dbs.RedisClient)
// 		val, err := dbs.RedisClient.Do(ctx, "JSON.GET", "game_rooms").Result()
// 		if err != nil {
// 			log.Println("Error: ", err)
// 		}

// 		roomList := MsgResponse{
// 			Type:    "roomList",
// 			Payload: val,
// 		}
// 		payload, err := json.Marshal(roomList)
// 		if err != nil {
// 			log.Println("Error: ", err)
// 		}

// 		err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
// 			username1, _ := q.Get("username")
// 			username2, _ := s.Get("username")
// 			return username1 == username2
// 		})
// 		if err != nil {
// 			log.Println("Error: ", err)
// 		}

// 		authResponse := MsgResponse{
// 			Type:    "auth",
// 			Payload: "1",
// 		}
// 		payload, err = json.Marshal(authResponse)
// 		if err != nil {
// 			log.Println("Error: ", err)
// 		} else {
// 			err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
// 				username1, _ := q.Get("username")
// 				username2, _ := s.Get("username")
// 				return username1 == username2
// 			})
// 			if err != nil {
// 				log.Println("Error: ", err)
// 			}
// 		}
// 	}
// 	return nil
// }
