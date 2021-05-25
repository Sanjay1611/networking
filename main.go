package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	maddr "github.com/multiformats/go-multiaddr"
)

var protocolID string = "badgerDNS/1.0"
var db *badger.DB
var startTime time.Time

type multiAddressList []maddr.Multiaddr

func (ml *multiAddressList) String() string {
	strs := make([]string, len(*ml))
	for i, addr := range *ml {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (ml *multiAddressList) Set(value string) error {
	address, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*ml = append(*ml, address)
	return nil
}

func newHost(port int) host.Host {
	host, err := libp2p.New(context.Background(), libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("We are hosted at", host.ID())
	return host
}

func bootstrap(host host.Host, bootstrapPeers multiAddressList) (*dht.IpfsDHT, *discovery.RoutingDiscovery) {
	hostDHT, err := dht.New(context.Background(), host)
	if err != nil {
		panic(err)
	}

	if err = hostDHT.Bootstrap(context.Background()); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(context.Background(), *peerinfo); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	routingDiscovery := discovery.NewRoutingDiscovery(hostDHT)
	return hostDHT, routingDiscovery
}

func main() {
	// fmt.Println(2)
	// var x []int

	cidflag := flag.Bool("c", false, "destination peer address")
	port := flag.Int("p", 9900, "Port on which it will listen")
	dir := flag.String("dir", "/tmp/badger", "Badger DB directory")
	var bootstrapPeers multiAddressList = dht.DefaultBootstrapPeers
	flag.Var(&bootstrapPeers, "peers", "")
	flag.Parse()
	host := newHost(*port)
	var err error

	hostDHT, routingDiscovery := bootstrap(host, bootstrapPeers)

	db, err = badger.Open(badger.DefaultOptions(*dir))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	buf := bufio.NewReader(os.Stdin)

	if *cidflag == false {
		writeMode(buf, routingDiscovery, host)
	} else {
		readMode(buf, routingDiscovery, host, hostDHT)
	}
}
