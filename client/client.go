package main

import (
	"context"
	"fmt"
	proto "go-grpc-prac/proto"
	"io"
	"log"
	"net/http"
	"strconv"

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
	stream, err := client.ServerReply(context.TODO())

	if err != nil {
		fmt.Println("Something error")
		return
	}
	send, receive := 0, 0
	for i := 0; i < 10; i++ {
		err := stream.Send(&proto.HelloRequest{
			SomeString: "message " + strconv.Itoa(i+1) + " from clients",
		})
		if err != nil {
			fmt.Println("failed to send message from client")
			return
		}
		send++
	}

	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println("Server response: ", message)
		receive++
	}

	c.JSON(http.StatusOK, gin.H{
		"message_sent":     send,
		"message_received": receive,
	})
}
