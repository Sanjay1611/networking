package main

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
)

type PingArgs struct {
	CID string
}

type PingReply struct {
	Name string
}

type PingService struct{}

func (t *PingService) Ping(ctx context.Context, argType PingArgs, replyType *PingReply) error {
	fmt.Println("Received a Ping call")

	name, err := getNameFromCID(db, argType.CID)
	if err != nil {
		return err
	}
	(*replyType).Name = name
	return nil
}

func makeRPCserver(host host.Host) {
	rpcHost := gorpc.NewServer(host, protocol.ID(protocolID))
	rpcHost.Register(&PingService{})
}

func callRPC(client host.Host, peerID peer.ID, cidstring string) (string, error) {
	rpcClient := gorpc.NewClient(client, protocol.ID(protocolID))
	args := &PingArgs{
		CID: cidstring,
	}
	resp := &PingReply{}
	err := rpcClient.Call(peerID, "PingService", "Ping", args, resp)
	if err != nil {
		return "", err
	}

	fmt.Println("Name for cid", cidstring, "is", resp.Name)
	fmt.Println("Time to fetch the name is", time.Since(startTime))
	return resp.Name, nil
}
