package main

import "net"
import "os"
import "fmt"
import "log"

func masternode() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Write host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	go masterconnect(c)
}

func masterconnect(c net.Conn) {
	defer c.Close()

	if _, err := c.Write([]byte("Master node message")); err != nil {
		log.Fatal(err)
	}
}