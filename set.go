package main

import badger "github.com/dgraph-io/badger/v3"

func setCIDwithName(db *badger.DB, cidstring, name string) error{
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(cidstring), []byte(name))
		return err
	}) 
}