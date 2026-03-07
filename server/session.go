package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type ShellSession struct {
	ID   string
	Type string

	Conn *websocket.Conn

	Net  net.Conn
	Pty  *os.File
	Cmd  *exec.Cmd

	RemoteAddr string
}

type Session struct {
    ID   string `json:"id"`
    IP   string `json:"ip"`
    Type string `json:"type"`
}

type SessionManager struct {
	sessions map[string]*ShellSession
	mu       sync.Mutex
	counter  int
}

func (sm *SessionManager) Count() int {

	sm.mu.Lock()
	defer sm.mu.Unlock()

	return len(sm.sessions)
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*ShellSession),
	}
}

func (sm *SessionManager) nextID() string {
	sm.counter++
	return fmt.Sprintf("session-%d", sm.counter)
}

func (sm *SessionManager) ListSessions() []Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	list := []Session{}

	for _, s := range sm.sessions {
		list = append(list, Session {
			ID: s.ID,
			Type: s.Type,
			IP: s.RemoteAddr,
		})
	}

	return list;
}

func (sm *SessionManager) DeleteSession(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, id)

	fmt.Println("Session removed:", id)
}

func (sm *SessionManager) CreateLocalShell(conn *websocket.Conn) (*ShellSession, error) {

	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := sm.nextID()

	cmd := exec.Command("/bin/bash")

	ptyFile, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	session := &ShellSession{
		ID:   id,
		Type: "local",
		Pty:  ptyFile,
		Cmd:  cmd,
		Conn: conn,
	}

	sm.sessions[id] = session

	go session.bridgePTY()

	return session, nil
}

func (sm *SessionManager) CreateReverseShell(netConn net.Conn) {

	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := sm.nextID()

	session := &ShellSession{
		ID:         id,
		Type:       "reverse",
		Net:        netConn,
		RemoteAddr: netConn.RemoteAddr().String(),
	}

	sm.sessions[id] = session

	fmt.Println("New reverse shell:", id, netConn.RemoteAddr())

}

func (s *ShellSession) bridgePTY() {

	go func() {
		for {
			_, msg, err := s.Conn.ReadMessage()
			if err != nil {
				s.close()
				return
			}
			s.Pty.Write(msg)
		}
	}()

	buf := make([]byte, 1024)

	for {
		n, err := s.Pty.Read(buf)
		if err != nil {
			s.close()
			return
		}

		s.Conn.WriteMessage(websocket.TextMessage, buf[:n])
	}
}

func (s *ShellSession) bridgeReverse() {

	go func() {
		for {
			_, msg, err := s.Conn.ReadMessage()
			if err != nil {
				s.close()
				return
			}

			s.Net.Write(msg)
		}
	}()

	buf := make([]byte, 1024)

	for {
		n, err := s.Net.Read(buf)
		if err != nil {
			s.close()
			return
		}

		err = s.Conn.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			s.close()
			return
		}
	}
}

func (s *ShellSession) close() {
	if s.Pty != nil {
		s.Pty.Close()
	}

	if s.Cmd != nil && s.Cmd.Process != nil {
		s.Cmd.Process.Kill()
	}

	if s.Net != nil {
		s.Net.Close()
	}

	if s.Conn != nil {
		s.Conn.Close()
	}

	manager.DeleteSession(s.ID)
}