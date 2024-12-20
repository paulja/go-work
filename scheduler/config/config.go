package config

import (
	"os"
	"strconv"
	"time"
)

func GetLeaderPort() int {
	v, err := strconv.Atoi(os.Getenv("LEADER_PORT"))
	if err != nil {
		return 50051
	}
	return v
}

func GetRPCPort() int {
	v, err := strconv.Atoi(os.Getenv("RPC_PORT"))
	if err != nil {
		return 50052
	}
	return v
}

func GetHeartbeatTimeout() time.Duration {
	v, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIMEOUT"))
	if err != nil {
		return 30
	}
	return time.Duration(v)
}

func GetEnvironment() string {
	v := os.Getenv("ENV")
	if v != "" {
		return "development"
	}
	return v
}
