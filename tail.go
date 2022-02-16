package tail

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Core struct {
	storage *Storage
}

func NewCore(storage *Storage) *Core {
	return &Core{storage: storage}
}
func (core *Core) Put(b []byte, time time.Time) error {
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(v, uint64(time.Unix()))
	return core.storage.Put(b, v)
}

func (core *Core) Get(last []byte) ([]byte, error) {
	k, _, err := core.storage.Get(last)
	if err == ErrNotFound {
		err = nil
	}
	return k, err
}
func (core *Core) Clean(before time.Time) (int, error) {
	a := make([]byte, 8)
	var n int
	binary.BigEndian.PutUint64(a, uint64(before.Unix()))
	var last []byte
	for {
		key, b, err := core.storage.Get(last)
		if err != nil {
			if err == ErrNotFound {
				break
			}
			return n, err
		}
		if bytes.Compare(a, b) != -1 {
			n++
			core.storage.Purge(key)
		}
		last = key
	}
	return n, nil
}
