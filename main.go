package main

// import "net"
// import "log"
// import "io"
import "fmt"
// import "strconv"

// main function 
func command(k *Kademlia, testing bool) { 
	fmt.Println("cmd")
	if (testing) {
		for {}
	}
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
			k.nt.SendJoinMessage(ip)
		} else if (cmd == "-ping") {
			k.nt.SendPingAll();
		} else if (cmd == "-lookup") {
			fmt.Println("Enter id: ") 
			var id string 
			fmt.Scanln(&id) 
			fmt.Println("Enter ip: ") 
			var ip string 
			fmt.Scanln(&ip)
			var contact = NewContact(NewKademliaID(id),ip)
			k.LookupContact(&contact, &k.nt.rt.me)
		} else if (cmd == "-put") {
			fmt.Println("Enter content to store: ") 
			var data string 
			fmt.Scanln(&data) 
			hash := NewRandomKademliaID()
			k.Store(hash.String(),([]byte(data)))
		} else if (cmd == "-get") {
			fmt.Println("Enter the hash: ") 
			var hash string 
			fmt.Scanln(&hash)
			k.LookupData(hash) 
		} else if (cmd == "-exit") {
			k.nt.terminate = true //Stop listen loop
			k.nt.kademliaChannel <- ([]byte("+TERMINATE+")) //Stop datahandler
			break
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
	// fmt.Println("Enter port: ") 
  
	// var port string 
	
	// fmt.Scanln(&port) 
	// fmt.Println("jada")
	rt := NewRoutingTable(NewContact(NewRandomKademliaID(), GetIpAddress()+"8000"))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)

	ka := NewKademlia(nt)

	// iport, err := strconv.Atoi(port)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go nt.Listen(8000)
	go ka.DataHandler()
	addressToJoin := ka.nt.GetIpToJoin()
	if (addressToJoin != "") {
		go ka.nt.SendJoinMessage(addressToJoin+":8000")
	}
	command(ka,true)
}