package kademlia

import "net"
import "log"

func smartnode() {
	go smartconn()
}

func smartconn() {
	ln, err := net.Listen("tcp", "8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) { //Connection handler
			go smartdialer()
			c.Close()
		}(conn)
	}
}

func smartdialer() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("Hello, World!")); err != nil {
		log.Fatal(err)
	}
}