package main

// import "net"
import "log"
// import "io"
import "fmt"
import "strconv"

// main function 
func command(n *Network) { 
	fmt.Println("cmd")
	for {
		fmt.Println("Enter command (-h for help): ") 
  
		var cmd string 
	
		fmt.Scanln(&cmd) 
		if (cmd == "-h") {
			fmt.Println("Commands \n-h Help\n-join Join network\n-ping Ping\n-put Upload object\n-get Get object\n-exit Terminate node")
		} else if (cmd == "-join") {
			fmt.Println("Enter ip: ") 
			var ip string 
			fmt.Scanln(&ip) 
			n.SendJoinMessage(ip)
		} else if (cmd == "-ping") {
			n.SendPingAll();
		} else if (cmd == "-lookup") {
			fmt.Println("Enter id: ") 
			var id string 
			fmt.Scanln(&id) 
			fmt.Println("Enter ip: ") 
			var ip string 
			fmt.Scanln(&ip)
			var contact = NewContact(NewKademliaID(id),ip)
			n.SendFindContactMessage(&contact, &n.rt.me)
		} else if (cmd == "-put") {
			fmt.Println("Not implemented ") 
		} else if (cmd == "-get") {
			fmt.Println("Not implemented ") 
		} else if (cmd == "-exit") {
			fmt.Println("Not implemented ") 
		}
	}
} 

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
	fmt.Println("Enter port: ") 
  
	var port string 
	
	fmt.Scanln(&port) 
	fmt.Println("jada")
	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)

	ka := NewKademlia(nt)

	iport, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	go nt.Listen(iport)
	go ka.DataHandler()
	command(nt)
}