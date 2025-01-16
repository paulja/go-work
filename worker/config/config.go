package config

import (
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
	"time"
)

func GetLeaderAddr() string {
	v := os.Getenv("LEADER_ADDR")
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

var name string

func GetName() string {
	v := os.Getenv("WORKER_NAME")
	if v == "" {
		if len(name) == 0 {
			num := rand.IntN(1025)
			name = fmt.Sprintf("WORKER_%d", num)
		}
		return name
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

func GetAddr() string {
	return fmt.Sprintf("%s:%d", GetLocalAddr(), GetWorkerPort())
}

func GetWorkerPort() int {
	v, err := strconv.Atoi(os.Getenv("WORKER_PORT"))
	if err != nil {
		return 40041
	}
	return v
}

func GetLocalAddr() string {
	v := os.Getenv("LOCAL_ADDR")
	if v == "" {
		addr, err := findLocalIP()
		if addr == "" || err != nil {
			return "127.0.0.1"
		}
		return addr
	}
	return v
}

func GetHeartbeatTimeout() time.Duration {
	v, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIMEOUT"))
	if v <= 0 || err != nil {
		return 15
	}
	return time.Duration(v)
}

func GetEnvironment() string {
	v := os.Getenv("ENV")
	if v == "" {
		return "development"
	}
	return v
}

func findLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
					return ip.IP.To4().String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("failed to find ip address")
}
