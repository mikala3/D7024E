package main

import (
	"testing"
	"time"
	"reflect"
)

func TestPing(t *testing.T){
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Second)
	recivier := NewRandomKademliaID()
	
	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)

	go ka.DataHandler()

	ka.nt.kademliaChannel <- ([]byte("Ping<contact("+sender.String()+", localhost:8080)>contact("+recivier.String()+", localhost:8085)"))
	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("PingAccepted<contact("+recivier.String()+", localhost:8085)>contact("+sender.String()+", localhost:8080)")))) {		
		t.Errorf("Ping test failed"+string(msg))
	} else {
		t.Logf("success ping test")
	}

}

func TestJoin(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Second)
	recivier := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)

	go ka.DataHandler()

	ka.nt.kademliaChannel <- ([]byte("Join<contact("+recivier.String()+", localhost:8085)"))
	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("JoinAccepted<contact("+sender.String()+", localhost:8080)")))) {
		t.Errorf("Join test failed"+string(msg))
	} else {
		t.Logf("Success join test")
	}
	closestid := ka.nt.rt.FindClosestContacts(recivier,1)
	if (!closestid[0].ID.Equals(recivier)) {
		t.Errorf("Join test failed"+string(msg))
	}
}

func TestFind(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Second)
	recivier := NewRandomKademliaID()
	time.Sleep(2 * time.Second)
	rtContactId := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)

	ka.nt.rt.AddContact(NewContact(recivier, "localhost:8085"))

	go ka.DataHandler()

	rtContact := NewContact(rtContactId, "localhost:8090")

	go ka.LookupContact(&rtContact,&ka.nt.rt.me)

	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("Find<contact("+rtContactId.String()+", localhost:8090)>contact("+sender.String()+", localhost:8080)")))) {
		t.Errorf("Find test failed"+string(msg))
	} else {
		t.Logf("Success find test")
	}
	//Write to internal channel, findaccepted with contact searched
	//Check kbucket after
}

func TestFindAccepted(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Second)
	recivier := NewRandomKademliaID()
	time.Sleep(2 * time.Second)
	rtContact := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)

	ka.nt.rt.AddContact(NewContact(rtContact, "localhost:8090"))

	go ka.DataHandler()

	ka.nt.kademliaChannel <- ([]byte("Find<contact("+rtContact.String()+", localhost:8090)>contact("+recivier.String()+", localhost:8085)"))
	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("FindAccepted<contact("+recivier.String()+", localhost:8085)>contact("+rtContact.String()+", localhost:8090)>contact("+rtContact.String()+", localhost:8090)")))) {
		t.Errorf("Find accepted test failed"+string(msg))
	} else {
		t.Logf("Success find accepted test")
	}
}

func TestStoreObjects(t *testing.T) {
	//TODO
}

func TestPut(t *testing.T) {
	//TODO
}

func TestGet(t *testing.T) {
	//TODO
}

func TestExit(t *testing.T) {
	//TODO
}

func TestContact(t *testing.T) {
	//TODO
}

func TestBucket(t *testing.T) {
	//TODO
}

func TestSwitch(t *testing.T) {
	//TODO
}

func TestDistance(t *testing.T) {
	//TODO
}
