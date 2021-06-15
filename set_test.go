package main

import (
	"fmt"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
)

func BenchmarkSet(b *testing.B) {

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
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setValue(db, c.String(), name)
	}
	b.StopTimer()
}
