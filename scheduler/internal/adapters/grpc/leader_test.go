package grpc_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/paulja/go-work/proto/cluster/v1"
	"github.com/paulja/go-work/scheduler/config"
	grpcint "github.com/paulja/go-work/scheduler/internal/adapters/grpc"
	"github.com/paulja/go-work/scheduler/internal/adapters/membership"
	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestLeader(t *testing.T) {
	t.Run("can start and stop leader", func(t *testing.T) {
		logger := createLogger()
		m := membership.NewAdapter()
		l := grpcint.NewLeaderServer(logger, m)
		assert.NoError(t, l.Start(), "should be able to start leader")
		assert.NoError(t, l.Stop(), "should be able to stop leader")
	})
	t.Run("grpc: can join leader", func(t *testing.T) {
		client, stop := testSetupLeaderClient(t)
		defer stop()

		_, err := client.Join(context.Background(), &cluster.JoinRequest{
			Id:      "1",
			Address: "localhost",
		})
		assert.NoError(t, err, "should be able to join leader")
	})
	t.Run("grpc: invaild join leader request", func(t *testing.T) {
		client, stop := testSetupLeaderClient(t)
		defer stop()

		_, err := client.Join(context.Background(), &cluster.JoinRequest{
			Id: "",
		})
		assert.ErrorIs(t, err, status.Errorf(
			codes.InvalidArgument, domain.ErrIdRequired.Error()),
			"should return id required error",
		)
		_, err = client.Join(context.Background(), &cluster.JoinRequest{
			Id:      "1",
			Address: "",
		})
		assert.ErrorIs(t, err, status.Errorf(
			codes.InvalidArgument, domain.ErrAddressRequired.Error()),
			"should return address required error",
		)
	})
	t.Run("grpc: can leave leader", func(t *testing.T) {
		client, stop := testSetupLeaderClient(t)
		defer stop()

		ctx := context.Background()
		_, err := client.Join(ctx, &cluster.JoinRequest{
			Id:      "1",
			Address: "localhost",
		})
		assert.NoError(t, err, "should be able to join leader")
		_, err = client.Leave(ctx, &cluster.LeaveRequest{
			Id: "1",
		})
		assert.NoError(t, err, "should be able to leave leader")
	})
	t.Run("grpc: invalid leave leader request", func(t *testing.T) {
		client, stop := testSetupLeaderClient(t)
		defer stop()

		ctx := context.Background()
		_, err := client.Join(ctx, &cluster.JoinRequest{
			Id:      "1",
			Address: "localhost",
		})
		assert.NoError(t, err, "should be able to join leader")
		_, err = client.Leave(ctx, &cluster.LeaveRequest{
			Id: "",
		})
		assert.ErrorIs(t, err, status.Errorf(
			codes.InvalidArgument, domain.ErrIdRequired.Error()),
			"should return id required error",
		)
	})
	t.Run("grpc: can list members", func(t *testing.T) {
		client, stop := testSetupLeaderClient(t)
		defer stop()

		ctx := context.Background()
		_, err := client.Join(ctx, &cluster.JoinRequest{
			Id:      "1",
			Address: "localhost",
		})
		assert.NoError(t, err, "should be able to join leader")
		resp, err := client.Members(ctx, &cluster.MembersRequest{})
		assert.NoError(t, err, "should be able to list members")
		assert.Len(t, resp.Members, 1, "should have 1 member")
	})
	t.Run("grpc: set heartbeat", func(t *testing.T) {
		client, stop := testSetupLeaderClient(t)
		defer stop()

		os.Setenv("HEARTBEAT_TIMEOUT", "1") // set heartbeat timeout to 1 second
		ctx := context.Background()

		// join and check the member status is as expected
		_, err := client.Join(ctx, &cluster.JoinRequest{
			Id:      "1",
			Address: "localhost",
		})
		assert.NoError(t, err, "should be able to join leader")
		resp, err := client.Members(ctx, &cluster.MembersRequest{})
		assert.NoError(t, err, "should be able to list members")
		assert.Equal(t,
			"unknown, unknown", resp.Members[0].Status,
			"should have UNKNOWN UNKNOWN status",
		)

		// set heartbeat and check the member status is as expected
		_, err = client.Heartbeat(ctx, &cluster.HeartbeatRequest{
			Id:     "1",
			Status: cluster.HeartbeatStatus_IDLE,
		})
		assert.NoError(t, err, "should be able to set heartbeat")
		resp, err = client.Members(ctx, &cluster.MembersRequest{})
		assert.NoError(t, err, "should be able to list members")
		assert.Equal(t, "alive, idle", resp.Members[0].Status, "should have ALIVE IDLE status")

		// wait for the heartbeat timeout
		time.Sleep(1100 * time.Millisecond)

		// check the member status is as expected
		resp, err = client.Members(ctx, &cluster.MembersRequest{})
		assert.NoError(t, err, "should be able to list members")
		assert.Equal(t, "left, unknown", resp.Members[0].Status, "should have LEFT UNKNOWN status")
	})
	t.Run("grpc: invalid heartbeat request", func(t *testing.T) {
		client, stop := testSetupLeaderClient(t)
		defer stop()

		ctx := context.Background()
		_, err := client.Heartbeat(ctx, &cluster.HeartbeatRequest{
			Id: "",
		})
		assert.ErrorIs(t, err, status.Errorf(
			codes.InvalidArgument, domain.ErrIdRequired.Error()),
			"should return id required error",
		)
		_, err = client.Heartbeat(ctx, &cluster.HeartbeatRequest{
			Id:     "1",
			Status: cluster.HeartbeatStatus_UNSPECIFIED,
		})
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		_, err = client.Heartbeat(ctx, &cluster.HeartbeatRequest{
			Id:     "1",
			Status: 5,
		})
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		_, err = client.Heartbeat(ctx, &cluster.HeartbeatRequest{
			Id:     "1",
			Status: cluster.HeartbeatStatus_IDLE,
		})
		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}

func testSetupLeaderClient(t *testing.T) (cluster.LeaderServiceClient, func()) {
	logger := createLogger()

	m := membership.NewAdapter()
	l := grpcint.NewLeaderServer(logger, m)
	assert.NoError(t, l.Start(), "failed to start leader")
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", config.GetLeaderPort()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "failed to connect to leader")

	stopFunc := func() {
		// close remote listener
		assert.NoError(t, l.Stop(), "failed to stop leader")
		// close local connection
		assert.NoError(t, conn.Close(), "failed to close connection")
	}
	client := cluster.NewLeaderServiceClient(conn)
	return client, stopFunc
}
