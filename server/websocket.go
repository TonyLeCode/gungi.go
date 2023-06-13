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
	db "github.com/TonyLeCode/gungi.go/server/db/sqlc"
	"github.com/TonyLeCode/gungi.go/server/gungi"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
)

type client struct {
	username string
	route    string
	gameID   string
}

type clientList struct {
	clients map[*melody.Session]client
	mutex   sync.Mutex
}

func (c *clientList) AddClient(s *melody.Session, username string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.clients[s] = client{username: username}
}

func (c *clientList) RemoveClient(s *melody.Session) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.clients, s)
}

func (c *clientList) ChangeRoute(s *melody.Session, route string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, ok := c.clients[s]; ok {
		entry.route = route
		c.clients[s] = entry
	} else {
		return errors.New("unable to change route")
	}
	return nil
}

type Room struct {
	clients map[*melody.Session]bool
	mutex   sync.Mutex
}

func (r *Room) AddClient(s *melody.Session) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.clients[s] = true
}

func (r *Room) RemoveClient(s *melody.Session) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.clients, s)
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

func RedisRoomExists(r *redis.Client, roomKey string) error {
	ctx := context.Background()
	exists, err := r.Exists(ctx, roomKey).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		_, err = r.Do(ctx, "JSON.SET", roomKey, "$", "[]").Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func ws(m *melody.Melody, dbs *api.DBConn) echo.HandlerFunc {
	clientList := clientList{clients: make(map[*melody.Session]client)}
	gameRooms := make(map[string]*Room)
	mutex := sync.Mutex{}

	m.HandleConnect(func(s *melody.Session) {
		log.Println("Connected")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		// log.Println(s.Keys)
		clientList.RemoveClient(s)
		mutex.Lock()
		gameID := clientList.clients[s].gameID
		if gameID != "" {
			delete(gameRooms[gameID].clients, s)
		}
		mutex.Unlock()
		log.Println("Disconnect")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		ctx := context.Background()

		var unmarshal MsgPayload
		err := json.Unmarshal(msg, &unmarshal)
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		// Generic messages that don't require a route
		switch unmarshal.Type {
		case "auth":
			var token string
			err := json.Unmarshal(unmarshal.Payload, &token)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			claims, err := auth.AuthenticateSupabaseToken(token)
			if err != nil {
				log.Println(err)
				// 0 means failed
				authResponse := MsgResponse{
					Type:    "auth",
					Payload: "0",
				}
				payload, err := json.Marshal(authResponse)
				if err != nil {
					log.Println("Error: ", err)
				} else {
					err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
						return q == s
					})
					if err != nil {
						log.Println("Error: ", err)
					}
				}
				return
			}

			metadata := claims["user_metadata"].(map[string]interface{})
			username := metadata["username"].(string)
			clientList.AddClient(s, username)

			// 1 means authenticated
			authResponse := MsgResponse{
				Type:    "auth",
				Payload: "1",
			}
			payload, err := json.Marshal(authResponse)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			err = s.Write(payload)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			return

		case "route":
			var route string
			err = json.Unmarshal(unmarshal.Payload, &route)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			err = clientList.ChangeRoute(s, route)
			if err != nil {
				log.Println("Error: ", err)
				return
			}

			// TODO if gameRoom is empty, then delete
			mutex.Lock()
			clientList.mutex.Lock()
			if entry, ok := clientList.clients[s]; ok {
				if entry.gameID != "" {
					delete(gameRooms[entry.gameID].clients, s)
				}

				entry.gameID = ""
				clientList.clients[s] = entry
			} else {
				return
			}
			clientList.mutex.Unlock()
			mutex.Unlock()

			if route == "roomList" {
				payload, err := getRoomList(ctx, dbs, s, true)
				if err != nil {
					log.Println("Error: ", err)
					return
				}
				err = s.Write(payload)
				if err != nil {
					log.Println("Error: ", err)
					return
				}
			}

			if route == "game" {
			}
			return
		}

		client := clientList.clients[s]
		if client.route == "" {
			return
		}

		switch client.route {
		case "roomList":
			handleRoomListMessages(unmarshal, m, s, dbs, ctx, &clientList)
		case "game":
			err = handleGameMessages(unmarshal, m, s, dbs, ctx, &clientList, &gameRooms, &mutex)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	return func(c echo.Context) error {
		err := m.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			return err
		}
		return nil
	}
}

func getRoomList(ctx context.Context, dbs *api.DBConn, s *melody.Session, checkRoom bool) ([]byte, error) {
	if checkRoom {
		err := RedisRoomExists(dbs.RedisClient, "game_rooms")
		if err != nil {
			log.Println(err)
			return []byte{}, err
		}
	}

	val, err := dbs.RedisClient.Do(ctx, "JSON.GET", "game_rooms").Result()
	if err != nil {
		log.Println("Error: ", err)
		return []byte{}, err
	}

	//TODO don't marshal in function
	roomList := MsgResponse{
		Type:    "roomList",
		Payload: val,
	}
	payload, err := json.Marshal(roomList)
	if err != nil {
		return []byte{}, err
	}
	// err = s.Write(payload)
	// if err != nil {
	// 	log.Println("Error: ", err)
	// 	return err
	// }
	return payload, nil
}

func getGame(dbs *api.DBConn, roomID string) (api.GameWithMoves, error) {
	game, err := dbs.GetGame(roomID)
	if err != nil {
		log.Println("Error: ", err)
		return api.GameWithMoves{}, err
	}

	// msgResponse := MsgResponse{
	// 	Type:    "game",
	// 	Payload: game,
	// }
	return game, nil
}

type Move struct {
	FromPiece int `json:"fromPiece"`
	FromCoord int `json:"fromCoord"`
	MoveType  int `json:"moveType"`
	ToCoord   int `json:"toCoord"`
}

func handleGameMessages(msg MsgPayload, m *melody.Melody, s *melody.Session, dbs *api.DBConn, ctx context.Context, clientList *clientList, gameRooms *map[string]*Room, mutex *sync.Mutex) error {
	switch msg.Type {
	case "joinGame":
		var roomID string
		err := json.Unmarshal(msg.Payload, &roomID)
		if err != nil {
			return err
		}

		clientList.mutex.Lock()
		if entry, ok := clientList.clients[s]; ok {
			entry.gameID = roomID
			clientList.clients[s] = entry
		} else {
			return errors.New("could not find gameID")
		}
		clientList.mutex.Unlock()

		mutex.Lock()
		room, ok := (*gameRooms)[roomID]
		if !ok {
			room = &Room{clients: make(map[*melody.Session]bool)}
			(*gameRooms)[roomID] = room
		}
		mutex.Unlock()
		room.AddClient(s)

	case "ready":
	case "makeMove":
		gameID := clientList.clients[s].gameID
		if gameID == "" {
			return errors.New("could not find gameID")
		}

		game, err := getGame(dbs, gameID)
		if err != nil {
			return err
		}

		var move Move
		err = json.Unmarshal(msg.Payload, &move)
		if err != nil {
			return err
		}

		newBoard := gungi.Board{}
		err = newBoard.SetBoardFromFen(game.CurrentState)
		if err != nil {
			return err
		}
		newBoard.SetMoveHistory(game.History.String)
		newBoard.PrintBoard()
		log.Println(move)
		log.Println("fen: ", newBoard.BoardToFen())
		// log.Println(newBoard.TurnColor)
		// log.Println(gungi.GetColor(move.FromPiece))
		isValid := newBoard.MakeMove(move.FromPiece, gungi.IndexToSquare(move.FromCoord), move.MoveType, gungi.IndexToSquare(move.ToCoord))
		if !isValid {
			return errors.New("invalid move")
		}
		newBoard.PrintBoard()
		log.Println("fen: ", newBoard.BoardToFen())

		game.CurrentState = newBoard.BoardToFen()
		game.History.String = newBoard.SerializeHistory()
		game.History.Valid = true

		makeMoveParams := db.MakeMoveParams{
			ID:           game.ID,
			CurrentState: game.CurrentState,
			History:      game.History,
		}
		queries := db.New(dbs.PostgresDB)
		err = queries.MakeMove(ctx, makeMoveParams)
		if err != nil {
			return err
		}

		_, _, _, _, _, legalMoves := newBoard.GenerateLegalMoves()
		correctedLegalMoves := make(map[int][]int)
		for key, element := range legalMoves {
			correctedKey := gungi.SquareToIndex(key)
			for _, index := range element {
				correctedElement := gungi.SquareToIndex(index)
				correctedLegalMoves[correctedKey] = append(correctedLegalMoves[correctedKey], correctedElement)
			}
		}
		game.MoveList = correctedLegalMoves
		log.Println(correctedLegalMoves)

		msgResponse := MsgResponse{
			Type:    "game",
			Payload: game,
		}
		payload, err := json.Marshal(msgResponse)
		if err != nil {
			return err
		}
		for session := range (*gameRooms)[gameID].clients {
			session.Write(payload)
		}
	case "requestUndo":
	case "acceptUndo":
	case "resign":
	}
	return nil
}

func handleRoomListMessages(msg MsgPayload, m *melody.Melody, s *melody.Session, dbs *api.DBConn, ctx context.Context, clientList *clientList) error {
	switch msg.Type {
	case "createRoom":
		// TODO username and validate
		var roomPayload GameRoomType
		err := json.Unmarshal(msg.Payload, &roomPayload)
		if err != nil {
			log.Println("Error: ", err)
			return err
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

		err = RedisRoomExists(dbs.RedisClient, "game_rooms")
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

		_, err = dbs.RedisClient.Do(ctx, "JSON.ARRAPPEND", "game_rooms", "$", marshalled).Result()
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

		payload, err := getRoomList(ctx, dbs, s, true)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
		err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
			return clientList.clients[q].route == "roomList"
		})
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

	case "cancel":
		//TODO verify
		var roomid string
		err := json.Unmarshal(msg.Payload, &roomid)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

		roomList, err := RedisGetGameRooms(dbs.RedisClient)
		if err != nil {
			log.Println("Error: ", err)
			return err
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
			return err
		}
		newRoomList := MsgResponse{
			Type:    "roomList",
			Payload: string(marshalled),
		}
		payload, err := json.Marshal(newRoomList)
		if err != nil {
			log.Println("Error: ", err)
		}
		err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
			return clientList.clients[q].route == "roomList"
		})
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

	case "accept":
		var roomid string
		err := json.Unmarshal(msg.Payload, &roomid)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

		roomList, err := RedisGetGameRooms(dbs.RedisClient)
		if err != nil {
			log.Println("Error: ", err)
			return err
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
			return err
		}

		_, err = dbs.RedisClient.Do(ctx, "JSON.SET", "game_rooms", "$", marshalled).Result()
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
		newRoomList := MsgResponse{
			Type:    "roomList",
			Payload: string(marshalled),
		}
		payload, err := json.Marshal(newRoomList)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
		err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
			return clientList.clients[q].route == "roomList"
		})
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

		// Create new game
		var user1 string
		var user2 string
		sessionUsername := clientList.clients[s].username
		if acceptedRoom.Color == "w" {
			user1 = acceptedRoom.Host
			user2 = sessionUsername
		} else if acceptedRoom.Color == "b" {
			user1 = sessionUsername
			user2 = acceptedRoom.Host
		} else {
			rand.Seed(time.Now().UnixNano())
			randBool := rand.Intn(2) == 0
			if randBool {
				user1 = acceptedRoom.Host
				user2 = sessionUsername
			} else {
				user1 = sessionUsername
				user2 = acceptedRoom.Host
			}
		}
		gameID, err := dbs.CreateGame("9/9/9/9/9/9/9/9/9 9446222122211/9446222122211 w 00", acceptedRoom.Rules, acceptedRoom.Type, user1, user2)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
		acceptResponse := MsgResponse{
			Type:    "accepted",
			Payload: gameID,
		}
		payload, err = json.Marshal(acceptResponse)
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

		err = m.BroadcastFilter(payload, func(q *melody.Session) bool {
			username := clientList.clients[q].username
			return username == user1 || username == user2
		})
		if err != nil {
			log.Println("Error: ", err)
			return err
		}

	}
	return nil
}

func RedisGetGameRooms(r *redis.Client) ([]SerializedGameRoom, error) {
	ctx := context.Background()
	err := RedisRoomExists(r, "game_rooms")
	if err != nil {
		log.Println("Error: ", err)
		return []SerializedGameRoom{}, err
	}

	val, err := r.Do(ctx, "JSON.GET", "game_rooms").Result()
	if err != nil {
		log.Println("Error: ", err)
		return []SerializedGameRoom{}, err
	}

	var roomList []SerializedGameRoom
	err = json.Unmarshal([]byte(val.(string)), &roomList)
	if err != nil {
		log.Println("Error: ", err)
		return []SerializedGameRoom{}, err
	}

	return roomList, nil
}
