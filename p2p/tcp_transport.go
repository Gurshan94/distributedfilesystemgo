package p2p

import (
	"fmt"
	"net"
	"io"
	"errors"
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

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts 
	listener net.Listener
	rpcch chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 1024),
	}
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
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
	defer conn.Close()

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("Connection closed by peer: %s\n", conn.RemoteAddr())
			} else {
				fmt.Printf("Error decoding message: %s\n", err)
			}
			return
		}

		rpc.From = conn.RemoteAddr().String()
		t.rpcch <- rpc
	}
}