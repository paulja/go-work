package config

import "os"

func GetClusterAddr() string {
	v := os.Getenv("CLUSTER_ADDR")
	if v == "" {
		return "localhost:50051"
	}
	return v
}

func GetSchedulerAddr() string {
	v := os.Getenv("SCHEDULER_ADDR")
	if v == "" {
		return "localhost:50052"
	}
	return v
}

func GetServerName() string {
	v := os.Getenv("SERVER_NAME")
	if v == "" {
		return "localhost"
	}
	return v
}
