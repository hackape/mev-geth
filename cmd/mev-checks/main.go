package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	clientDial = flag.String(
		"client_dial", "ws://127.0.0.1:8545", "could be websocket or IPC",
	)
	cb = flag.String(
		"coinbase", "", "what coinbase to use",
	)
)

func program() error {
	client, err := ethclient.Dial(*clientDial)
	if err != nil {
		return err
	}

	ch := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(
		context.Background(), ch,
	)

	if err != nil {
		return err
	}

	for {
		select {
		case e := <-sub.Err():
			return e
		case incoming := <-ch:
			fmt.Println("header is ", incoming)
			if incoming.Number.Uint64() == 5 {
				if err := client.SendMegaBundle(
					context.Background(), &ethclient.MegaBundle{
						TransactionList: nil,
						Timestamp:       uint64(time.Now().Add(time.Second * 45).Unix()),
						Coinbase_diff:   3e18,
						Coinbase:        common.HexToAddress(*cb),
						ParentHash:      incoming.Root,
					},
				); err != nil {
					return err
				}
			}
		}
	}
}

func main() {
	flag.Parse()
	if err := program(); err != nil {
		log.Fatal(err)
	}
}
