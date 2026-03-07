package main

import (
	"log"
	"net/http"
    "encoding/json"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var manager = NewSessionManager()

/*
func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	session := manager.GetLastReverseShell()

	if session != nil {

		session.Conn = conn

		go session.bridgeReverse()

		log.Println("Web attached to reverse shell:", session.ID)

		return
	}

	localSession, err := manager.CreateLocalShell(conn)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("New local shell:", localSession.ID)
}
*/

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

func main() {

	// reverse shell listener
	go StartReverseShellListener("4444")

	//http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/api/sessions", SessionHandler)
    http.HandleFunc("/ws/session", AttachHandler)

	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)

	log.Println("OpenShellServer running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}