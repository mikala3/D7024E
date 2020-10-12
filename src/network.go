package main

import (
	//"os/exec"
	"net"
	"log"
	"strconv"
	"fmt"
)

type Network struct {
	rt *RoutingTable
	kademliaChannel chan []byte
	externalChannel chan []byte
	testing bool
}

// NewNetwork returns a new instance of a RNetwork
func NewNetwork(rt *RoutingTable, kc chan []byte, ex chan []byte) *Network {
	network := &Network{}
	network.rt = rt
	network.kademliaChannel = kc
	network.externalChannel = ex
	network.testing = false
	return network
}

func (network *Network) ListenToIp(ip string, port int) {
	listenip := ip + ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", listenip)
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
			data := make([]byte, 1024)
			_, err := c.Read(data) //Read data sent
			if err != nil {
				panic(err)
			}
			//msg := <- network.externalChannel
			network.kademliaChannel <- data
			c.Close()
		}(conn)
	}
}

func (network *Network) Listen(port int) {
	fmt.Println("listen")
	listenip := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", listenip)
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
			data := make([]byte, 1024)
			_, err := c.Read(data) //Read data sent
			if err != nil {
				panic(err)
			}
			//msg := <- network.externalChannel
			network.kademliaChannel <- data //Handle data
			c.Close()
		}(conn)
	}
}


//fmt.Println(strings.TrimSuffix("localhost:8080)", ")"))
func (network *Network) SendPingAccepted(contact *Contact, sender *Contact) {
	if (network.testing) {
		if network.rt.me.ID.Equals(contact.ID) {
			//Ping accepted recivied, ping bounced back.
			fmt.Println("Ping bounced back from "+sender.String())
		} else {
			network.externalChannel <- ([]byte("PingAccepted<"+contact.String()+">"+sender.String()))
		}
	} else {
		if network.rt.me.ID.Equals(contact.ID) {
			//Ping accepted recivied, ping bounced back.
			fmt.Println("Ping bounced back from "+sender.String())
		} else {
			coid := network.rt.FindClosestContacts(contact.ID, 1)
			for co := 0; co < len(coid); co++ {
				conn, err := net.Dial("tcp", coid[co].Address)
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()
	
				if _, err := conn.Write([]byte("PingAccepted<"+contact.String()+">"+sender.String())); err != nil {
					log.Fatal(err)
				}
				// network.externalChannel <- ([]byte("PingAccepted<"+contact.String()+">"+sender.String()))
				// conn.Close()
			}
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact, sender *Contact) {
	if (network.testing) {
		if network.rt.me.ID.Equals(contact.ID) {
			fmt.Println("Ping recivied from "+sender.String())
			network.SendPingAccepted(sender, &network.rt.me)
		} else {
			network.externalChannel <- ([]byte("Ping<"+contact.String()+">"+sender.String()))
		}
	} else {
		if network.rt.me.ID.Equals(contact.ID) {
			fmt.Println("Ping recivied from "+sender.String())
			network.SendPingAccepted(sender, &network.rt.me)
		} else {
			coid := network.rt.FindClosestContacts(contact.ID, 1)
			for co := 0; co < len(coid); co++ {
				conn, err := net.Dial("tcp", coid[co].Address)
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				if _, err := conn.Write([]byte("Ping<"+contact.String()+">"+sender.String())); err != nil {
					log.Fatal(err)
				}
				// network.externalChannel <- ([]byte("Ping<"+contact.String()+">"+sender.String()))
				// conn.Close()
			}
		}
	}
}

func (network *Network) SendPingAll() {
	// execOut, _ := exec.Command("ping",contact.Address,"-c 3","-w 10").Output()
	// if strings.Contains(string(execOut), "Destination Host Unreachable") {
	// 	log.Fatal("Destination Host Unreachable")
	// }
	if (network.testing) {
		coid := network.rt.FindClosestContacts(network.rt.me.ID, 3)
		for co := 0; co < len(coid); co++ {
			network.externalChannel <- ([]byte("Ping<"+coid[co].String()+">"+network.rt.me.String()))
		}
	} else {
		coid := network.rt.FindClosestContacts(network.rt.me.ID, 3)
		for co := 0; co < len(coid); co++ {
			conn, err := net.Dial("tcp", coid[co].Address)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			if _, err := conn.Write([]byte("Ping<"+coid[co].String()+">"+network.rt.me.String())); err != nil {
				log.Fatal(err)
			}
			// network.externalChannel <- ([]byte("Ping<"+coid[co].String()+">"+network.rt.me.String()))
			// conn.Close()
		}
	}
}

func (network *Network) SendFindAccepted(contact *Contact, sender *Contact) {
	if (network.testing) {
		coid := network.rt.FindClosestContacts(contact.ID, alpha)
		senderstring := ""
		for co := 0; co < len(coid); co++ {
			senderstring = senderstring + coid[co].String() +">"
		}
		network.externalChannel <- ([]byte("FindAccepted<"+sender.String()+">"+contact.String()+">"+senderstring[: (len(senderstring)-1) ]))
	} else {
		coid := network.rt.FindClosestContacts(contact.ID, alpha)
		conn, err := net.Dial("tcp", sender.Address)
		senderstring := ""
		for co := 0; co < len(coid); co++ {
			senderstring = senderstring + coid[co].String() +">"
		}
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("FindAccepted<"+sender.String()+">"+contact.String()+">"+senderstring[: (len(senderstring)-1) ])); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("FindAccepted<"+sender.String()+">"+contact.String()+">"+senderstring[: (len(senderstring)-1) ]))
		// conn.Close()
	}
}

func (network *Network) SendFindContactMessage(contact *Contact, recivier *Contact, sender *Contact) {
	if (network.testing) {
		network.externalChannel <- ([]byte("Find<"+contact.String()+">"+sender.String()))
	} else {
		conn, err := net.Dial("tcp", recivier.Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Find<"+contact.String()+">"+sender.String())); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("Find<"+contact.String()+">"+sender.String()))
		// conn.Close()
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

//Needs more work later, need to match with data to find if we are close enough with hash
func (network *Network) SendStoreMessage(data []byte, contact *Contact) {
	if (network.testing) {
		if network.rt.me.ID.Equals(contact.ID) {
			dataString := string(data)
			network.externalChannel <- ([]byte("Store<"+contact.String()+dataString))
		} else {
			network.externalChannel <- ([]byte("Find<"+contact.String()))
		}
	} else {
		// TODO
		if network.rt.me.ID.Equals(contact.ID) {
			//We are the node that is trying to be found
			/* k closest nodes to the target node */
			k_closest_nodes := network.rt.FindClosestContacts(contact.ID, 20)
			for co := 0; co < len(k_closest_nodes); co++ {
				conn, err := net.Dial("tcp", k_closest_nodes[co].Address)
				dataString := string(data)
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()
		
				if _, err := conn.Write([]byte("Store<"+contact.String()+dataString)); err != nil {
					log.Fatal(err)
				}
				// network.externalChannel <- ([]byte("Store<"+contact.String()+dataString))
				// conn.Close()
			} 
		} else {
			coid := network.rt.FindClosestContacts(contact.ID, 20)
			for co := 0; co < len(coid); co++ {
				conn, err := net.Dial("tcp", coid[co].Address)
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()
		
				if _, err := conn.Write([]byte("Find<"+contact.String())); err != nil {
					log.Fatal(err)
				}
				// network.externalChannel <- ([]byte("Find<"+contact.String()))
				// conn.Close()
			} 
		}
	}

}

func (network *Network) SendJoinMessage(ip string) {
	if (network.testing) {
		network.externalChannel <- ([]byte("Join<"+network.rt.me.String()))
	} else {
		conn, err := net.Dial("tcp", ip)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Join<"+network.rt.me.String())); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("Join<"+network.rt.me.String()))
		// conn.Close()
	}
}

func (network *Network) SendJoinAcceptedMessage(ip string) {
	if (network.testing) {
		network.externalChannel <- ([]byte("JoinAccepted<"+network.rt.me.String()))
	} else {
		conn, err := net.Dial("tcp", ip)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("JoinAccepted<"+network.rt.me.String())); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("JoinAccepted<"+network.rt.me.String()))
		// conn.Close()
	}
}

// func (network *Network) SendFindAccepted(contact *Contact, sender *Contact) {
// 	if network.rt.me.ID.Equals(contact.ID) {
// 		//Ping accepted recivied, ping bounced back.
// 		fmt.Println("Find bounced back from "+sender.String())
// 		network.rt.AddContact(*sender)
// 	} else {
// 		coid := network.rt.FindClosestContacts(contact.ID, 1)
// 		for co := 0; co < len(coid); co++ {
// 			conn, err := net.Dial("tcp", coid[co].Address)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer conn.Close()

// 			if _, err := conn.Write([]byte("FindAccepted<"+contact.String()+">"+sender.String())); err != nil {
// 				log.Fatal(err)
// 			}
// 			// network.externalChannel <- ([]byte("FindAccepted<"+contact.String()+">"+sender.String()))
// 			// conn.Close()
// 		}
// 	}
// }

// func (network *Network) SendFindContactMessage(contact *Contact, sender *Contact) {
// 	if network.rt.me.ID.Equals(contact.ID) {
// 		fmt.Println("Lookup recivied from "+sender.String())
// 		network.SendFindAccepted(sender, &network.rt.me)
// 	} else {
// 		coid := network.rt.FindClosestContacts(contact.ID, 1)
// 		for co := 0; co < len(coid); co++ {
// 			conn, err := net.Dial("tcp", coid[co].Address)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer conn.Close()

// 			if _, err := conn.Write([]byte("Find<"+contact.String()+">"+sender.String())); err != nil {
// 				log.Fatal(err)
// 			}
// 			// network.externalChannel <- ([]byte("Find<"+contact.String()+">"+sender.String()))
// 			// conn.Close()
// 		}
// 	}
// }