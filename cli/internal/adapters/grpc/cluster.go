package grpc

import (
	"context"

	"github.com/paulja/go-work/cli/config"
	"github.com/paulja/go-work/proto/cluster/v1"
	"github.com/paulja/go-work/shared/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClusterClient struct {
	conn   *grpc.ClientConn
	client cluster.LeaderServiceClient
}

func (c *ClusterClient) Connect() error {
	cliTLS, err := tls.CliTLSConfig(config.GetServerName())
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(
		config.GetClusterAddr(),
		grpc.WithTransportCredentials(credentials.NewTLS(cliTLS)),
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
