package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	pb "github.com/ynsssss/ethe/client/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func BenchmarkGetAccounts(b *testing.B) {
	conn, err := grpc.NewClient(
		":8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewAccountServiceClient(conn)

	benchmarks := []int{
		100,
		1000,
		10000,
	}

	tokens := []string{
		"0x0000000000b3F879cb30FE243b4Dfee438691c04", // gastoken
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // weth
		"0xdac17f958d2ee523a2206206994597c13d831ec7", // usdt
		"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", // usdc
	}

	addresses := make([]string, 0, 10000)
	for range 10000 {
		privateKey, err := generateWallet()
		if err != nil {
			log.Fatal(err)
		}

		address, err := getAddress(privateKey)
		if err != nil {
			log.Fatal(err)
		}

		addresses = append(addresses, address)
	}
	for _, bm := range benchmarks {
		b.Run(fmt.Sprint(bm), func(b *testing.B) {
			benchmarkGetAccounts(client, addresses[:bm], tokens[rand.Intn(3)])
		})
	}
}

func benchmarkGetAccounts(client pb.AccountServiceClient, addresses []string, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	stream, err := client.GetAccounts(ctx)
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}
	go func() {
		var i int
		for {
			_, err := stream.Recv()
			i++
			if i == len(addresses) {
				cancel()
			}
			if err != nil {
				if grpc.Code(err) == codes.Canceled {
					break
				}
				log.Fatalf("Failed to receive: %v", err)
			}
		}
	}()

	for _, address := range addresses {
		req := &pb.GetAccountsRequest{
			EthereumAddress:   address,
			Erc20TokenAddress: token,
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send address: %v", err)
		}
	}

}
