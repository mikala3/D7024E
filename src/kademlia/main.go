package kademlia

// import "net"
// import "log"
// import "io"
// import "fmt"

func main() {
	// fmt.Println("Starting node")
	// ln, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer ln.Close()
	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	go func(c net.Conn) { //Connection handler
	// 		io.Copy(c, c)
	// 		c.Close()
	// 	}(conn)
	// }
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	nt := NewNetwork(rt)

	go nt.Listen(8080)
}

func dumbnode() {
	go conn()
}

func conn() {
	
}