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


/*func TestExit(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))
	time.Sleep(2 * time.Millisecond)
	recivier := NewRandomKademliaID()
	
	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)

	go ka.DataHandler()
	ka.nt.terminate = true
	ka.nt.kademliaChannel <- ([]byte("+TERMINATE+"))

	ka.nt.kademliaChannel <- ([]byte("Ping<contact("+sender.String()+", localhost:8080)>contact("+recivier.String()+", localhost:8085)"))
	msg := <- ka.nt.externalChannel
	if (!reflect.DeepEqual(msg, ([]byte("PingAccepted<contact("+recivier.String()+", localhost:8085)>contact("+sender.String()+", localhost:8080)")))) {		
		t.Errorf("Ping test failed"+string(msg))
	} else {
		t.Logf("success ping test")
	}
} */

func TestStoreModule(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	hash := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	rand := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)
	ka.storage.Store(hash.String(),"TestStore")
	ka.storage.Store(rand.String(),"WrongStore")
	data := ka.storage.Get(hash.String())
	fmt.Println("Data: "+string(data))
	data2 := ka.storage.Get(rand.String())
	if (string(data) != "TestStore") {
		t.Errorf("TestStoreModule test failed: data "+string(data))
	}
	if (string(data2) == "TestStore") {
		t.Errorf("TestStoreModule test failed: data "+string(data))
	}
	if (string(data) == "WrongStore") {
		t.Errorf("TestStoreModule test failed: data "+string(data))
	}
	if (string(data2) == "TestStore") {
		t.Errorf("TestStoreModule test failed: data "+string(data))
	} else {
		t.Logf("Success store module test")
	}
}

func TestCheckModule(t *testing.T) {
	var port string = "8080"

	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	hash := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	rand := NewRandomKademliaID()

	rt := NewRoutingTable(NewContact(sender, "localhost:"+port))

	kc := make(chan []byte)
	ex := make(chan []byte)
	nt := NewNetwork(rt,kc,ex)
	nt.testing = true;

	ka := NewKademlia(nt)
	ka.storage.Store(hash.String(),"TestStore")
	if (!ka.storage.Check(hash.String())) {
		t.Errorf("TestCheckModule test failed: hash "+hash.String())
	}
	if (ka.storage.Check(rand.String())) {
		t.Errorf("TestCheckModule test failed: rand "+rand.String())
	} else {
		t.Logf("Success check module test")
	}
}

func TestParseTwoContacts(t *testing.T) {
	sender := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivier := NewRandomKademliaID()
	data := ([]byte("FindData<contact("+sender.String()+", localhost:8080)>contact("+recivier.String()+", localhost:8085)>"))
	contactarr := parseTwoContacts(data,9)
	if (!contactarr[0].ID.Equals(sender)) {
		t.Errorf("TestParseTwoContacts test failed: contactarr[0].ID "+contactarr[0].ID.String())
	}
	if (contactarr[0].Address != "localhost:8080") {
		t.Errorf("TestParseTwoContacts test failed: contactarr[0].Address "+contactarr[1].Address)
	}
	if (!contactarr[1].ID.Equals(recivier)) {
		t.Errorf("TestParseTwoContacts test failed: contactarr[1].ID "+contactarr[1].ID.String())
	}
	if (contactarr[1].Address != "localhost:8085") {
		t.Errorf("TestParseTwoContacts test failed: contactarr[1].Address "+contactarr[1].Address)
	} else {
		t.Logf("Success parse two test")
	}
}

func TestContains(t *testing.T) {
	senderId := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivierId := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	rtContactId := NewRandomKademliaID()

	sender := NewContact(senderId, "localhost:8080")
	recivier := NewContact(recivierId, "localhost:8085")
	rtContact := NewContact(rtContactId, "localhost:8090")

	var List1 []Contact
	var List2 []Contact

	List1 = []Contact{sender,recivier,rtContact}
	List2 = []Contact{sender,recivier}
	
	if (!Contains(List1,rtContact)) {
		fmt.Println(List1)
		t.Errorf("TestContains test failed "+rtContact.String())
	}
	if (Contains(List2,rtContact)) {
		t.Errorf("TestContains test failed "+rtContact.String())
	} else {
		t.Logf("Success contains test")
	}
}

func TestContainsSame(t *testing.T) {
	senderId := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	recivierId := NewRandomKademliaID()
	time.Sleep(2 * time.Millisecond)
	rtContactId := NewRandomKademliaID()

	sender := NewContact(senderId, "localhost:8080")
	recivier := NewContact(recivierId, "localhost:8085")
	rtContact := NewContact(rtContactId, "localhost:8090")

	List1 := []Contact{sender,recivier,rtContact}
	List2 := []Contact{sender,recivier,rtContact}

	List3 := []Contact{sender,rtContact,recivier}

	List4 := []Contact{sender,rtContact}

	if (!ContainsSame(List1,List2)) {
		t.Errorf("TestContainsSame test failed")
	}
	if (!ContainsSame(List2,List3)) {
		t.Errorf("TestContainsSame test failed")
	}
	if (ContainsSame(List2,List4)) {
		t.Errorf("TestContainsSame test failed")
	} else {
		t.Logf("Success contains same test")
	}
}

