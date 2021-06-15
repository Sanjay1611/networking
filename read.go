package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

var memDB *badger.DB

func readMode(buf *bufio.Reader, routingDiscovery *discovery.RoutingDiscovery, hostDHT *dht.IpfsDHT) {
	var err error
	memDB, err = badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		fmt.Println("enter cid")
		cidstring, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		// startTime = time.Now()
		go fetchName(routingDiscovery, cidstring, hostDHT)
	}
}

func fetchName(routingDiscovery *discovery.RoutingDiscovery, cidstring string, hostDHT *dht.IpfsDHT) {
	defer measureTime()()
	cidstring = strings.Trim(cidstring, " \n")
	if cidstring == "" {
		fmt.Println("Invalid cidstring")
		return
	}

	name, err := getValue(memDB, cidstring)
	if err == nil {
		fmt.Println("Name for cid", cidstring, "is", name)
		return
	} else if err != nil && err != badger.ErrKeyNotFound {
		log.Println("Error accessing local db ", err)
	}

	fmt.Printf("searching for peers with %s\n", cidstring)
	peerChan, err := routingDiscovery.FindPeers(context.Background(), cidstring)
	if err != nil {
		panic(err)
	}
	for peer := range peerChan {
		if peer.ID == hostDHT.Host().ID() || len(peer.Addrs) == 0 {
			continue
		}
		fmt.Println("Connecting to:", peer)

		addrInfo, err := hostDHT.FindPeer(context.Background(), peer.ID)
		if err != nil {
			fmt.Println("cannot find address of peer", peer)
			continue
		}

		if err := hostDHT.Host().Connect(context.Background(), addrInfo); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Connected to:", peer)
		name, err := callRPC(hostDHT.Host(), peer.ID, cidstring)
		if err != nil {
			log.Println(err)
		} else {
			memDB.Update(func(txn *badger.Txn) error {
				entry := badger.NewEntry([]byte(cidstring), []byte(name)).WithTTL(time.Second * 15)
				err := txn.SetEntry(entry)
				return err
			})
			return
		}
	}
}
