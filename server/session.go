package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"openshell/internal/logger"

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

	mu     sync.Mutex
	OutBuf chan []byte //Cache
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

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*ShellSession),
	}
}

func (sm *SessionManager) nextID() string {
	sm.counter++
	return fmt.Sprintf("session-%d", sm.counter)
}

func (sm *SessionManager) Count() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return len(sm.sessions)
}

func (sm *SessionManager) ListSessions() []Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	list := []Session{}
	for _, s := range sm.sessions {
		list = append(list, Session{
			ID:   s.ID,
			Type: s.Type,
			IP:   s.RemoteAddr,
		})
	}
	return list
}

func (sm *SessionManager) DeleteSession(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, id)
	logger.Info("Session removed: %s", id)
}

//PTY
func (sm *SessionManager) CreateLocalShell(conn *websocket.Conn) (*ShellSession, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := sm.nextID()
	cmd := exec.Command("/bin/bash", "-i")
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
		OutBuf: make(chan []byte, 1024*10), // buffer
	}

	sm.sessions[id] = session

	go session.bridgePTY()

	ptyFile.Write([]byte("\n"))

	return session, nil
}

func (s *ShellSession) bridgePTY() {
	buf := make([]byte, 1024)
	for {
		n, err := s.Pty.Read(buf)
		if err != nil {
			s.closePTYOnly()
			return
		}
		s.mu.Lock()
		if s.Conn != nil {
			s.Conn.WriteMessage(websocket.TextMessage, buf[:n])
		} else {
			select {
			case s.OutBuf <- buf[:n]:
			default:
			}
		}
		s.mu.Unlock()
	}
}

func (s *ShellSession) attachWebSocket(conn *websocket.Conn) {
    s.mu.Lock()
    s.Conn = conn
    s.mu.Unlock()

    go func() {
        for {
            select {
            case data, ok := <-s.OutBuf:
                if !ok {
                    return
                }
                conn.WriteMessage(websocket.TextMessage, data)
            default:
                return
            }
        }
    }()

    go func() {
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                s.mu.Lock()
                s.Conn = nil
                s.mu.Unlock()
                return
            }
            s.Pty.Write(msg)
        }
    }()
}

func (s *ShellSession) closePTYOnly() {
	if s.Pty != nil {
		s.Pty.Close()
	}
	if s.Cmd != nil && s.Cmd.Process != nil {
		s.Cmd.Process.Kill()
	}
	manager.DeleteSession(s.ID)
}

//Reverse shell
func (sm *SessionManager) CreateReverseShell(netConn net.Conn) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := sm.nextID()
	session := &ShellSession{
		ID:         id,
		Type:       "reverse",
		Net:        netConn,
		RemoteAddr: netConn.RemoteAddr().String(),
		OutBuf:     make(chan []byte, 1024*10),
	}

	sm.sessions[id] = session

	go session.bridgeReverse()

	logger.Success("New reverse shell: %s %s", id, netConn.RemoteAddr())
}

func (s *ShellSession) bridgeReverse() {
	buf := make([]byte, 1024)
	go func() {
		for {
			if s.Conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			_, msg, err := s.Conn.ReadMessage()
			if err != nil {
				s.mu.Lock()
				s.Conn = nil
				s.mu.Unlock()
				return
			}
			s.Net.Write(msg)
		}
	}()

	for {
		n, err := s.Net.Read(buf)
		if err != nil {
			s.Net.Close()
			return
		}
		s.mu.Lock()
		if s.Conn != nil {
			s.Conn.WriteMessage(websocket.TextMessage, buf[:n])
		} else {
			select {
			case s.OutBuf <- buf[:n]:
			default:
			}
		}
		s.mu.Unlock()
	}
}

func (s *ShellSession) attachReverse(conn *websocket.Conn) {
	s.attachWebSocket(conn)
}

func (s *ShellSession) closeNetOnly() {
	if s.Net != nil {
		s.Net.Close()
	}
	manager.DeleteSession(s.ID)
}