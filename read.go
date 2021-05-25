package main

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func readMode(buf *bufio.Reader, routingDiscovery *discovery.RoutingDiscovery, host host.Host, hostDHT *dht.IpfsDHT) {
	for {
		fmt.Println("enter cid")
		cidstring, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		startTime = time.Now()
		cidstring = strings.Trim(cidstring, " \n")
		if cidstring == "" {
			fmt.Println("Invalid cidstring")
			continue
		}
		go func() {
			fmt.Printf("searching for peers with %s\n", cidstring)
			peerChan, err := routingDiscovery.FindPeers(context.Background(), cidstring)
			if err != nil {
				panic(err)
			}
			for peer := range peerChan {
				if peer.ID == host.ID() || len(peer.Addrs) == 0 {
					continue
				}
				fmt.Println("Connecting to:", peer)

				addrInfo, err := hostDHT.FindPeer(context.Background(), peer.ID)
				if err != nil {
					fmt.Println("cannot find address of peer", peer)
					continue
				}

				if err := host.Connect(context.Background(), addrInfo); err != nil {
					fmt.Println(err)
				}

				fmt.Println("Connected to:", peer)
				callRPC(host, peer.ID, cidstring)
			}
		}()
	}
}
