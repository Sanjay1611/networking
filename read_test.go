package main

import (
	"fmt"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func BenchmarkFetchName(b *testing.B) {

	options := badger.DefaultOptions("tmp/badger")
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	host := newHost(12000)
	var bootstrapPeers multiAddressList = dht.DefaultBootstrapPeers
	hostDHT, routingDiscovery := bootstrap(host, bootstrapPeers)
	cidstring := "bafkreiacd6yejuocc26acwqe7tuqm773ccjfgnu4rqfif77fsqnkj3l7u4"

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fetchName(routingDiscovery, cidstring, hostDHT)
	}
	b.StopTimer()
}
