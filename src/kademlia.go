package main

const alpha = 3
const k = 4

type Kademlia struct {
	nt *Network
	closestContact *Contact
	oldclosestContact *Contact
	shortlist []Contact
	alreadycontacted []Contact
	index int
	firstrun bool
}

// NewNetwork returns a new instance of a RNetwork
func NewKademlia(nt *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.nt = nt
	return kademlia
}

//Send out first parallel search for contact
func (kademlia *Kademlia) LookupContact(target *Contact, sender *Contact) {
	if (kademlia.firstrun && len(kademlia.shortlist) < k) {
		kademlia.firstrun = false
		kademlia.alreadycontacted = kademlia.alreadycontacted[:0]
		kademlia.shortlist = kademlia.nt.rt.FindClosestContacts(target.ID, alpha)
		kademlia.closestContact = &kademlia.shortlist[0]
		for co := 0; co < len(kademlia.shortlist); co++ {
			kademlia.nt.SendFindContactMessage(target, &kademlia.shortlist[co], sender)
		}
		kademlia.index = kademlia.index + 1
	} else if (len(kademlia.shortlist) < k) {
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
			kademlia.nt.rt.AddContact(contact)
			kademlia.shortlist = append(kademlia.shortlist, contact)
			distanceclosest := target.ID.CalcDistance(kademlia.closestContact.ID)
			distancecontact := target.ID.CalcDistance(contact.ID)
			if (distancecontact.Less(distanceclosest)){
				kademlia.closestContact = &contact
			}
		}
	} else if (kademlia.index == alpha) { //Done with parallel search
		if (kademlia.closestContact.ID.Equals(kademlia.oldclosestContact.ID)) { 
			//Ids match, no new closer node found during parallel search
			kademlia.index = 0
			kademlia.firstrun = true
		} else { //Continue with parallel searches
			kademlia.index = 0
			kademlia.oldclosestContact = kademlia.closestContact
			kademlia.LookupContact(target,sender)
		}
	} else { //Done with lookup
		kademlia.index = 0
		kademlia.firstrun = true
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
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
}

func Contains(a []Contact, x Contact) bool {
    for _, n := range a {
        if x == n {
            return true
        }
    }
    return false
}
