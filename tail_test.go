package tail

import (
	"encoding/binary"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestCore(t *testing.T) {
	storage, err := OpenStorage(filepath.Join(t.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	if err != nil {
		panic(err)
	}
	defer storage.Close()
	core := NewCore(storage)
	now := time.Now()
	for n := range make([]struct{}, 100) {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(n))
		err = core.Put(b, now)
		if err != nil {
			t.Fatal(err)
		}
	}
	<-time.NewTimer(2 * time.Second).C
	for n := range make([]struct{}, 40) {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(n))
		err = core.Put(b, time.Now())
		if err != nil {
			t.Fatal(err)
		}
	}
	err = core.Clean(now)
	if err != nil {
		t.Fatal(err)
	}
	var last []byte
	for n := range make([]struct{}, 100) {
		v, err := core.Get(last)
		if len(v) == 0 {
			if n == 40 {
				break
			}
			t.FailNow()
		}
		if err != nil {
			t.Fatal(err)
		}
		if binary.BigEndian.Uint64(v) != uint64(n) {
			t.FailNow()
		}
		last = v
	}
}
