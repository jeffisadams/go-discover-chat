package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	chatHost "github.com/jeffisadams/go-discover-chat/src/host"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	host := chatHost.Create(ctx)

	host.SetStreamHandler("/chat/1.0", handleStream)

	h1discService, err := discovery.NewMdnsService(ctx, host, time.Second, "_host-discovery")
	if err != nil {
		log.Fatal(err)
	}
	defer h1discService.Close()

	h1handler := &rspHandler{host}
	h1discService.RegisterNotifee(h1handler)

	time.Sleep(time.Second * 5)
	store := host.Peerstore()

	// Say hello to your friends
	fmt.Println(store.Peers())
	for _, p := range store.Peers()[1:] {
		if p.Pretty() != host.ID().Pretty() {
			fmt.Printf("Sending data to: %v", p.Pretty())
			stream, err := host.NewStream(ctx, p, "/chat/1.0")
			if err != nil {
				panic(err)
			}

			// Create a buffered stream so that read and writes are non blocking.
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			// Create a thread to read and write data.
			go writeData(rw)
			go readData(rw)
		}
	}

	select {}
}

type rspHandler struct {
	host host.Host
}

func (rh *rspHandler) HandlePeerFound(pi pstore.PeerInfo) {
	// Connect will add the host to the peerstore and dial up a new connection
	// fmt.Println(fmt.Sprintf("\nhost %v connecting to %v... (blocking)", rh.host.ID(), pi.ID))

	err := rh.host.Connect(context.Background(), pi)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error when connecting peers: %v", err))
		return
	}
}

func printKnownPeers(h host.Host) {
	fmt.Println(fmt.Sprintf("\nhost %v knows:", h.ID().Pretty()))
	for _, p := range h.Peerstore().Peers() {
		fmt.Println(fmt.Sprintf(" >> %v", p.Pretty()))
	}
}

// Stream handling code
func handleStream(s net.Stream) {
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}
