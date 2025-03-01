package main

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/host"
	discovery "github.com/libp2p/go-libp2p-discovery"
	mh "github.com/multiformats/go-multihash"
)

func writeMode(buf *bufio.Reader, routingDiscovery *discovery.RoutingDiscovery, host host.Host) {
	go makeRPCserver(host)
	for {
		fmt.Println("enter name")
		name, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		go writeNameWithCID(routingDiscovery, name)
	}
}

func writeNameWithCID(routingDiscovery *discovery.RoutingDiscovery, name string) {
	defer measureTime()()
	name = strings.Trim(name, " \n")
	if name == "" {
		fmt.Println("Blank content not allowed")
		return
	}
	fmt.Println("name is", []byte(name))
	c, err := generateCID(name)
	if err != nil {
		fmt.Println("could not generate cid")
		return
	}

	err = setValue(db, c.String(), name)
	if err != nil {
		fmt.Println("could not save cid", err)
		return
	}

	discovery.Advertise(context.Background(), routingDiscovery, c.String())
	fmt.Printf("Successfully announced cid %s!\n", c.String())
}

func generateCID(name string) (cid.Cid, error) {
	pref := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}

	return pref.Sum([]byte(name))
}

func measureTime() func() {
	startTime := time.Now()
	return func() {
		fmt.Println("Time taken is", time.Since(startTime))
	}
}
