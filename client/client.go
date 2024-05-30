package main

import (
	"context"
	"fmt"
	proto "go-grpc-prac/proto"
	"io"
	"net/http"
	"time"

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
	stream, err := client.ServerReply(context.TODO(), &proto.HelloRequest{SomeString: "Hiiii"})

	if err != nil {
		fmt.Println("Something error")
		return
	}
	count := 0
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println("Server response: ", message)
		time.Sleep(2 * time.Second)
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"message_count": count,
	})
}
