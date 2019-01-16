package host

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
)

func Create(ctx context.Context) host.Host {
	port := generatePort()

	hma, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	h, err := libp2p.New(ctx, libp2p.ListenAddrs(hma))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Host %v Listening on Port: %d \n", h.ID().Pretty(), port)
	return h
}

func generatePort() int {
	source := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(source)

	return r1.Intn(10000)
}
