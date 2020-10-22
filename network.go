package main

import (
	//"os/exec"
	"net"
	"log"
	"strconv"
	"fmt"
	"os/exec"
)

type Network struct {
	rt *RoutingTable
	kademliaChannel chan []byte
	externalChannel chan []byte
	testing bool
	terminate bool
}

// NewNetwork returns a new instance of a RNetwork
func NewNetwork(rt *RoutingTable, kc chan []byte, ex chan []byte) *Network {
	network := &Network{}
	network.rt = rt
	network.kademliaChannel = kc
	network.externalChannel = ex
	network.testing = false
	network.terminate = false
	return network
}

func GetIpAddress() string {
	command, err := exec.Command("awk 'END{print $1}' /etc/hosts").Output()
	if (err != nil) {
		fmt.Println(err)
		return ""
	}
	return string(command)
}

func (network *Network) GetIpToJoin() string {
	command, err := exec.Command("ping swarm_kademliaNodes.1 | jq '.ip'").Output()
	if (err != nil) {
		fmt.Println(err)
		return ""
	}
	if (string(command) == network.rt.me.Address) {
		return ""
	}
	return string(command)
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
	//fmt.Println("listen")
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
		network.externalChannel <- ([]byte("PingAccepted<"+contact.String()+">"+network.rt.me.String()))
	} else {
		conn, err := net.Dial("tcp", contact.Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("PingAccepted<"+contact.String()+">"+network.rt.me.String())); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("PingAccepted<"+contact.String()+">"+sender.String()))
		// conn.Close()
	}
}

func (network *Network) SendPingMessage(contact *Contact, sender *Contact) {
	if (network.testing) {
		network.externalChannel <- ([]byte("Ping<"+contact.String()+">"+network.rt.me.String()))
	} else {
		conn, err := net.Dial("tcp", contact.Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Ping<"+contact.String()+">"+network.rt.me.String())); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("Ping<"+contact.String()+">"+sender.String()))
		// conn.Close()
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

// NEW FORMAT FOR FIND, OLD WAS REEEEEAAAAALLLYYYY BAD
// Find<contact(SEARCHED)>contact(SENDER)
// FindAccepted<contact(SEARCHED)>contact(SENDER)>contact(BUCKET1)>contact(BUCKET2) ...
// Recivier is no longer included since they are directly connected, the sender is still sent by function call to network.

func (network *Network) SendFindAccepted(contact *Contact, sender *Contact) {
	if (network.testing) {
		coid := network.rt.FindClosestContacts(contact.ID, alpha)
		senderstring := ""
		for co := 0; co < len(coid); co++ {
			senderstring = senderstring + coid[co].String() +">"
		}
		network.externalChannel <- ([]byte("FindAccepted<"+contact.String()+">"+network.rt.me.String()+">"+senderstring[: (len(senderstring)-1) ]))
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

		if _, err := conn.Write([]byte("FindAccepted<"+contact.String()+">"+network.rt.me.String()+">"+senderstring[: (len(senderstring)-1) ])); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("FindAccepted<"+sender.String()+">"+contact.String()+">"+senderstring[: (len(senderstring)-1) ]))
		// conn.Close()
	}
}

func (network *Network) SendFindContactMessage(contact *Contact, recivier *Contact, sender *Contact) {
	if (network.testing) {
		network.externalChannel <- ([]byte("Find<"+contact.String()+">"+network.rt.me.String()))
	} else {
		conn, err := net.Dial("tcp", recivier.Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Find<"+contact.String()+">"+network.rt.me.String())); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("Find<"+contact.String()+">"+sender.String()))
		// conn.Close()
	}
}

func (network *Network) SendLookupDataMessage(contact *Contact, sender *Contact, hash string) {
	if (network.testing) {
		network.externalChannel <- ([]byte("FindData<"+contact.String()+">"+sender.String()+">"+hash))
	} else {
		conn, err := net.Dial("tcp", contact.Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("FindData<"+contact.String()+">"+sender.String()+">"+hash)); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("Store<"+contact.String()+dataString))
		// conn.Close()
	} 
}

func (network *Network) SendFoundDataMessage(contact *Contact, hash string, data []byte) {
	if (network.testing) {
		dataString := string(data)
		network.externalChannel <- ([]byte("FoundData<"+contact.String()+">"+network.rt.me.String()+">"+hash+">"+dataString))
	} else {
		conn, err := net.Dial("tcp", contact.Address)
		dataString := string(data)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("FoundData<"+contact.String()+">"+network.rt.me.String()+">"+hash+">"+dataString)); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("Store<"+contact.String()+dataString))
		// conn.Close()
	} 
}

//Needs more work later, need to match with data to find if we are close enough with hash
func (network *Network) SendStoreDataMessage(contact *Contact, hash string, data []byte) {
	if (network.testing) {
		dataString := string(data)
		network.externalChannel <- ([]byte("Data<"+network.rt.me.String()+">"+hash+">"+dataString))
	} else {
		conn, err := net.Dial("tcp", contact.Address)
		dataString := string(data)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Data<"+network.rt.me.String()+">"+hash+">"+dataString)); err != nil {
			log.Fatal(err)
		}
		// network.externalChannel <- ([]byte("Store<"+contact.String()+dataString))
		// conn.Close()
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