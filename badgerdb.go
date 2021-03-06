package kvbase

import (
	"encoding/json"
	"errors"
	"github.com/dgraph-io/badger"
	"strings"
)

type BadgerBackend struct {
	Backend
	Connection *badger.DB
}

func NewBadgerDB(source string) (Backend, error) {
	if source == "" {
		source = "data"
	}

	opts := badger.DefaultOptions(source)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	database := BadgerBackend{
		Connection: db,
	}

	return &database, nil
}

func (database *BadgerBackend) Count(bucket string) (int, error) {
	db := database.Connection
	counter := 0

	return counter, db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			if strings.HasPrefix(string(it.Item().Key()), bucket) {
				counter++
			}
		}
		return nil
	})
}

func (database *BadgerBackend) Create(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err == nil {
		return errors.New("key already exists")
	}

	return database.write(bucket, key, model)
}

func (database *BadgerBackend) Delete(bucket string, key string) error {
	db := database.Connection

	if _, err := database.view(bucket, key); err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(bucket + "_" + key))
	})
}

func (database *BadgerBackend) Get(bucket string, model interface{}) (*map[string]interface{}, error) {
	db := database.Connection
	results := make(map[string]interface{})

	return &results, db.View(func(txn *badger.Txn) error {
		prefix := []byte(bucket + "_")
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			if err := item.Value(func(value []byte) error {
				key := strings.Replace(string(item.Key()), bucket+"_", "", -1)

				err := json.Unmarshal(value, &model)
				if err != nil {
					return err
				}

				results[key] = model

				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (database *BadgerBackend) Read(bucket string, key string, model interface{}) error {
	data, err := database.view(bucket, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &model)
}

func (database *BadgerBackend) Update(bucket string, key string, model interface{}) error {
	if _, err := database.view(bucket, key); err != nil {
		return err
	}

	return database.write(bucket, key, model)
}

func (database *BadgerBackend) view(bucket string, key string) ([]byte, error) {
	db := database.Connection
	var data []byte

	return data, db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(bucket + "_" + key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			data = append([]byte{}, val...)

			return nil
		})
	})
}

func (database *BadgerBackend) write(bucket string, key string, model interface{}) error {
	db := database.Connection

	data, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(bucket+"_"+key), data)
	})
}