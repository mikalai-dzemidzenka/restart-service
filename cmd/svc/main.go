package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/renort/restart-service/api/proto/pb"
	"google.golang.org/grpc"
)

func getGwAddress() string {
	builder := strings.Builder{}
	builder.WriteString(os.Getenv("GRPC_ADDRESS"))
	builder.WriteString(":")
	builder.WriteString(os.Getenv("GRPC_PORT"))
	return builder.String()
}

func main() {
	conn, err := grpc.Dial(getGwAddress(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection error: %v\n", err)
	}
	defer conn.Close()
	client := pb.NewRestartServiceClient(conn)
	stream, err := client.SendMessage(context.Background())
	//TODO backoff
	if err != nil {
		log.Fatalf("Send message error: %v\n", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Client.SendMessage: %v\n", in.Body)
		}
	}()
	if err := stream.Send(&pb.Message{Body: "client"}); err != nil {
		log.Fatalln(err)
	}
	stream.CloseSend()
	<-waitc
	fmt.Println("Shutting down...")
}
