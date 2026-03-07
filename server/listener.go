package main

import (
	"fmt"
	"net"
)

func StartReverseShellListener(port string) {

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reverse shell listener on port", port)

	for {

		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go manager.CreateReverseShell(conn)
	}
}