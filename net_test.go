package tail

import (
	"context"
	"fmt"
	"github.com/phayes/freeport"
	"net"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestNet(t *testing.T) {
	storage, err := OpenStorage(filepath.Join(t.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	if err != nil {
		panic(err)
	}
	defer storage.Close()
	service := NewServer(context.Background(), NewCore(storage))
	defer service.Close()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", freeport.GetPort()))
	if err != nil {
		panic(err)
	}
	go service.Listen(addr.String())
	runtime.Gosched()
	poller, err := NewPoller(":0")
	if err != nil {
		t.Fatal(err)
	}
	defer poller.Close()
	err = Push(addr.String(), []byte("hi"))
	if err != nil {
		t.Fatal(err)
	}
	var last []byte
	b, err := poller.Poll(last, addr)
	fmt.Println(b, err)
}
