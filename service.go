package tail

import (
	"context"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	core   *Core
}

func NewServer(ctx context.Context, core *Core) *Server {
	ctx, cancel := context.WithCancel(ctx)
	server := &Server{ctx: ctx, cancel: cancel, core: core}
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
func (server *Server) serveTCP(conn net.Listener) error {
	http := fasthttp.Server{}
	http.ReadTimeout = 2 * time.Second
	http.WriteTimeout = 2 * time.Second
	http.ReadTimeout = 2 * time.Second
	http.IdleTimeout = 2 * time.Second
	http.MaxRequestBodySize = 2048
	http.DisableHeaderNamesNormalizing = true
	http.NoDefaultDate = true
	http.NoDefaultServerHeader = true
	http.NoDefaultContentType = true
	http.TCPKeepalivePeriod = 2 * time.Second
	http.Handler = func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Method()) == fasthttp.MethodPost {
			err := server.core.Put(ctx.Request.Body(), time.Now())
			if err != nil {
				http.ErrorHandler(ctx, err)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusNoContent)
		} else {
			b, err := server.core.Get(ctx.Request.Body())
			if err != nil {
				http.ErrorHandler(ctx, err)
				return
			}
			if len(b) == 0 {
				ctx.SetStatusCode(fasthttp.StatusNoContent)
			}
			ctx.SetBody(b)
		}
	}
	http.ErrorHandler = func(ctx *fasthttp.RequestCtx, err error) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
	}
	return http.Serve(conn)
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
