package main

import (
	"context"
	//"flag"
	"log"
	"net"
	"net/netip"
	//"strings"
	"time"
	//"os"
	"crypto/tls"

	"github.com/nats-io/nats.go"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsec-ethz/scion-apps/pkg/quicutil"
	"github.com/quic-go/quic-go"
)

type customDialer struct {
	ctx             context.Context
	nc              *nats.Conn
	scionAddr       string
	connectTimeout  time.Duration
	connectTimeWait time.Duration
}

func (cd *customDialer) Dial(network, address string) (net.Conn, error) {
	// Simplistic policy
	policy, err := pan.PolicyFromCommandline("", "latency", false)
	log.Println(policy)
	log.Println("this is the address", address)
	log.Println("this is the network", cd.scionAddr)
	// The address can be of the form of a SCION address (i.e. of the form "ISD-AS,[IP]:port")
	fullAddr := cd.scionAddr + "," + address
	addr, err := pan.ResolveUDPAddr(context.TODO(), fullAddr)
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{quicutil.SingleStreamProto},
		//NextProtos:         []string{"hello-quic"},
	}
	quicConf := &quic.Config{
		KeepAlivePeriod: 15 * time.Second,
	}
	log.Println("resolveUPDAdrr:", addr)
	// TODO: debug why quic is not connecting to loopback on remote scion host
	ql, err := pan.DialQUIC(context.Background(), netip.AddrPort{}, addr, policy, nil, "", tlsCfg, quicConf)
	if err != nil {
		return nil, err
	}
	//var dialed quicutil.SingleStream
	stream, err := quicutil.NewSingleStream(ql)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

/*
func (cd *customDialer) Dial(network, address string) (net.Conn, error) {
    ctx, cancel := context.WithTimeout(cd.ctx, cd.connectTimeout)
    defer cancel()

    for {
		log.Println("Attempting to connect to address: ", address)
		log.Println("Attempting to connect to network: ", network)
        if ctx.Err() != nil {
            return nil, ctx.Err()
        }

        select {
        case <-cd.ctx.Done():
            return nil, cd.ctx.Err()
        default:
            d := &net.Dialer{}
            if conn, err := d.DialContext(ctx, network, address); err == nil {
                log.Println("Connected to NATS successfully")
                return conn, nil
            } else {
                time.Sleep(cd.connectTimeWait)
            }
        }
    }
}
*/
