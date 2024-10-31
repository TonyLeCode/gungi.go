package ws

import (
	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"

	"github.com/whitemonarch/gungi-server/server/internal/api"
)

type GameRoomType struct {
	Host        string `json:"host"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Rules       string `json:"rules"`
}

type CreateGameRoomRequest struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Rules       string `json:"rules"`
}

type ClientMsg struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ServerMsg struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type makeGameMoveMsg struct {
	FromPiece int `json:"fromPiece"`
	FromCoord int `json:"fromCoord"`
	MoveType  int `json:"moveType"`
	ToCoord   int `json:"toCoord"`
}

func WSHandler(m *melody.Melody, dbConn *api.DBConn) echo.HandlerFunc {
	listeners := Listeners{}
	sessions := Sessions{}

	m.HandleConnect(func(s *melody.Session) {
		log.Println("Connected", s)
		sessions.AddSpectator(s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("Disconnect", s, sessions[s].ID)
		sessions.RemoveUser(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var clientMsg ClientMsg
		err := json.Unmarshal(msg, &clientMsg)
		if err != nil {
			log.Println("Error: ", err)
			return
		}

		if clientMsg.Type != "auth" && sessions[s] == nil {
			return
		}

		sessionHandlers := SessionHandlers{
			Listeners: &listeners,
			Sessions:  &sessions,
			ClientMsg: clientMsg,
			Session:   s,
			DBConn:    dbConn,
		}

		msgHandlers := MsgHandlersMap{
			"auth":             sessionHandlers.handleAuth,
			"joinGame":         sessionHandlers.handleJoinGame,
			"leaveGame":        sessionHandlers.handleLeaveGame,
			"makeGameMove":     sessionHandlers.handleMakeGameMove,
			"gameResign":       sessionHandlers.handleGameResign,
			"requestGameUndo":  sessionHandlers.handleRequestGameUndo,
			"responseGameUndo": sessionHandlers.handleResponseGameUndo,
			"completeGameUndo": sessionHandlers.handleCompleteGameUndo,
			"joinPlay":         sessionHandlers.handleJoinPlay,
			"createPlayRoom":   sessionHandlers.handleCreatePlayRoom,
			"leavePlay":        sessionHandlers.handleLeavePlay,
			"acceptPlayRoom":   sessionHandlers.handleAcceptPlayRoom,
			"cancelPlayRoom":   sessionHandlers.handleCancelPlayRoom,
		}

		if handler, ok := msgHandlers[clientMsg.Type]; ok {
			handler(s)
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

func getGame(dbs *api.DBConn, GamePublicID string) (api.GameWithUndo, error) {
	game, err := dbs.GetGame(GamePublicID)
	if err != nil {
		log.Println("Error: ", err)
		return api.GameWithUndo{}, err
	}
	return game, nil
}
