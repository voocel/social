package gateway

import (
	"errors"
	"github.com/google/uuid"
	"social/pkg/log"
	"social/pkg/network"
	"sync"
)

type Session struct {
	sid  string
	uid  int64
	conn network.Conn
	rw   sync.RWMutex
}

func newSession(conn network.Conn) *Session {
	return &Session{
		sid:  uuid.New().String(),
		conn: conn,
	}
}

func (s *Session) Reset() {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.uid = 0
	s.conn = nil
}

func (s *Session) CID() int64 {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.Cid()
}

func (s *Session) UID() int64 {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.uid
}

func (s *Session) Bind(uid int64) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.uid = uid
	s.conn.Bind(uid)
}

func (s *Session) Unbind(uid int64) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.uid = 0
	s.conn.Bind(0)
}

func (s *Session) Close() error {
	s.rw.Lock()
	defer s.rw.Unlock()

	return s.conn.Close()
}

func (s *Session) LocalIP() string {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.LocalIP()
}

func (s *Session) LocalAddr() string {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.LocalAddr()
}

func (s *Session) RemoteIP() string {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.RemoteIP()
}

func (s *Session) RemoteAddr() string {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.RemoteAddr()
}

// Send 发送消息（同步）
func (s *Session) Send(msg []byte) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.Send(msg)
}

// Push 发送消息（异步）
func (s *Session) Push(msg []byte) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.conn.Push(msg)
}

type SessionGroup struct {
	rw         sync.RWMutex
	uidSession map[int64]*Session
	cidSession map[int64]*Session
}

func NewSessionGroup() *SessionGroup {
	return &SessionGroup{
		rw:         sync.RWMutex{},
		uidSession: make(map[int64]*Session),
		cidSession: make(map[int64]*Session),
	}
}

func (g *SessionGroup) AddSession(s *Session) {
	g.rw.Lock()
	defer g.rw.Unlock()
	g.cidSession[s.CID()] = s
	if uid := s.UID(); uid > 0 {
		g.uidSession[uid] = s
	}
}

func (g *SessionGroup) GetSessionByCid(cid int64) (*Session, error) {
	g.rw.RLock()
	defer g.rw.RUnlock()
	sess, ok := g.cidSession[cid]
	if !ok {
		return nil, errors.New("cid session not found")
	}
	return sess, nil
}

func (g *SessionGroup) GetSessionByUid(uid int64) (*Session, error) {
	g.rw.RLock()
	defer g.rw.RUnlock()
	sess, ok := g.uidSession[uid]
	if !ok {
		return nil, errors.New("uid session not found")
	}
	return sess, nil
}

func (g *SessionGroup) PushByUid(uid int64, msg []byte) error {
	s, err := g.GetSessionByUid(uid)
	if err != nil {
		return err
	}
	return s.Push(msg)
}

func (g *SessionGroup) PushByCid(cid int64, msg []byte) error {
	s, err := g.GetSessionByCid(cid)
	if err != nil {
		return err
	}
	return s.Push(msg)
}

func (g *SessionGroup) Broadcast(msg []byte) int {
	var n int
	for _, session := range g.cidSession {
		err := session.Push(msg)
		if err != nil {
			log.Errorf("Broadcast push err: %v", err)
			continue
		}
		n++
	}
	return n
}

func (g *SessionGroup) RemoveByCid(cid int64) error {
	g.rw.Lock()
	defer g.rw.Unlock()
	sess, ok := g.cidSession[cid]
	if !ok {
		return errors.New("cid session not found")
	}
	if uid := sess.UID(); uid > 0 {
		delete(g.uidSession, uid)
	}
	delete(g.cidSession, cid)
	return nil
}

func (g *SessionGroup) RemoveByUid(uid int64) error {
	g.rw.Lock()
	defer g.rw.Unlock()
	sess, ok := g.uidSession[uid]
	if !ok {
		return errors.New("uid session not found")
	}
	delete(g.cidSession, sess.CID())
	delete(g.uidSession, uid)
	return nil
}
