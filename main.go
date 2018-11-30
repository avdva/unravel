// Copyright 2018 Aleksandr Demakin. All rights reserved.

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avdva/unravel/card"
	"github.com/avdva/unravel/card/printer"
	"github.com/avdva/unravel/hash"
	"github.com/avdva/unravel/server"
)

func main() {
	flagAddr := flag.String("addr", ":5000", "should be a valid addr to serve on")
	flagHash := flag.String("hash", "pjw", "a hash algorithm")
	flag.Parse()
	srv := server.New(*flagAddr, makeHandler(*flagHash))
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		srv.Stop()
	}()
	log.Fatal(srv.Serve())
}

func makeHandler(hashAlg string) card.Handler {
	hasher := hash.MakeHasher(hashAlg)
	if hasher == nil {
		log.Fatal("invalid hash alg")
	}
	l := log.New(os.Stderr, "unravel", log.LstdFlags)
	return printer.New(os.Stdout, hasher, l)
}
