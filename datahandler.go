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
			contactarr := parseTwoContacts(b,5);
			go kademlia.nt.SendPingAccepted(&contactarr[1], &contactarr[0])
			//newstring = contact(id, address)
		} else if bytes.Contains(b, []byte("PingAccepted<")) {
			contactarr := parseTwoContacts(b,13);
			fmt.Println("Ping bounced back from: "+contactarr[1].ID.String())
			//newstring = contact(id, address)
		} else if bytes.Contains(b, []byte("Find<")) {
			contactarr := parseTwoContacts(b,5);
			go kademlia.nt.SendFindAccepted(&contactarr[0], &contactarr[1])
		} else if bytes.Contains(b, []byte("FindAccepted<")) {
			var newdata []byte = b[13:]
			newstring := string(newdata)
			//fmt.Println(newstring)
			split := strings.Split(newstring, ">")
			stringarr := strings.Split(split[0][8:(len(split[0])-1)], ",")
			stringarr2 := strings.Split(split[1][8:(len(split[1])-1)], ",")
			id := stringarr[0]
			address := strings.Split(stringarr[1][1:], ")")
			id2 := stringarr2[0]
			address2 := strings.Split(stringarr2[1][1:], ")")
			contact := NewContact(NewKademliaID(id),address[0])
			contact2 := NewContact(NewKademliaID(id2),address2[0])
			var contactlist []string
			for co := 2; co < len(split); co++ {
				stringarr := strings.Split(split[co][8:(len(split[co])-1)], ",")
				id := stringarr[0]
				address := strings.Split(stringarr[1][1:], ")")
				contactlist = append(contactlist,id,address[0])
			}
			go kademlia.LookupContactAccepted(&contact, &contact2, contactlist)
		} else if bytes.Contains(b, []byte("FindData<")) {
			contactarr := parseTwoContacts(b,9);
			newstring := string(b)
			//fmt.Println(newstring)
			split := strings.Split(newstring, ">")
			hash := split[2][:40]
			if (kademlia.storage.Check(string(hash))) {
				data := kademlia.storage.Get(string(hash))
				fmt.Println("Passed data: "+string(data))
				go kademlia.nt.SendFoundDataMessage(&contactarr[1],hash,data)
			} else {
				data := kademlia.storage.Get(string(hash))
				contact := NewContact(NewKademliaID(hash),"localhost:0000")
				newclosest := kademlia.nt.rt.FindClosestContacts(contact.ID,1)
				if ((!newclosest[0].ID.Equals(kademlia.nt.rt.me.ID)) || (!newclosest[0].ID.Equals(contactarr[1].ID))) { //No looping back to ourself and contact searching for data
					//fmt.Println(newstring)
					fmt.Println("Failed data: "+string(data)+ " hash: "+string(hash)+"##ENDS##")
					go kademlia.nt.SendLookupDataMessage(&newclosest[0],&contactarr[1],hash)
				}
			}
		} else if bytes.Contains(b, []byte("FoundData<")) {
			contactarr := parseTwoContacts(b,10);
			newstring := string(b)
			split := strings.Split(newstring, ">")
			fmt.Println("Got data from: "+contactarr[1].String() + " Hash: "+split[2]+" Data: "+split[3]) 
		} else if bytes.Contains(b, []byte("Data<")) {
			var newdata []byte = b[5:]
			newstring := string(newdata)
			//fmt.Println(newstring)
			split := strings.Split(newstring, ">")
			stringarr := strings.Split(split[0][8:(len(split[0])-1)], ",")
			id := stringarr[0]
			address := strings.Split(stringarr[1][1:], ")")
			contact := NewContact(NewKademliaID(id),address[0])
			kademlia.storage.Store(string(split[1]),string(split[2]))
			data := kademlia.storage.Get(split[1])
			fmt.Println("STORING DATA, DATA: "+string(data)+" HASH: "+string(split[1])+" Address: "+kademlia.nt.rt.me.String())
			go kademlia.nt.SendStoreDataAcceptedMessage(&contact,string(split[1]),data)
		} else if bytes.Contains(b, []byte("DataAccepted<")) {
			var newdata []byte = b[13:]
			newstring := string(newdata)
			//fmt.Println(newstring)
			split := strings.Split(newstring, ">")
			stringarr := strings.Split(split[0][8:(len(split[0])-1)], ",")
			id := stringarr[0]
			address := strings.Split(stringarr[1][1:], ")")
			contact := NewContact(NewKademliaID(id),address[0])
			fmt.Println("ACCEPTED DATA, DATA: "+string(split[2])+" HASH: "+string(split[1])+" Address: "+contact.String())
		} else if bytes.Contains(b, []byte("Join<")) {
			var newdata []byte = b[5:]
			newstring := string(newdata)
			//fmt.Println(newstring)
			stringarr := strings.Split(newstring[8:(len(newstring)-1)], ",")
			id := stringarr[0]
			address := strings.Split(stringarr[1][1:], ")")
			go kademlia.JoinRecivied(address[0], id)
		} else if bytes.Contains(b, []byte("JoinAccepted<")) {
			var newdata []byte = b[13:]
			newstring := string(newdata)
			//fmt.Println(newstring)
			stringarr := strings.Split(newstring[8:(len(newstring)-1)], ",")
			id := stringarr[0]
			address := strings.Split(stringarr[1][1:], ")")
			go kademlia.JoinAccepted(address[0], id)
		} else if bytes.Contains(b, []byte("+TERMINATE+")) {
			break
		} else {
			fmt.Println("Something incorect with incoming message!")
		}
	}
}

func parseTwoContacts(bytearr []byte, index int) [2]Contact {
	var newdata []byte = bytearr[index:]
	newstring := string(newdata)
	split := strings.Split(newstring,">")
	stringarr := strings.Split(split[0][8:(len(split[0])-1)],",")
	stringarr2 := strings.Split(split[1][8:(len(split[1])-1)],",")
	id := stringarr[0]
	address := strings.Split(stringarr[1][1:], ")")
	id2 := stringarr2[0]
	address2 := strings.Split(stringarr2[1][1:], ")")
	contact := NewContact(NewKademliaID(id),address[0])
	contact2 := NewContact(NewKademliaID(id2),address2[0])
	contactarr := [2]Contact{contact,contact2}
	return contactarr
}