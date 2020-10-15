package main

import (
	"testing"
	"time"
	"reflect"
	"fmt"
)

func TestPing(t *testing.T){
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
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
	time.Sleep(2 * time.Millisecond)
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

// NEW FORMAT FOR FIND, OLD WAS REEEEEAAAAALLLYYYY BAD
// Find<contact(SEARCHED)>contact(SENDER)
// FindAccepted<contact(SEARCHED)>contact(SENDER)>contact(BUCKET1)>contact(BUCKET2) ...
// Recivier is no longer included since they are directly connected, the sender is still sent by function call to network.

func TestFind(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivier := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
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
		t.Errorf("Find test failed: msg "+string(msg))
	}
	//Write to internal channel, findaccepted with contact searched
	//Check kbucket after

	//Send find accept with contact searched for in bucket (sent from recivier)
	ka.nt.kademliaChannel <- ([]byte("FindAccepted<contact("+rtContactId.String()+", localhost:8090)>contact("+recivier.String()+", localhost:8085)>contact("+rtContactId.String()+", localhost:8090)"))
	fmt.Println("Before channel")
	fmt.Println(ka.kaalpha)
	msg2 := <- ka.nt.externalChannel
	fmt.Println("After channel")
	if (!reflect.DeepEqual(msg2, ([]byte("Find<contact("+rtContactId.String()+", localhost:8090)>contact("+sender.String()+", localhost:8080)")))) {
		t.Errorf("Find test failed: msg2 "+string(msg2))
	}
	ka.nt.kademliaChannel <- ([]byte("FindAccepted<contact("+rtContactId.String()+", localhost:8090)>contact("+rtContactId.String()+", localhost:8090)>contact("+recivier.String()+", localhost:8085)"))
	for {
		if (ka.firstrun == true) {break}
	}
	contactsInRt := ka.nt.rt.FindClosestContacts(rtContactId,10)
	if (!contactsInRt[0].ID.Equals(rtContactId)) {
		t.Errorf("Find test failed: rt "+string(contactsInRt[0].String()))
	} else {
		t.Logf("Success find test")
	}
}

func TestFindAccepted(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivier := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	rtContact := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)

	ka.nt.rt.AddContact(NewContact(recivier, "localhost:8085"))
	ka.nt.rt.AddContact(NewContact(rtContact, "localhost:8090"))

	go ka.DataHandler()

	ka.nt.kademliaChannel <- ([]byte("Find<contact("+rtContact.String()+", localhost:8090)>contact("+recivier.String()+", localhost:8085)"))
	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("FindAccepted<contact("+rtContact.String()+", localhost:8090)>contact("+sender.String()+", localhost:8080)>contact("+rtContact.String()+", localhost:8090)>contact("+recivier.String()+", localhost:8085)")))) {
		t.Errorf("Find accepted test failed"+string(msg))
	} else {
		t.Logf("Success find accepted test")
	}
}

func TestStoreObjects(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivier := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	hash := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))
	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)
	ka.nt.rt.AddContact(NewContact(recivier, "localhost:8085"))

	go ka.DataHandler()

	ka.nt.kademliaChannel <- ([]byte("Data<contact("+recivier.String()+", localhost:8085)>"+hash.String()+">dettaardata"))

	ka.nt.kademliaChannel <- ([]byte("FindData<contact("+sender.String()+", localhost:8080)>contact("+recivier.String()+", localhost:8085)>"+hash.String()))
	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("FoundData<contact("+recivier.String()+", localhost:8085)>contact("+sender.String()+", localhost:8080)>"+hash.String()+">dettaardata")))) {
		t.Errorf("Store test failed"+string(msg))
	} else {
		t.Logf("Success store test")
	}
}

func TestStoreNotFound(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivier := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	hash := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))
	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)
	ka.nt.rt.AddContact(NewContact(recivier, "localhost:8085"))

	go ka.DataHandler()

	ka.nt.kademliaChannel <- ([]byte("FindData<contact("+sender.String()+", localhost:8080)>contact("+recivier.String()+", localhost:8085)>"+hash.String()))
	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("FindData<contact("+recivier.String()+", localhost:8085)>contact("+recivier.String()+", localhost:8085)>"+hash.String())))) {
		t.Errorf("Store not found test failed"+string(msg))
	} else {
		t.Logf("Success store not found test")
	}
}

func TestFullStore(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivier := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	rtContactId := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)

	ka.nt.rt.AddContact(NewContact(recivier, "localhost:8085"))

	go ka.DataHandler()

	go ka.Store(rtContactId.String(), ([]byte("supersecret")))

	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("Find<contact("+rtContactId.String()+", localhost:0000)>contact("+sender.String()+", localhost:8080)")))) {
		t.Errorf("FullStore test failed: msg "+string(msg))
	}
	//Write to internal channel, findaccepted with contact searched
	//Check kbucket after

	//Send find accept with contact searched for in bucket (sent from recivier)
	ka.nt.kademliaChannel <- ([]byte("FindAccepted<contact("+rtContactId.String()+", localhost:0000)>contact("+recivier.String()+", localhost:8085)>contact("+rtContactId.String()+", localhost:8090)"))
	msg2 := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg2, ([]byte("Find<contact("+rtContactId.String()+", localhost:0000)>contact("+sender.String()+", localhost:8080)")))) {
		t.Errorf("FullStore test failed: msg2 "+string(msg2))
	}
	ka.nt.kademliaChannel <- ([]byte("FindAccepted<contact("+rtContactId.String()+", localhost:0000)>contact("+rtContactId.String()+", localhost:8090)>contact("+recivier.String()+", localhost:8085)"))
	for {
		if (ka.firstrun == true) {break}
	}
	contactsInRt := ka.nt.rt.FindClosestContacts(rtContactId,10)
	if (!contactsInRt[0].ID.Equals(rtContactId)) {
		t.Errorf("FullStore test failed: rt "+string(contactsInRt[0].String()))
	} else {
		t.Logf("Success FullStore test")
	}

	msg3 := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg3, ([]byte("Data<contact("+sender.String()+", localhost:8080)>"+rtContactId.String()+">supersecret")))) {
		t.Errorf("FullStore test failed: msg3 "+string(msg3))
	}

	msg4 := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg4, ([]byte("Data<contact("+sender.String()+", localhost:8080)>"+rtContactId.String()+">supersecret")))) {
		t.Errorf("FullStore test failed: msg4 "+string(msg4))
	}
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
