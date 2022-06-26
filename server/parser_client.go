package server

import (
	"fmt"

	gen "github.com/s-vvardenfell/QuinoaServer/generated"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//здесь дб ParserClient
type ParserClient struct {
	gen.RedisCacheServiceClient
}

func NewParserClient(host, port string) *ParserClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("cannot connect to host <%s> and port <%s>: %v", host, port, err)
	}
	return &ParserClient{
		gen.NewRedisCacheServiceClient(conn),
	}
}
