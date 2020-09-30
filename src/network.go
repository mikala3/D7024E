package main

import (
	//"os/exec"
	"strings"
	"net"
	"log"
	"strconv"
	"bytes"
	"fmt"
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
		fmt.Println(newstring)
		split := strings.Split(newstring,">")
		stringarr := strings.Split(split[0][8:(len(split[0])-1)],",")
		stringarr2 := strings.Split(split[1][8:(len(split[1])-1)],",")
		id := stringarr[0]
		address := strings.Split(stringarr[1][1:], ")")
		id2 := stringarr2[0]
		address2 := strings.Split(stringarr2[1][1:], ")")
		contact := NewContact(NewKademliaID(id),address[0])
		contact2 := NewContact(NewKademliaID(id2),address2[0])
		network.SendPingMessage(&contact, &contact2)
		//newstring = contact(id, address)
	} else if bytes.Contains(b, []byte("PingAccepted<")) {
		var newdata []byte = b[13:]
		newstring := string(newdata)
		fmt.Println(newstring)
		split := strings.Split(newstring, ">")
		stringarr := strings.Split(split[0][8:(len(split[0])-1)], ",")
		stringarr2 := strings.Split(split[1][8:(len(split[1])-1)], ",")
		id := stringarr[0]
		address := strings.Split(stringarr[1][1:], ")")
		id2 := stringarr2[0]
		address2 := strings.Split(stringarr2[1][1:], ")")
		contact := NewContact(NewKademliaID(id),address[0])
		contact2 := NewContact(NewKademliaID(id2),address2[0])
		network.SendPingAccepted(&contact, &contact2)
		//newstring = contact(id, address)
	} else if bytes.Contains(b, []byte("Find<")) {
		var newdata []byte = b[5:]
		newstring := string(newdata)
		fmt.Println(newstring)
		split := strings.Split(newstring, ">")
		stringarr := strings.Split(split[0][8:(len(split[0])-1)], ",")
		stringarr2 := strings.Split(split[1][8:(len(split[1])-1)], ",")
		id := stringarr[0]
		address := strings.Split(stringarr[1][1:], ")")
		id2 := stringarr2[0]
		address2 := strings.Split(stringarr2[1][1:], ")")
		contact := NewContact(NewKademliaID(id),address[0])
		contact2 := NewContact(NewKademliaID(id2),address2[0])
		network.SendFindContactMessage(&contact, &contact2)
	} else if bytes.Contains(b, []byte("FindAccepted<")) {
		var newdata []byte = b[13:]
		newstring := string(newdata)
		fmt.Println(newstring)
		split := strings.Split(newstring, ">")
		stringarr := strings.Split(split[0][8:(len(split[0])-1)], ",")
		stringarr2 := strings.Split(split[1][8:(len(split[1])-1)], ",")
		id := stringarr[0]
		address := strings.Split(stringarr[1][1:], ")")
		id2 := stringarr2[0]
		address2 := strings.Split(stringarr2[1][1:], ")")
		contact := NewContact(NewKademliaID(id),address[0])
		contact2 := NewContact(NewKademliaID(id2),address2[0])
		network.SendFindAccepted(&contact, &contact2)
	} else if bytes.Contains(b, []byte("FindData<")) {
		//var newdata []byte = b[9:]
		//newstring := string(newdata)
	} else if bytes.Contains(b, []byte("Data<")) {
		//var newdata []byte = b[5:]
		//newstring := string(newdata)
	} else if bytes.Contains(b, []byte("Join<")) {
		var newdata []byte = b[5:]
		newstring := string(newdata)
		fmt.Println(newstring)
		stringarr := strings.Split(newstring[8:(len(newstring)-1)], ",")
		id := stringarr[0]
		address := strings.Split(stringarr[1][1:], ")")
		network.rt.AddContact(NewContact(NewKademliaID(id), address[0]))
		network.SendJoinAcceptedMessage(address[0])
	} else if bytes.Contains(b, []byte("JoinAccepted<")) {
		var newdata []byte = b[13:]
		newstring := string(newdata)
		fmt.Println(newstring)
		stringarr := strings.Split(newstring[8:(len(newstring)-1)], ",")
		id := stringarr[0]
		address := strings.Split(stringarr[1][1:], ")")
		network.rt.AddContact(NewContact(NewKademliaID(id), address[0]))

	} else {
		fmt.Println("Something incorect with incoming message!")
	}
}
//fmt.Println(strings.TrimSuffix("localhost:8080)", ")"))
func (network *Network) SendPingAccepted(contact *Contact, sender *Contact) {
	if network.rt.me.ID.Equals(contact.ID) {
		//Ping accepted recivied, ping bounced back.
		fmt.Println("Ping bounced back from "+sender.String())
	} else {
		coid := network.rt.FindClosestContacts(contact.ID, 3)
		for co := 0; co < len(coid); co++ {
			conn, err := net.Dial("tcp", coid[co].Address)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			if _, err := conn.Write([]byte("PingAccepted<"+contact.String()+">"+sender.String())); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact, sender *Contact) {
	if network.rt.me.ID.Equals(contact.ID) {
		fmt.Println("Ping recivied from "+sender.String())
		network.SendPingAccepted(sender, &network.rt.me)
	} else {
		coid := network.rt.FindClosestContacts(contact.ID, 3)
		for co := 0; co < len(coid); co++ {
			conn, err := net.Dial("tcp", coid[co].Address)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			if _, err := conn.Write([]byte("Ping<"+contact.String()+">"+sender.String())); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (network *Network) SendPingAll() {
	// execOut, _ := exec.Command("ping",contact.Address,"-c 3","-w 10").Output()
	// if strings.Contains(string(execOut), "Destination Host Unreachable") {
	// 	log.Fatal("Destination Host Unreachable")
	// }

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
	}
}

func (network *Network) SendFindAccepted(contact *Contact, sender *Contact) {
	if network.rt.me.ID.Equals(contact.ID) {
		//Ping accepted recivied, ping bounced back.
		fmt.Println("Find bounced back from "+sender.String())
		network.rt.AddContact(*sender)
	} else {
		coid := network.rt.FindClosestContacts(contact.ID, 3)
		for co := 0; co < len(coid); co++ {
			conn, err := net.Dial("tcp", coid[co].Address)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			if _, err := conn.Write([]byte("FindAccepted<"+contact.String()+">"+sender.String())); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (network *Network) SendFindContactMessage(contact *Contact, sender *Contact) {
	if network.rt.me.ID.Equals(contact.ID) {
		fmt.Println("Lookup recivied from "+sender.String())
		network.SendFindAccepted(sender, &network.rt.me)
	}
	coid := network.rt.FindClosestContacts(contact.ID, 3)
	for co := 0; co < len(coid); co++ {
		conn, err := net.Dial("tcp", coid[co].Address)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte("Find<"+contact.String()+">"+network.rt.me.String())); err != nil {
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