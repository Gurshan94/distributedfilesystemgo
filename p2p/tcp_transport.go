package p2p

import (
	"fmt"
	"net"
	"sync"
)

//tcppeer represents a remote node over a tcp established network connection
type TCPPeer struct {
	conn net.Conn
 
	// outbound if we make the connection to the peer
	outbound bool

}

func NewTCPPeer (conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener net.Listener

	mu sync.RWMutex
	peers map[string]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err!=nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err!=nil {
			fmt.Printf("TCp accept error: %s\n",err)
		}

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("new TCP connection:%+v\n",peer)
}