package tail

import (
	"encoding/binary"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	db, err := OpenStorage(filepath.Join(t.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for n := range make([]struct{}, 100) {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(n))
		err = db.Put(b, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	var last []byte
	for n := range make([]struct{}, 100) {
		k, _, err := db.Get(last)
		if err != nil {
			t.Fatal(err)
		}
		if binary.BigEndian.Uint64(k) != uint64(n) {
			t.FailNow()
		}
		last = k
	}
}
