package grpc

import (
	"context"

	"github.com/paulja/go-work/cli/config"
	"github.com/paulja/go-work/proto/cluster/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClusterClient struct {
	conn   *grpc.ClientConn
	client cluster.LeaderServiceClient
}

func (c *ClusterClient) Connect() error {
	conn, err := grpc.NewClient(
		config.GetClusterAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	c.conn = conn
	c.client = cluster.NewLeaderServiceClient(conn)
	return nil
}

func (c *ClusterClient) Close() error {
	return c.conn.Close()
}

func (c *ClusterClient) GetMembers() ([]*cluster.Member, error) {
	resp, err := c.client.Members(context.Background(), &cluster.MembersRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Members, nil
}
