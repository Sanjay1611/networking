package main

import badger "github.com/dgraph-io/badger/v3"

func getNameFromCID(db *badger.DB, cidstring string) (string, error) {
	var name []byte

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(cidstring))
		if err != nil {
			return err
		}
		name, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return string(name), err
}
