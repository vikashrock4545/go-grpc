package main

import (
	"context"
	"fmt"
	proto "go-grpc-prac/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.ExampleClient

func main() {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewExampleClient(conn)
	// implement rest-api
	r := gin.Default()
	r.GET("/send", clientConnectionServer)
	r.Run(":8000")
}

func clientConnectionServer(c *gin.Context) {
	req := []*proto.HelloRequest{
		{SomeString: "Request 1"},
		{SomeString: "Request 2"},
		{SomeString: "Request 3"},
		{SomeString: "Request 4"},
		{SomeString: "Request 5"},
		{SomeString: "Request 6"},
		{SomeString: "Request 7"},
	}

	stream, err := client.ServerReply(context.TODO())
	if err != nil {
		fmt.Println("Something error")
		return
	}
	for _, re := range req {
		err = stream.Send(re)
		if err != nil {
			fmt.Println("request not fulfilled")
			return
		}
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("there is some error ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message_count": response,
	})
}
