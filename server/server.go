package main

import (
	"errors"
	"fmt"
	proto "go-grpc-prac/proto"
	"io"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedExampleServer
}

func main() {
	listener, tcpErr := net.Listen("tcp", ":9000")
	if tcpErr != nil {
		panic(tcpErr)
	}

	srv := grpc.NewServer()
	proto.RegisterExampleServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) ServerReply(strem proto.Example_ServerReplyServer) error {
	for i := 0; i < 5; i++ {
		err := strem.Send(&proto.HelloResponse{
			Reply: "message " + strconv.Itoa(i+1) + " from servers",
		})
		if err != nil {
			return errors.New("unable to send data from server ")
		}
	}
	for {
		req, err := strem.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println(req.SomeString)
	}
	return nil
}
