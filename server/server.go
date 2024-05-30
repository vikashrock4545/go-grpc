package main

import (
	"fmt"
	proto "go-grpc-prac/proto"
	"net"
	"time"

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

func (s *server) ServerReply(req *proto.HelloRequest, strem proto.Example_ServerReplyServer) error {
	fmt.Println(req.SomeString)
	time.Sleep(5 * time.Second)
	reply := []*proto.HelloResponse{
		{Reply: "hello1"},
		{Reply: "hello2"},
		{Reply: "hello3"},
		{Reply: "hello4"},
		{Reply: "hello5"},
		{Reply: "hello6"},
	}
	for _, msg := range reply {
		err := strem.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
