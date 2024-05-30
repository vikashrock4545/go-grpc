package main

import (
	"context"
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
	r.GET("/send-message-to-server/:message", clientConnectionServer)
	r.Run(":8000")
}

func clientConnectionServer(c *gin.Context) {
	variableName := c.Param("message")

	req := &proto.HelloRequest{SomeString: variableName}

	client.ServerReply(context.TODO(), req)
	c.JSON(http.StatusOK, gin.H{
		"message": "message sent successfully to server " + variableName,
	})
}
