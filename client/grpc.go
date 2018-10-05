package client

import (
	"context"
	"time"

	"github.com/orvice/kit/log"
	"github.com/orvice/monitor-client/proto"
	"github.com/orvice/monitor-server/mod"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	NodeID      int32
	ctx         context.Context
	cancel      func()
	grpcAddr    string
	messageRecv chan mod.Packet
	logger      log.Logger
}

func NewGrpcClient(nodeID int32, grpcAddr string, msgRecv chan mod.Packet,
	l log.Logger) *GrpcClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &GrpcClient{
		NodeID:      nodeID,
		ctx:         ctx,
		cancel:      cancel,
		grpcAddr:    grpcAddr,
		messageRecv: msgRecv,
		logger:      l,
	}

}

func (g *GrpcClient) Run() error {
	for {
		if g.RunStream() != nil {
		}
		time.Sleep(time.Second)
		continue
	}
}

func (g *GrpcClient) RunStream() error {
	conn, err := grpc.Dial(g.grpcAddr, grpc.WithInsecure())
	if err != nil {
		g.logger.Errorf("connect failed %v", err)
		return err
	}
	defer conn.Close()

	cli := monitorClient.NewMonitorClientClient(conn)

	stream, err := cli.Stream(context.Background(), &monitorClient.StreamRequest{})
	if err != nil {
		g.logger.Errorf("stream failed %v", err)
		return err
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		g.messageRecv <- mod.Packet{
			NodeID:  g.NodeID,
			Message: []byte(resp.Body),
		}
	}

	return nil
}
