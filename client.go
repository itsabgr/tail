package tail

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io"
	"net"
	"time"
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	_ = conn.SetNoDelay(true)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Write(b)
	if err != nil {
		return err
	}
	_ = conn.CloseWrite()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
	b, err = io.ReadAll(conn)
	if err != nil {
		return err
	}
	if !bytes.Equal(b, []byte("OK")) {
		return errors.New(string(b))
	}
	return nil
}
