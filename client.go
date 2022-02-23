package tail

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"net"
	"net/http"
	"strconv"
)

type Poller struct {
	conn net.PacketConn
}

func NewPoller(addr string) (*Poller, error) {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err

	}
	return &Poller{conn: conn}, nil
}
func (poller *Poller) Close() error {
	return poller.conn.Close()
}

func (poller *Poller) Poll(last []byte, addr net.Addr) ([]byte, error) {
	_, err := poller.conn.WriteTo(last, addr)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 2048)
	n, _, err := poller.conn.ReadFrom(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func Push(addr string, b []byte) error {
	resp, err := http.DefaultClient.Post(addr, "text/plain", bytes.NewReader(b))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
	}
	return nil
}
