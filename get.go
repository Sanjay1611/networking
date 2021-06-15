package main

import badger "github.com/dgraph-io/badger/v3"

func getValue(db *badger.DB, key string) (string, error) {
	var val []byte

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return string(val), err
}
