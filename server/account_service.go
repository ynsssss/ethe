package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/ynsssss/ethe/server/ethclient"
	pb "github.com/ynsssss/ethe/server/genproto"
	"github.com/ynsssss/ethe/server/signature"
)

type accountService struct {
	pb.UnimplementedAccountServiceServer
	ethclient *ethclient.EthereumBlockchainClient
}

var (
	ErrCouldNotVerifySignature = errors.New("error while validating signature")
	ErrInvalidSignature        = errors.New("invalid signature")
)

func (s *accountService) GetAccount(
	ctx context.Context,
	req *pb.GetAccountRequest,
) (*pb.GetAccountResponse, error) {
	// message is the ethereum address
	valid, err := signature.ValidateSignature(
		req.EthereumAddress,
		req.CryptoSignature,
		req.EthereumAddress,
	)
	if !valid {
		if err != nil {
			return nil, fmt.Errorf(ErrCouldNotVerifySignature.Error(), err)
		}
		return nil, ErrInvalidSignature
	}

	gastokenBalance, nonce, err := s.ethclient.GetAccountData(ctx, req.EthereumAddress)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		GastokenBalance: gastokenBalance,
		WalletNonce:     nonce,
	}, nil
}

func (s *accountService) GetAccounts(stream pb.AccountService_GetAccountsServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		balance, err := s.ethclient.GetERC20Balance(
			context.Background(),
			req.EthereumAddress,
			req.Erc20TokenAddress,
		)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.GetAccountsResponse{
			EthereumAddress: req.EthereumAddress,
			Erc20Balance:    balance.String(),
		})
		if err != nil {
			return err
		}
	}
}
