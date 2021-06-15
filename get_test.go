package main

import (
	"fmt"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
)

func BenchmarkGet(b *testing.B) {

	options := badger.DefaultOptions("tmp/badger")
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	name := "friends"
	c, err := generateCID(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	setValue(db, c.String(), name)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		getValue(db, "bafkreib5cwdhwhdgkubpseag55yzattjl6fvryaso4w7hyqhgyibezi5se")
	}
	b.StopTimer()
}
