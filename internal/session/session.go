package session

import (
	"net"
	"time"
)

type Session struct {
	ID string
	Conn net.Conn
	RemoteIP string
	CreatedAt time.Time
}