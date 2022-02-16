package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/itsabgr/go-handy"
	"github.com/itsabgr/tail"
	"github.com/phayes/freeport"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

var flagAddr = flag.String("addr", "", "service address")
var flagDir = flag.String("dir", "", "service directory")
var flagCleanPeriod = flag.Duration("clean", time.Hour*24*90, "clean records inserted before this period")

func init() {
	flag.Parse()
}

var ctx, cancel = context.WithCancel(context.Background())

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)
	go func() {
		log.Println("FATAL", "signal", <-c)
		cancel()
	}()
}
func main() {
	defer handy.Catch(func(recovered interface{}) {
		log.Fatalln(recovered)
	})
	if *flagDir == "" {
		*flagDir = filepath.Join(os.TempDir(), fmt.Sprintf("tail.%d", time.Now().UnixNano()))
		defer os.ReadDir(*flagDir)
	}
	if *flagAddr == "" {
		*flagAddr = fmt.Sprintf("0.0.0.0:%d", freeport.GetPort())
	}
	log.Println("INFO", "directory", *flagDir)
	storage, err := tail.OpenStorage(*flagDir)
	if err != nil {
		log.Fatalln(err)
	}
	defer storage.Close()
	defer log.Println("INFO", "closing")
	core := tail.NewCore(storage)
	before := time.Now().Add(-*flagCleanPeriod)
	log.Println("INFO", "cleaning before", before.Format(time.UnixDate))
	n, err := core.Clean(before)
	log.Println("INFO", "cleaned", n)
	if err != nil {
		log.Fatalln(err)
	}
	service := tail.NewServer(ctx, core)
	defer service.Close()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", freeport.GetPort()))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("INFO", "address", *flagAddr)
	service.Listen(addr.String())
}
