package p2p

// Peer is an interface that represent the remote node.
type Peer interface {}

// Transport is an interface that handles the communication
// between the remotenodes in the network
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}