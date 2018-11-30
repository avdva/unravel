// Copyright 2018 Aleksandr Demakin. All rights reserved.

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avdva/unravel/server"
)

func main() {
	flagAddr := flag.String("serve addr", ":5000", "should be a valid addr to serve on")
	flag.Parse()
	srv := server.New(*flagAddr)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		srv.Stop()
	}()
	log.Fatal(srv.Serve())
}
