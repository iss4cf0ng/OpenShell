package listener

import (
	"fmt"
	"net"

	"openshellserver/internal/session"
)

func StartListener(port string, manager *session.Manager) {

	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		panic(err)
	}

	fmt.Println("Listening on port", port)

	for {

		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		s := manager.AddSession(conn)

		fmt.Println("New session:", s.ID, s.RemoteIP)

		go handleConnection(s, manager)

	}

}

func handleConnection(s *session.Session, manager *session.Manager) {

	defer func() {
		manager.RemoveSession(s.ID)
		s.Conn.Close()
		fmt.Println("Session closed:", s.ID)
	}()

	buf := make([]byte, 1024)

	for {

		_, err := s.Conn.Read(buf)

		if err != nil {
			break
		}

	}

}