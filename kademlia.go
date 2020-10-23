package main

import (
	"fmt"
)

const alpha = 5
const k = 6

type Kademlia struct {
	nt *Network
	closestContact *Contact
	oldclosestContact *Contact
	shortlist []Contact
	alreadycontacted []Contact
	index int
	kaalpha int
	firstrun bool
	//ping bool
	//pingContact *Contact
	storage *Storage
}

// NewNetwork returns a new instance of a RNetwork
func NewKademlia(nt *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.nt = nt
	kademlia.firstrun = true;
	//kademlia.ping = false;
	kademlia.index = 0;
	kademlia.kaalpha = 0;
	kademlia.storage = NewStorage()
	return kademlia
}

//Send out first parallel search for contact
func (kademlia *Kademlia) LookupContact(target *Contact, sender *Contact) {
	if (kademlia.firstrun && len(kademlia.shortlist) < k) {
		kademlia.firstrun = false
		kademlia.alreadycontacted = kademlia.alreadycontacted[:0]
		kademlia.shortlist = kademlia.nt.rt.FindClosestContacts(target.ID, alpha)
		kademlia.kaalpha = len(kademlia.shortlist)
		kademlia.closestContact = &kademlia.shortlist[0]
		kademlia.oldclosestContact = &kademlia.shortlist[0]
		for co := 0; co < len(kademlia.shortlist); co++ {
			kademlia.nt.SendFindContactMessage(target, &kademlia.shortlist[co], sender)
			kademlia.alreadycontacted = append(kademlia.alreadycontacted, kademlia.shortlist[co])
		}
	} else if (len(kademlia.shortlist) < k) {
		kademlia.kaalpha = len(kademlia.shortlist) - len(kademlia.alreadycontacted)
		for co := 0; co < len(kademlia.shortlist); co++ {
			if (!Contains(kademlia.alreadycontacted, kademlia.shortlist[co])) {
				kademlia.nt.SendFindContactMessage(target, &kademlia.shortlist[co], sender)
				kademlia.alreadycontacted = append(kademlia.alreadycontacted, kademlia.shortlist[co])
			}
		}
	}
}	

//Contacts are in format [id, address, id2, address2 ...]
//Contacts have responded
func (kademlia *Kademlia) LookupContactAccepted(target *Contact, sender *Contact, contacts []string) {
	if (len(kademlia.shortlist) < k ) {
		kademlia.index = kademlia.index + 1
		for co := 0; co < len(contacts); co=co+2 {
			id := contacts[co]
			address := contacts[co+1]
			contact := NewContact(NewKademliaID(id),address)
			kademlia.shortlist = append(kademlia.shortlist, contact)
			distanceclosest := target.ID.CalcDistance(kademlia.closestContact.ID)
			distancecontact := target.ID.CalcDistance(contact.ID)
			if (distancecontact.Less(distanceclosest)){
				kademlia.closestContact = &contact
			}
			if (!contact.ID.Equals(kademlia.nt.rt.me.ID)){ //Needs to be changed so that it pings contacts in bucket and updates accordingly
				// kademlia.ping = true;
				// kademlia.Ping(&contact,&kademlia.nt.rt.me)
				kademlia.nt.rt.AddContact(contact)
			}
		}
	} 
	if ((kademlia.index == alpha) || (kademlia.index == kademlia.kaalpha)) { //Done with parallel search
		if (kademlia.closestContact.ID.Equals(kademlia.oldclosestContact.ID)) { 
			//Ids match, no new closer node found during parallel search
			fmt.Println("Done with parallel")
			kademlia.index = 0
			kademlia.kaalpha = 0
			kademlia.firstrun = true
		} else { //Continue with parallel searches
			fmt.Println("Starting next parallel")
			kademlia.index = 0
			kademlia.oldclosestContact = kademlia.closestContact
			kademlia.LookupContact(target,sender)
		}
	} else { //Done with lookup
		fmt.Println("Done with parallel: else")
		kademlia.index = 0
		kademlia.kaalpha = 0
		kademlia.firstrun = true
	}
	if (ContainsSame(kademlia.shortlist,kademlia.alreadycontacted)) {
		fmt.Println("Done with parallel")
		kademlia.index = 0
		kademlia.kaalpha = 0
		kademlia.firstrun = true
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	if (kademlia.storage.Check(hash)) {
		fmt.Println(string(kademlia.storage.Get(hash)))
	} else {
		contact := NewContact(NewKademliaID(hash),"localhost:0000")
		newclosest := kademlia.nt.rt.FindClosestContacts(contact.ID,1)
		fmt.Println("Lookup sent to: "+newclosest[0].String())
		kademlia.nt.SendLookupDataMessage(&newclosest[0],&kademlia.nt.rt.me,hash)
		
	}
}

func (kademlia *Kademlia) Store(hash string, data string) {
	contact := NewContact(NewKademliaID(hash),"localhost:0000")
	fmt.Println("Hash: "+hash)
	fmt.Println("Data: "+string(data))
	kademlia.index = 0
	kademlia.kaalpha = 0
	kademlia.firstrun = true
	kademlia.LookupContact(&contact,&kademlia.nt.rt.me)
	for {if (kademlia.firstrun == true) {break}}
	newclosest := kademlia.nt.rt.FindClosestContacts(contact.ID,k)
	for co := 0; co < len(newclosest); co++ {
		fmt.Println("Sending data to: "+newclosest[co].String())
		kademlia.nt.SendStoreDataMessage(&newclosest[co],hash,data)
	}
}

func (kademlia *Kademlia) Ping(target *Contact, sender *Contact) {
	go kademlia.nt.SendPingMessage(target, sender)
}

func (kademlia *Kademlia) Join(ip string, id string) {
	kademlia.nt.rt.AddContact(NewContact(NewKademliaID(id), ip))
	kademlia.nt.SendJoinAcceptedMessage(ip)
}

func (kademlia *Kademlia) JoinAccepted(ip string, id string) {
	kademlia.nt.rt.AddContact(NewContact(NewKademliaID(id), ip))
	kademlia.index = 0
	kademlia.kaalpha = 0
	kademlia.firstrun = true
	kademlia.LookupContact(&kademlia.nt.rt.me,&kademlia.nt.rt.me)
}

// func (kademlia *Kademlia) AddContact(contact *Contact) {
// 	bucketIndex := kademlia.nt.rt.getBucketIndex(contact.ID)
// 	bucket := kademlia.nt.rt.buckets[bucketIndex]
// 	if (bucket.Len()<bucketSize) {
// 		kademlia.ping = true
// 		kademlia.pingContact = contact
// 		kademlia.Ping(bucket.GetFirst(),&kademlia.nt.rt.me)
// 	} else {
// 		kademlia.nt.rt.AddContact(*contact)
// 	}
// }

func ContainsSame(a []Contact, x []Contact) bool {
    for _, n := range a {
        if Contains(x,n) {
			if (n.String() == a[len(a)-1].String()) {
				return true
			}
        } else {
			return false
		}
	}
	return false
}

func Contains(a []Contact, x Contact) bool {
    for _, n := range a {
		if x.String() == n.String() {
			return true
        }
    }
    return false
}
