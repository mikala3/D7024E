package main

const alpha = 3
const k = 4

type Kademlia struct {
	nt *Network
	closestContact *Contact
	index int
	isold bool
}

// NewNetwork returns a new instance of a RNetwork
func NewKademlia(nt *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.nt = nt
	return kademlia
}

//Send out first parallel search for contact
func (kademlia *Kademlia) LookupContact(target *Contact, sender *Contact) {
	kademlia.nt.SendFindContactMessage(target, sender)
}

//Contacts are in format [id, address, id2, address2 ...]
//Contacts have responded
func (kademlia *Kademlia) LookupContactAccepted(target *Contact, sender *Contact, contacts []string) {
	for co := 0; co < len(contacts); co=co+2 {
		id := contacts[co]
		address := contacts[co+1]
		contact := NewContact(NewKademliaID(id),address)
		kademlia.nt.rt.AddContact(contact)
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
