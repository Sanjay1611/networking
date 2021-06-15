package main

import badger "github.com/dgraph-io/badger/v3"

func setValue(db *badger.DB, key, value string) error {
	return db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
}
