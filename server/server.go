package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/ynsssss/ethe/server/genproto"
	"github.com/ynsssss/ethe/server/ethclient"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	apiKey := os.Getenv("INFURA_API_KEY")
	client, err := ethclient.NewEthereumBlockchainClient(context.Background(), apiKey)
	if err != nil {
		log.Fatal(err)
	}
	accService := accountService{
		ethclient: &client,
	}

	pb.RegisterAccountServiceServer(grpcServer, &accService)

	log.Println("gRPC server is running on port 8080...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
