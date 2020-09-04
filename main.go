package main

import "net"
import "log"
import "io"
import "fmt"

func main() {
	fmt.Println("Starting node")
	go dumbnode()
}

func dumbnode() {
	go conn()
}

func conn() {
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
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}