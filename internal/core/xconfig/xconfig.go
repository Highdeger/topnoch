package xconfig

import (
	"../xpath"
	"github.com/dgraph-io/badger"
)

func ConfigPut(key, value string) {
	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	if err == nil {
		e := db.Sync()
		if e != nil {
			panic(e)
		}
	} else {
		panic(err)
	}
}

func ConfigGet(key, defaultValue string) (string, error) {
	result := defaultValue

	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(key))
		if e == nil {
			v, _ := item.ValueCopy(nil)
			result = string(v)
			return nil
		} else {
			return e
		}
	})
	if err == nil {
		return result, nil
	} else {
		return defaultValue, err
	}
}
