package ws

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

const ROOM_KEY = "roomlist"

type Sessions map[*melody.Session]*User
type User struct {
	ID           uuid.UUID
	GamePublicID string
	Spectator    bool
	Unsub        func()
}

func (ss *Sessions) AddUser(s *melody.Session, id uuid.UUID) {
	(*ss)[s].ID = id
	(*ss)[s].Spectator = false
}
func (ss *Sessions) AddSpectator(s *melody.Session) {
	(*ss)[s] = &User{
		ID:           uuid.Nil,
		Unsub:        nil,
		GamePublicID: "",
		Spectator:    true,
	}
}

func (ss *Sessions) ChangeUnsub(s *melody.Session, unsub func()) {
	if _, ok := (*ss)[s]; !ok {
		return
	}
	if (*ss)[s].Unsub != nil {
		(*ss)[s].Unsub()
	}
	(*ss)[s].Unsub = unsub
}

func (ss *Sessions) Unsub(s *melody.Session) {
	if _, ok := (*ss)[s]; !ok {
		return
	}
	if (*ss)[s].Unsub != nil {
		(*ss)[s].Unsub()
	}
	(*ss)[s].Unsub = nil
}

func (ss *Sessions) ChangeGame(s *melody.Session, gamePublicID string) {
	log.Println("changing gameid to: ", gamePublicID)
	if _, ok := (*ss)[s]; !ok {
		log.Println("not okay")
		return
	}
	(*ss)[s].GamePublicID = gamePublicID
}

func (ss *Sessions) RemoveUser(s *melody.Session) {
	if (*ss)[s] == nil {
		return
	}
	if (*ss)[s].Unsub != nil {
		(*ss)[s].Unsub()
	}
	delete(*ss, s)
}

type Listeners struct {
	listeners map[string][]*melody.Session
	mu        sync.RWMutex
}

func (l *Listeners) addListener(eventKey string, s *melody.Session) func() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.listeners == nil {
		l.listeners = make(map[string][]*melody.Session)
	}

	if _, ok := l.listeners[eventKey]; !ok {
		l.listeners[eventKey] = make([]*melody.Session, 0)
	}
	log.Println("appending", s, l.listeners[eventKey])
	l.listeners[eventKey] = append(l.listeners[eventKey], s)
	log.Println("appended", l.listeners[eventKey])
	return func() {
		l.removeListener(eventKey, s)
	}
}

func (l *Listeners) removeListener(eventKey string, s *melody.Session) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.listeners[eventKey]; !ok {
		return
	}

	for i, session := range l.listeners[eventKey] {
		if session == s {
			l.listeners[eventKey] = append(l.listeners[eventKey][:i], l.listeners[eventKey][i+1:]...)
			return
		}
	}

	if len(l.listeners[eventKey]) == 0 {
		delete(l.listeners, eventKey)
		return
	}
}

func (l *Listeners) AddListenerRooms(s *melody.Session) func() {
	unsub := l.addListener(ROOM_KEY, s)
	return unsub
}

func (l *Listeners) RemoveListenerRooms(s *melody.Session) {
	l.removeListener(ROOM_KEY, s)
}

func (l *Listeners) AddListenerGame(s *melody.Session, gamePublicID string) func() {
	eventKey := "game-" + gamePublicID
	unsub := l.addListener(eventKey, s)
	return unsub
}

func (l *Listeners) RemoveListenerGame(s *melody.Session, gamePublicID string) {
	eventKey := "game-" + gamePublicID
	l.removeListener(eventKey, s)
}

func (l *Listeners) EmitMsg(msg []byte, eventKey string) error {
	listeners := l.listeners[eventKey]
	for _, s := range listeners {
		if err := s.Write(msg); err != nil {
			return err
		}
	}
	return nil
}

func (l *Listeners) EmitMsgFilter(msg []byte, eventKey string, fn func(s *melody.Session) bool) error {
	listeners := l.listeners[eventKey]
	for _, s := range listeners {
		if !fn(s) {
			continue
		}
		if err := s.Write(msg); err != nil {
			return err
		}
	}
	return nil
}

func (l *Listeners) EmitGameMsg(msg []byte, gamePublicID string) error {
	eventKey := "game-" + gamePublicID
	err := l.EmitMsg(msg, eventKey)
	if err != nil {
		return err
	}
	return nil
}

func (l *Listeners) EmitGameMsgFilter(msg []byte, gamePublicID string, fn func(s *melody.Session) bool) error {
	eventKey := "game-" + gamePublicID
	err := l.EmitMsgFilter(msg, eventKey, fn)
	if err != nil {
		return err
	}
	return nil
}

func (l *Listeners) EmitRoomMsg(msg []byte) error {
	listeners := l.listeners[ROOM_KEY]
	for _, s := range listeners {
		if err := s.Write(msg); err != nil {
			return err
		}
	}
	return nil
}

func (l *Listeners) EmitRoomMsgFilter(msg []byte, fn func(s *melody.Session) bool) error {
	err := l.EmitMsgFilter(msg, ROOM_KEY, fn)
	if err != nil {
		return err
	}
	return nil
}
