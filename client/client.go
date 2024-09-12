package main

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	pb "github.com/ynsssss/ethe/client/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"grpc-server:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAccountServiceClient(conn)

	callGetAccount(client)

	callGetAccounts(client)
}

func callGetAccount(client pb.AccountServiceClient) {

	privateKey, err := generateWallet()
	if err != nil {
		log.Fatal(err)
	}

	address, err := getAddress(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	//message is the address
	data := []byte(address)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	sign := hexutil.Encode(signature)

	req := &pb.GetAccountRequest{
		EthereumAddress: address,
		CryptoSignature: sign,
	}

	res, err := client.GetAccount(context.Background(), req)
	if err != nil {
		log.Fatalf("GetAccount failed: %v", err)
	}
	log.Printf("Gastoken Balance: %s, Wallet Nonce: %d", res.GastokenBalance, res.WalletNonce)
}

func callGetAccounts(client pb.AccountServiceClient) {
	stream, err := client.GetAccounts(context.Background())
	if err != nil {
		log.Fatalf("could not start GetAccounts: %v", err)
	}

	addresses := make([]string, 0, 3)
	for range 3 {
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

	for _, address := range addresses {
		req := &pb.GetAccountsRequest{
			EthereumAddress:   address,
			Erc20TokenAddress: "0x0000000000b3F879cb30FE243b4Dfee438691c04",
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatalf("error sending: %v", err)
		}

		res, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving: %v", err)
		}
		log.Printf("Address: %s, ERC20 Balance: %s", res.EthereumAddress, res.Erc20Balance)
	}
}
