package server

import (
	"fmt"

	gen "github.com/s-vvardenfell/QuinoaServer/generated"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RedisClient struct {
	gen.RedisCacheServiceClient
}

func NewRedisClient(host, port string) *RedisClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("cannot connect to host <%s> and port <%s>: %v", host, port, err)
	}
	return &RedisClient{
		gen.NewRedisCacheServiceClient(conn),
	}
}
