package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/renort/restart-service/api/proto/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRestartServiceServer
}

func (*server) SendMessage(stream pb.RestartService_SendMessageServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		fmt.Printf("Server.SendMessage: %v\n", in)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr := ":" + os.Getenv("HOST_PORT")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	pb.RegisterRestartServiceServer(s, &server{})

	var group errgroup.Group
	group.Go(func() error {
		return s.Serve(lis)
	})

	mux := runtime.NewServeMux()
	err = pb.Register(ctx, mux, addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	group.Go(func() error {
		return http.ListenAndServe(":"+os.Getenv("GATEWAY_PORT"), mux)
	})

	if err = group.Wait(); err != nil {
		log.Fatalln(err)
	}

}
