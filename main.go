package main

import (
	"fmt"
	"log"

	"github.com/Gurshan94/distributedfilesystemgo/p2p"
)

func Onpeer (peer p2p.Peer) error {
	return nil
}

func main() {
	tcpOpts:= p2p.TCPTransportOpts{
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		ListenAddr:    ":3000",
		OnPeer: Onpeer,
	}

	tr:=p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
		msg:= <- tr.Consume()
		fmt.Printf("%+v\n",msg)
	    }
    }()

	if err:=tr.ListenAndAccept(); err!=nil {
		log.Fatal(err)
	}

	select{}
}