package main


import (
	//"os/exec"
	"strings"
	"bytes"
	"fmt"
)

func (kademlia *Kademlia) DataHandler() {
	for {
		b := <- kademlia.nt.kademliaChannel
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
			kademlia.nt.SendPingMessage(&contact, &contact2)
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
			kademlia.nt.SendPingAccepted(&contact, &contact2)
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
			kademlia.nt.SendFindContactMessage(&contact, &contact2)
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
			kademlia.nt.SendFindAccepted(&contact, &contact2)
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
			kademlia.nt.rt.AddContact(NewContact(NewKademliaID(id), address[0]))
			kademlia.nt.SendJoinAcceptedMessage(address[0])
		} else if bytes.Contains(b, []byte("JoinAccepted<")) {
			var newdata []byte = b[13:]
			newstring := string(newdata)
			fmt.Println(newstring)
			stringarr := strings.Split(newstring[8:(len(newstring)-1)], ",")
			id := stringarr[0]
			address := strings.Split(stringarr[1][1:], ")")
			kademlia.nt.rt.AddContact(NewContact(NewKademliaID(id), address[0]))

		} else {
			fmt.Println("Something incorect with incoming message!")
		}
	}
}