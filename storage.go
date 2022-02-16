package tail

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Storage struct {
	db *leveldb.DB
}

var ErrNotFound = leveldb.ErrNotFound

func OpenStorage(path string) (storage *Storage, err error) {
	storage = &Storage{}
	storage.db, err = leveldb.OpenFile(path, &opt.Options{
		Compression: opt.NoCompression,
	})
	return storage, err
}
func (storage *Storage) Close() error {
	return storage.db.Close()
}
func (storage *Storage) Put(key []byte, val []byte) error {
	if val == nil {
		val = []byte{}
	}
	return storage.db.Put(key, val, &opt.WriteOptions{Sync: true})

}

func (storage *Storage) Get(start []byte) (key, val []byte, err error) {
	iter := storage.db.NewIterator(&util.Range{Start: nil, Limit: nil}, &opt.ReadOptions{})
	defer iter.Release()
	if false == iter.Next() {
		return nil, nil, leveldb.ErrNotFound
	}
	if bytes.Compare(iter.Key(), start) == 0 {
		if false == iter.Next() {
			return nil, nil, leveldb.ErrNotFound
		}
	}
	key = clone(iter.Key())
	val = clone(iter.Value())
	return key, val, iter.Error()
}

func (storage *Storage) Fold(start, end []byte, forEach func(key, val []byte) error) error {
	iter := storage.db.NewIterator(&util.Range{Start: start, Limit: end}, &opt.ReadOptions{})
	defer iter.Release()
	for iter.Next() {
		err := forEach(iter.Key(), iter.Value())
		if err != nil {
			return err
		}
	}
	return iter.Error()
}

func clone(b []byte) []byte {
	dst := make([]byte, len(b))
	copy(dst, b)
	return dst
}
func (storage *Storage) Purge(key []byte) error {
	return storage.db.Delete(key, &opt.WriteOptions{Sync: true})
}
