package ws

import (
	"sync"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

const ROOM_KEY = "roomlist"

type Sessions map[*melody.Session]*User
type User struct {
	ID        uuid.UUID
	GameID    uuid.UUID
	Spectator bool
	Unsub     func()
}

func (ss *Sessions) AddUser(s *melody.Session, id uuid.UUID) {
	(*ss)[s] = &User{
		ID:        id,
		Unsub:     nil,
		GameID:    uuid.Nil,
		Spectator: false,
	}
}
func (ss *Sessions) AddSpectator(s *melody.Session) {
	(*ss)[s] = &User{
		ID:        uuid.Nil,
		Unsub:     nil,
		GameID:    uuid.Nil,
		Spectator: true,
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

func (ss *Sessions) ChangeGame(s *melody.Session, gameID uuid.UUID) {
	if _, ok := (*ss)[s]; !ok {
		return
	}
	(*ss)[s].GameID = gameID
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
	l.listeners[eventKey] = append(l.listeners[eventKey], s)
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

func (l *Listeners) AddListenerGame(s *melody.Session, gameID uuid.UUID) func() {
	eventKey := "game-" + gameID.String()
	unsub := l.addListener(eventKey, s)
	return unsub
}

func (l *Listeners) RemoveListenerGame(s *melody.Session, gameID uuid.UUID) {
	eventKey := "game-" + gameID.String()
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

func (l *Listeners) EmitGameMsg(msg []byte, gameID uuid.UUID) error {
	eventKey := "game-" + gameID.String()
	err := l.EmitMsg(msg, eventKey)
	if err != nil {
		return err
	}
	return nil
}

func (l *Listeners) EmitGameMsgFilter(msg []byte, gameID uuid.UUID, fn func(s *melody.Session) bool) error {
	eventKey := "game-" + gameID.String()
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
