package server

import (
	"fmt"

	gen "github.com/s-vvardenfell/QuinoaServer/generated"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ParserClient struct {
	gen.ParserServiceClient
}

func NewParserClient(host, port string) *ParserClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("cannot connect to host <%s> and port <%s>: %v", host, port, err)
	}
	return &ParserClient{
		gen.NewParserServiceClient(conn),
	}
}
