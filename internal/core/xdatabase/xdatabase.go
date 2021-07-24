package xdatabase

import (
	"../xjson"
	"../xpath"
	structPure "../xstruct/model_independent"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"strconv"
	"strings"
	"time"
)

// ObjectStore puts any interface in database
func ObjectStore(arg interface{}) error {
	freshData := xjson.StructToJson(arg)
	structName := structPure.GetStructNameOfInterface(arg)
	structName = strings.TrimPrefix(structName, "*")
	totalNumber := 0
	totalKey := fmt.Sprintf("Total:%s", structName)

	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(totalKey))
		if e == nil {
			v, _ := item.ValueCopy(nil)
			totalNum, e1 := strconv.Atoi(string(v))
			if e1 != nil {
				panic(e1)
			} else {
				totalNumber = totalNum
			}
		} else if e == badger.ErrKeyNotFound {
			e1 := txn.Set([]byte(totalKey), []byte("0"))
			if e1 != nil {
				panic(e1)
			}
			e1 = txn.Commit()
			if e1 != nil {
				panic(e1)
			}
		} else {
			panic(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.Update(func(txn *badger.Txn) error {
		for i := 0; i < totalNumber + 1; i++ {
			k := fmt.Sprintf("%s:%d", structName, i)
			item, e := txn.Get([]byte(k))
			if e == nil {
				// key has been found
				e1 := item.Value(func(val []byte) error {
					if string(val) == freshData {
						return errors.New(fmt.Sprintf("duplicate object [%s (%s)]", item.String(), structName))
					}
					return nil
				})
				if e1 != nil {
					return e1
				}
			} else if e == badger.ErrKeyNotFound {
				// key not found & is after the last
				if i == totalNumber {
					e1 := txn.Set([]byte(k), []byte(freshData))
					if e1 != nil {
						panic(e1)
					}
					e1 = txn.Set([]byte(totalKey), []byte(strconv.Itoa(totalNumber+1)))
					e1 = txn.Commit()
					if e1 != nil {
						panic(e1)
					}
					break
				}
			} else {
				panic(e)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.Sync()
	if err != nil {
		panic(err)
	}

	return nil
}

// ObjectUpdateOnKey updates any interface in database
func ObjectUpdateOnKey(arg interface{}, key string) error {
	freshData := xjson.StructToJson(arg)
	structName := structPure.GetStructNameOfInterface(arg)
	structName = strings.TrimPrefix(structName, "*")

	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		_, e := txn.Get([]byte(key))
		if e == nil {
			// key has been found
			e1 := txn.Set([]byte(key), []byte(freshData))
			if e1 != nil {
				panic(e1)
			}
			e1 = txn.Commit()
			if e1 != nil {
				panic(e1)
			}
		} else if e == badger.ErrKeyNotFound {
			// key is fresh
			return errors.New(fmt.Sprintf("can't find the key '%s' (%s)", key, structName))
		} else {
			panic(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.Sync()
	if err != nil {
		panic(err)
	}

	return nil
}

// ObjectGetKey gets the key of any interface
func ObjectGetKey(arg interface{}) (key string, err error) {
	key = ""
	freshData := xjson.StructToJson(arg)
	structName := structPure.GetStructNameOfInterface(arg)
	structName = strings.TrimPrefix(structName, "*")
	totalNumber := 0
	foundKey := ""

	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_ = db.View(func(txn *badger.Txn) error {
		totalKey := fmt.Sprintf("Total:%s", structName)
		item, e := txn.Get([]byte(totalKey))
		if e == nil {
			v, _ := item.ValueCopy(nil)
			totalNum, e1 := strconv.Atoi(string(v))
			if e1 != nil {
				panic(e1)
			} else {
				totalNumber = totalNum
			}
		} else if e == badger.ErrKeyNotFound {
			// key not found
		} else {
			panic(e)
		}
		return nil
	})

	found := false
	_ = db.View(func(txn *badger.Txn) error {
		for i := 0; i < totalNumber; i++ {
			k := fmt.Sprintf("%s:%d", structName, i)
			item, e := txn.Get([]byte(k))
			if e == nil {
				// key has been found
				_ = item.Value(func(val []byte) error {
					if string(val) == freshData {
						found = true
						foundKey = k
					}
					return nil
				})
				if found {
					break
				}
			} else if e == badger.ErrKeyNotFound {
				// key not found
			} else {
				panic(e)
			}
		}
		return nil
	})

	if found {
		return foundKey, nil
	} else {
		return "", errors.New(fmt.Sprintf("object not found (%s)", arg))
	}
}

// ObjectGetAll fetches all of any interface
func ObjectGetAll(structName string) (values []interface{}, keys []string, err error) {
	result := make([]interface{}, 0)
	resultKeys := make([]string, 0)
	totalNumber := 0

	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_ = db.View(func(txn *badger.Txn) error {
		totalKey := fmt.Sprintf("Total:%s", structName)
		item, e := txn.Get([]byte(totalKey))
		if e == nil {
			// key has been found
			v, _ := item.ValueCopy(nil)
			totalNum, e1 := strconv.Atoi(string(v))
			if e1 != nil {
				panic(e1)
			} else {
				totalNumber = totalNum
			}
		} else if e == badger.ErrKeyNotFound {
			// key not found
		} else {
			panic(e)
		}
		return nil
	})

	_ = db.View(func(txn *badger.Txn) error {
		for i := 0; i < totalNumber; i++ {
			k := fmt.Sprintf("%s:%d", structName, i)
			item, e := txn.Get([]byte(k))
			if e == nil {
				// key has been found
				_ = item.Value(func(val []byte) error {
					var r interface{}
					xjson.JsonToStruct(string(val), &r)
					result = append(result, r)
					resultKeys = append(resultKeys, k)
					return nil
				})
			} else if e == badger.ErrKeyNotFound {
				// key not found
			} else {
				panic(e)
			}
		}
		return nil
	})

	return result, resultKeys, nil
}

// ObjectGetByKey fetches an interface by id
func ObjectGetByKey(key string) (value interface{}, err error) {
	var result interface{}

	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(key))
		if e == nil {
			// key has been found
			v, _ := item.ValueCopy(nil)
			xjson.JsonToStruct(string(v), &result)
			return nil
		} else if e == badger.ErrKeyNotFound {
			// key not found
			return e
		} else {
			panic(e)
		}
	})
	if err == nil {
		return result, nil
	} else {
		return nil, err
	}
}

// ObjectDeleteByKey deletes an interface by id
func ObjectDeleteByKey(key string) error {
	db, err := badger.Open(badger.DefaultOptions(xpath.GetInternalDatabaseFilepath()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		_, e := txn.Get([]byte(key))
		if e == nil {
			// key has been found
			e1 := txn.Delete([]byte(key))
			if e1 != nil {
				panic(e1)
			}
			e1 = txn.Commit()
			if e1 != nil {
				panic(e1)
			}
			return nil
		} else if e == badger.ErrKeyNotFound {
			// key not found
			return e
		} else {
			panic(e)
		}
	})
	if err != nil {
		return err
	} else {
		err = db.Sync()
		if err != nil {
			panic(err)
		}

		return nil
	}
}

func ParamStore(value, sensorName, nodeKey string) {
	db, err := badger.Open(badger.DefaultOptions(xpath.GetNodeDatabaseFilepath(nodeKey)))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_ = db.Update(func(txn *badger.Txn) error {
		stamp := nowJSTimestamp(time.Now())
		k := fmt.Sprintf("%s:%d", sensorName, stamp)
		e := txn.Set([]byte(k), []byte(value))
		if e != nil {
			panic(e)
		}
		e = txn.Set([]byte(fmt.Sprintf("%s:Last", sensorName)), []byte(fmt.Sprintf("%d=%s", stamp, value)))
		if e != nil {
			panic(e)
		}
		e = txn.Commit()
		if e != nil {
			panic(e)
		}
		return nil
	})

	err = db.Sync()
	if err != nil {
		panic(err)
	}
}

func ParamGetAll(sensorName, nodeKey string) (values []string, timestamps []int64) {
	result := make([]string, 0)
	resultTimestamp := make([]int64, 0)

	db, err := badger.Open(badger.DefaultOptions(xpath.GetNodeDatabaseFilepath(nodeKey)))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(sensorName)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			_ = item.Value(func(val []byte) error {
				result = append(result, string(val))
				stamp, _ := strconv.ParseInt(strings.Split(string(key), ":")[1], 10, 64)
				resultTimestamp = append(resultTimestamp, stamp)
				return nil
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return result, resultTimestamp
}

func ParamGetLast(sensorName, nodeKey string) (value string, timestamp int64, err error) {
	var (
		result string
		resultTimestamp int64
	)

	db, err := badger.Open(badger.DefaultOptions(xpath.GetNodeDatabaseFilepath(nodeKey)))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		k := fmt.Sprintf("%s:Last", sensorName)
		item, e := txn.Get([]byte(k))
		if e == nil {
			byts, _ := item.ValueCopy(nil)
			temp := strings.Split(string(byts), "=")[0]
			resultTimestamp, _ = strconv.ParseInt(temp, 10, 64)
			result = strings.Split(string(byts), "=")[1]
			return nil
		} else if e == badger.ErrKeyNotFound {
			return e
		} else {
			panic(e)
		}
	})
	if err == nil {
		return result, resultTimestamp, nil
	} else {
		return "", 0, err
	}
}
