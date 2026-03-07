package main

import (
	"openshellserver/internal/listener"
	"openshellserver/internal/session"
)

func main() {
	manager := session.NewManager()
	listener.StartListener("4444", manager)
}