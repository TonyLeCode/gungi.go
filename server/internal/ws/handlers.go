package ws

import (
	"encoding/json"
	"log"

	"github.com/olahol/melody"
	"github.com/whitemonarch/gungi-server/server/internal/auth"
)

// auth
// leaveGame
// makeGameMove
// gameResign
// requestGameUndo
// responseGameUndo
// completeGameUndo
// joinPlay
// leavePlay
// createPlayRoom
// acceptPlayRoom
// cancelPlayRoom

func authHandler(s *melody.Session, m *melody.Melody, unmarshal MsgPayload) {
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
	id := claims["sub"].(string)
	s.Set("username", username)
	s.Set("id", id)

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
}
