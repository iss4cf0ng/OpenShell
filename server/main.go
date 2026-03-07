package main

import (
    "fmt"
	"log"
	"net/http"
    "encoding/json"
    "net"
    "crypto/tls"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var manager = NewSessionManager()

func SessionHandler(w http.ResponseWriter, r *http.Request) {
    list := manager.ListSessions()
    json.NewEncoder(w).Encode(list)
}

func AttachHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    session := manager.sessions[id]
    if session == nil {
        return
    }

    session.Conn = conn

    if session.Net != nil {
        go session.bridgeReverse()
    }

    if session.Pty != nil {
        go session.bridgePTY()
    }
}

func StartTLSReverseShell(port string) {
    cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := tls.Listen("tcp", ":"+port, config)
	if err != nil {
		panic(err)
	}

	fmt.Println("TLS reverse shell listener on", port)

	for {

		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go func(c net.Conn) {

			manager.CreateReverseShell(c)

		}(conn)
	}
}

func main() {

	// reverse shell listener
	go StartReverseShellListener("4444") //Normal TCP
    go StartTLSReverseShell("4445") //TLS

    http.HandleFunc("/api/sessions", SessionHandler)
    http.HandleFunc("/ws/session", AttachHandler)

	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)

	log.Println("OpenShellServer running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}