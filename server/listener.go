package main

import (
	"net"
	"crypto/tls"

	"openshell/internal/logger"
)

//TCP listener
func StartReverseShell(port string) {

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	logger.Info("[*] Reverse shell listener on port: %s", port)

	for {

		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go manager.CreateReverseShell(conn)
	}
}

//TLS listener
func StartTLSReverseShell(port string) {
    cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := tls.Listen("tcp", ":" + port, config)
	if err != nil {
		panic(err)
	}

	logger.Info("[*] TLS reverse shell listener on: %s", port)

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