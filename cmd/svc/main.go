package main

import (
	"log"
	"os"
	"strings"

	"github.com/renort/restart-service/api/proto/pb"
	"google.golang.org/grpc"
)

func getGwAddress() string {
	builder := strings.Builder{}
	builder.WriteString(os.Getenv("GW_ADDRESS"))
	builder.WriteString(":")
	builder.WriteString(os.Getenv("GW_PORT"))
	return builder.String()
}

func main() {
	conn, err := grpc.Dial(getGwAddress(), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewRestartServiceClient(conn)
	client.
}
