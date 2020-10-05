package main

type Kademlia struct {
	nt *Network
}

// NewNetwork returns a new instance of a RNetwork
func NewKademlia(nt *Network) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.nt = nt
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) Ping(target *Contact) {
	// TODO
}
