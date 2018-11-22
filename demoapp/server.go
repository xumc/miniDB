package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"net/http"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Execute(ctx context.Context, in *SqlRequest) (*SqlReply, error) {
	return &SqlReply{}, nil
}

type httpHandler struct{}

func (h httpHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("<h1>hello</h1>"))
}

func Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	m := cmux.New(lis)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	grpcS := grpc.NewServer()
	RegisterExecutorServer(grpcS, &server{})
	reflection.Register(grpcS)

	httpS := &http.Server{
		Handler: &httpHandler{},
	}

	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	if err := m.Serve(); err != nil {
		fmt.Print("server done")
	}
}
