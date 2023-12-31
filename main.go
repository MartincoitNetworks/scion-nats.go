package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Parent context cancels connecting/reconnecting altogether.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var nc *nats.Conn

	argsOnly := os.Args[1:]
	if len(argsOnly) == 0 {
		log.Fatal("Please pass SCION Address:Port (e.g. 17-ffaa:1:1,[127.0.0.1]:4222)")
	}

	fullAddr := strings.Split(argsOnly[0], ",")
	localAddr := fullAddr[len(fullAddr)-1]
	scionAddr := fullAddr[0]

	cd := &customDialer{
		ctx:             ctx,
		scionAddr:       scionAddr,
		connectTimeout:  10 * time.Second,
		connectTimeWait: 1 * time.Second,
	}
	opts := []nats.Option{
		nats.SetCustomDialer(cd),
		nats.ReconnectWait(5 * time.Second),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Println("Reconnected to", c.ConnectedUrl())
		}),
		nats.DisconnectHandler(func(c *nats.Conn) {
			log.Println("Disconnected from NATS")
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			log.Println("NATS connection is closed.")
		}),
		nats.NoReconnect(),
	}
	go func() {
		nc, err = nats.Connect(localAddr, opts...)
	}()

WaitForEstablishedConnection:
	for {
		if err != nil {
			log.Fatal(err)
		}

		// Wait for context to be canceled either by timeout
		// or because of establishing a connection...
		select {
		case <-ctx.Done():
			break WaitForEstablishedConnection
		default:
		}

		if nc == nil || !nc.IsConnected() {
			log.Println("Connection not ready")
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break WaitForEstablishedConnection
	}
	if ctx.Err() != nil {
		log.Fatal(ctx.Err())
	}

	for {
		if nc.IsClosed() {
			break
		}
		if err := nc.Publish("hello", []byte("world")); err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		log.Println("Published message")
		time.Sleep(1 * time.Second)
	}

	// Disconnect and flush pending messages
	if err := nc.Drain(); err != nil {
		log.Println(err)
	}
	log.Println("Disconnected")
}
