package tcp

import (
	"github.com/google/uuid"
	"time"
)

// Session struct
type Session struct {
	sid      string
	uid      string
	conn     *Conn
	lastTime int64
	extraMap map[string]interface{}
}

// NewSession create a new session
func NewSession(conn *Conn) *Session {
	id := uuid.New()
	session := &Session{
		sid:      id.String(),
		uid:      "",
		conn:     conn,
		lastTime: time.Now().Unix(),
		extraMap: make(map[string]interface{}),
	}

	return session
}

// GetSessionID get a session ID
func (s *Session) GetSessionID() string {
	return s.sid
}

// BindUserID bind a user ID to session
func (s *Session) BindUserID(uid string) {
	s.uid = uid
}

// GetUserID get user ID
func (s *Session) GetUserID() string {
	return s.uid
}

// GetConn get Conn pointer
func (s *Session) GetConn() *Conn {
	return s.conn
}

// SetConn set a Conn to session
func (s *Session) SetConn(conn *Conn) {
	s.conn = conn
}

// UpdateTime update the message last time
func (s *Session) UpdateTime() {
	s.lastTime = time.Now().Unix()
}

// GetExtraMap get the extra data
func (s *Session) GetExtraMap(key string) interface{} {
	if v, ok := s.extraMap[key]; ok {
		return v
	}

	return nil
}

// SetExtraMap set the extra data
func (s *Session) SetExtraMap(key string, value interface{}) {
	s.extraMap[key] = value
}
