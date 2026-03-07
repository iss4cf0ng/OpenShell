package session

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Manager struct {
	Sessions map[string]*Session
	mutex sync.Mutex
	counter int
}

func NewManager() *Manager {
	return &Manager {
		Sessions: make(map[string]*Session),
	}
}

func (m *Manager) AddSession(conn net.Conn) *Session {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.counter++

	id := fmt.Sprintf("%d", m.counter)

	session := &Session {
		ID: id,
		Conn: conn,
		RemoteIP: conn.RemoteAddr().String(),
		CreatedAt: time.Now(),
	}

	m.Sessions[id] = session

	return session
}

func (m *Manager) RemoveSession(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.Sessions, id)
}

func (m *Manager) ListSessions() []*Session {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	list := []*Session{}

	for _, s := range m.Sessions {
		list = append(list, s)
	}

	return list
}