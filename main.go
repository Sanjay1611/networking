package main

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	maddr "github.com/multiformats/go-multiaddr"
)

var protocolID string = "badgerDNS/1.0"
var db *badger.DB

// var startTime time.Time

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

func newHost(port int, cidflag *bool) host.Host {
	var host host.Host
	privateKey, err := ioutil.ReadFile(PRIVATE_KEY_PATH)
	if err == nil {
		nonce := privateKey[:12]
		privateKey = privateKey[12:]

		var password string
		fmt.Printf("Please enter the password to open the file\n")
		fmt.Scanf("%s", &password)

		block, err := aes.NewCipher([]byte(password))
		if err != nil {
			log.Fatalln(err)
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("len of PRIVATE KEY ", aesgcm.Overhead())
		privateKeyHost, err := aesgcm.Open(nil, nonce, privateKey, nil)
		if err != nil {
			log.Fatalln(err)
		}

		hostKey, err := crypto.UnmarshalPrivateKey([]byte(privateKeyHost))
		if err != nil {
			log.Println("Cannot unmarshal private key")
		} else {
			host, err = libp2p.New(context.Background(), libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)), libp2p.Identity(hostKey))
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		host, err = libp2p.New(context.Background(), libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)))
		if err != nil {
			log.Fatalln(err)
		}

		privKey := host.Peerstore().PrivKey(host.ID())
		sk, err := crypto.MarshalPrivateKey(privKey)
		if err != nil {
			log.Println(err)
		}

		var password string
		fmt.Printf("Please enter a password to encrypt the file\n")
		fmt.Scanf("%s", &password)

		block, err := aes.NewCipher([]byte(password))
		if err != nil {
			log.Fatalln(err)
		}

		nonce := make([]byte, 12)
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err)
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Overhead is", aesgcm.Overhead())
		encryptedText := aesgcm.Seal(nonce, nonce, sk, nil)

		err = ioutil.WriteFile(PRIVATE_KEY_PATH, encryptedText, 0600)
		if err != nil {
			log.Fatalln(err)
		}
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

	var err error
	db, err = badger.Open(badger.DefaultOptions(*dir))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	host := newHost(*port, cidflag)
	hostDHT, routingDiscovery := bootstrap(host, bootstrapPeers)

	buf := bufio.NewReader(os.Stdin)

	if *cidflag == false {
		writeMode(buf, routingDiscovery, host)
	} else {
		readMode(buf, routingDiscovery, hostDHT)
	}
}
