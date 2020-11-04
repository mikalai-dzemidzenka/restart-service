package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/renort/restart-service/api/proto/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type server struct {
	mu        sync.Mutex
	multiplex []chan pb.Message
	pb.UnimplementedRestartServiceServer
}

func (s *server) SendMessage(stream pb.RestartService_SendMessageServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		fmt.Printf("Server.SendMessage: %v\n", in)
		ch := make(chan pb.Message)
		s.mu.Lock()
		s.multiplex = append(s.multiplex, ch)
		s.mu.Unlock()
		message := <-ch
		defer func() {
			s.mu.Lock()
			close(ch)
			s.multiplex = s.multiplex[:len(s.multiplex)-1]
			s.mu.Unlock()
		}()
		if err = stream.Send(&message); err != nil {
			return err
		}

		//TODO sync???? think about possible errors
	}
}

func main() {
	addr := ":" + os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	servInstance := &server{}
	pb.RegisterRestartServiceServer(s, servInstance)

	var group errgroup.Group
	group.Go(func() error {
		return s.Serve(lis)
	})

	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		var message pb.Message
		c.BindJSON(&message)
		servInstance.mu.Lock()
		for _, ch := range servInstance.multiplex {
			ch <- message
		}
		servInstance.mu.Unlock()
		//TODO response?
		c.JSON(200, gin.H{"status": "ok"})
	})
	group.Go(func() error {
		return r.Run(":" + os.Getenv("GW_PORT"))
	})

	if err = group.Wait(); err != nil {
		log.Fatalln(err)
	}
}
