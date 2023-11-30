## SCiON Enabled NATS Client

This is a simple client for sending messages to [SCiON enabled NATS Servers/Clusters](https://github.com/MartincoitNetworks/nats-server) using the [Path Aware Network (PAN) Library](https://github.com/netsec-ethz/scion-apps/tree/3afc9a9118080aa78e6a6435c06549a8a6c0bd23/pkg/pan)

### Building

The following will result in a `scionnats` binary:

```
$ nix develop
$ go-build
```

### Run

Pass the address for the NATS Server to connect to:

```
$ ./scionnats 17-ffaa:1:1,[127.0.0.1]:4222
```
> Assumes the SCiON services are configured

The client will continuiously send messages to the `hello` subject

## More Info

- https://github.com/scionproto/scion -  SCiON Internet Architecture
- https://www.scionlab.org/ - An easy way to get access to a SCiON environment
- https://www.scion.org/ - SCiON Association
