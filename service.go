package tail

import (
	"context"
	"io"
	"net"
	"time"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	core   *Core
}

func NewServer(ctx context.Context, core *Core) *Server {
	context, cancel := context.WithCancel(ctx)
	server := &Server{ctx: context, cancel: cancel, core: core}

	return server
}
func (server *Server) Close() error {
	server.cancel()
	return nil
}
func (server *Server) handleUDP(packetConn net.PacketConn) error {
	b := make([]byte, 2048)
	for {
		n, from, err := packetConn.ReadFrom(b)
		if err != nil {
			return err
		}
		v, _ := server.core.Get(b[:n])
		if len(v) == 0 {
			continue
		}
		_, _ = packetConn.WriteTo(v, from)
	}
}

func (server *Server) handleTCPConn(conn net.Conn) error {
	defer conn.Close()
	err := conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return err
	}
	
	b, err := io.ReadAll(io.LimitReader(conn, 2048))
	if err != nil {
		return err
	}
	err = server.core.Put(b, time.Now())
	err2 := conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err2 != nil {
		return err
	}
	if err != nil {
		_, _ = io.WriteString(conn, err.Error())
		return err
	}
	_, err = io.WriteString(conn, "OK")
	return err

}
func (server *Server) serveTCP(conn net.Listener) error {
	for {
		conn, err := conn.Accept()
		if err != nil {
			return err
		}
		_ = server.handleTCPConn(conn)
	}
}
func (server *Server) Listen(addr string) error {
	packetConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer packetConn.Close()
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	go server.serveTCP(conn)
	go server.handleUDP(packetConn)
	<-server.ctx.Done()
	return server.ctx.Err()
}
