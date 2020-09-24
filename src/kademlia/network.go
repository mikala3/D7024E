package kademlia

import (
	//"os/exec"
	"strings"
	"net"
	"log"
	"strconv"
	"bytes"
)

type Network struct {
	rt *RoutingTable
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewNetwork(rt *RoutingTable) *Network {
	network := &Network{}
	network.rt = rt
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
			data := make([]byte, 128)
			_, err := c.Read(data) //Read data sent
			if err != nil {
				panic(err)
			}
			go network.ListenDataHandler(data) //Handle data
			c.Close()
		}(conn)
	}
}

func (network *Network) Listen(port int) {
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
			data := make([]byte, 128)
			_, err := c.Read(data) //Read data sent
			if err != nil {
				panic(err)
			}
			go network.ListenDataHandler(data) //Handle data
			c.Close()
		}(conn)
	}
}

func (network *Network) ListenDataHandler(b []byte) {
	if bytes.Contains(b, []byte("Ping<")) {
		var newdata []byte = b[5:]
		newstring := string(newdata)
		stringarr := strings.Split(",",newstring[8:(len(newstring)-1)])
		id := stringarr[0]
		address := stringarr[1]
		contact := NewContact(NewKademliaID(id),address[1:])
		network.SendPingMessage(&contact)
		//newstring = contact(id, address)
	} else if bytes.Contains(b, []byte("Find<")) {
		var newdata []byte = b[5:]
		newstring := string(newdata)
		stringarr := strings.Split(",",newstring[8:(len(newstring)-1)])
		id := stringarr[0]
		address := stringarr[1]
		contact := NewContact(NewKademliaID(id),address[1:])
		network.SendFindContactMessage(&contact)
	} else if bytes.Contains(b, []byte("FindData<")) {
		//var newdata []byte = b[9:]
		//newstring := string(newdata)
	} else if bytes.Contains(b, []byte("Data<")) {
		//var newdata []byte = b[5:]
		//newstring := string(newdata)
	} else if bytes.Contains(b, []byte("Join<")) {
		var newdata []byte = b[5:]
		newstring := string(newdata)
		stringarr := strings.Split(",",newstring[8:(len(newstring)-1)])
		id := stringarr[0]
		address := stringarr[1]
		network.rt.AddContact(NewContact(NewKademliaID(id), address))
		network.SendJoinAcceptedMessage(address)
	} else if bytes.Contains(b, []byte("JoinAccepted<")) {
		var newdata []byte = b[13:]
		newstring := string(newdata)
		stringarr := strings.Split(",",newstring[8:(len(newstring)-1)])
		id := stringarr[0]
		address := stringarr[1]
		network.rt.AddContact(NewContact(NewKademliaID(id), address))

	}
}

//Needs to be changed to fit model
func (network *Network) SendPingMessage(contact *Contact) {
	// execOut, _ := exec.Command("ping",contact.Address,"-c 3","-w 10").Output()
	// if strings.Contains(string(execOut), "Destination Host Unreachable") {
	// 	log.Fatal("Destination Host Unreachable")
	// }

	if network.rt.me.ID == contact.ID {
		//We are the node that is trying to be pinged
	}
	coid := network.rt.FindClosestContacts(contact.ID, 20)
	for co := 0; co < 3; co++ {
		conn, err := net.Dial("tcp", coid[co].Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Ping<"+contact.String())); err != nil {
			log.Fatal(err)
		}
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	if network.rt.me.ID == contact.ID {
		//We are the node that is trying to be found

	}
	coid := network.rt.FindClosestContacts(contact.ID, 20)
	for co := 0; co < 3; co++ {
		conn, err := net.Dial("tcp", coid[co].Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Find<"+contact.String())); err != nil {
			log.Fatal(err)
		}
	}


}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

//Needs more work later, need to match with data to find if we are close enough with hash
func (network *Network) SendStoreMessage(data []byte, contact *Contact) {
	// TODO
	if network.rt.me.ID == contact.ID {
		//We are the node that is trying to be found
		/* k closest nodes to the target node */
		k_closest_nodes := network.rt.FindClosestContacts(contact.ID, 20)
		for co := 0; co < 20; co++ {
			conn, err := net.Dial("tcp", k_closest_nodes[co].Address)
			dataString := string(data)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
	
			if _, err := conn.Write([]byte("Store<"+contact.String()+dataString)); err != nil {
				log.Fatal(err)
			}
		} 
	} else {
		coid := network.rt.FindClosestContacts(contact.ID, 20)
		for co := 0; co < 3; co++ {
			conn, err := net.Dial("tcp", coid[co].Address)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
	
			if _, err := conn.Write([]byte("Find<"+contact.String())); err != nil {
				log.Fatal(err)
			}
		} 
	}

}

func (network *Network) SendJoinMessage(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("Join<"+network.rt.me.String())); err != nil {
		log.Fatal(err)
	}
}

func (network *Network) SendJoinAcceptedMessage(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("JoinAccepted<"+network.rt.me.String())); err != nil {
		log.Fatal(err)
	}
}