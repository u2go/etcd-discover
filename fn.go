package etcdDiscover

import (
	"context"
	"net"
	"os"
)

func Ctx() context.Context {
	return context.Background()
}

func Ip() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func Hostname() (string, error) {
	return os.Hostname()
}
