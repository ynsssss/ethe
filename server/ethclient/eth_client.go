package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const gasTokenAddress = "0x0000000000b3F879cb30FE243b4Dfee438691c04"

const erc20ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"}]`

type EthereumBlockchainClient struct {
	eClient *ethclient.Client
}

func NewEthereumBlockchainClient(
	ctx context.Context,
	apiKey string,
) (EthereumBlockchainClient, error) {
	url := "https://mainnet.infura.io/v3/" + apiKey
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return EthereumBlockchainClient{}, err
	}
	return EthereumBlockchainClient{
		eClient: client,
	}, nil
}

func (client *EthereumBlockchainClient) GetAccountData(
	ctx context.Context,
	hexAddress string,
) (string, uint64, error) {
	ad := common.HexToAddress(hexAddress)
	nonce, err := client.eClient.NonceAt(ctx, ad, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get nonce: %v", err)
	}

	balance, err := client.GetERC20Balance(ctx, hexAddress, gasTokenAddress)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get balance: %v", err)
	}

	return balance.String(), nonce, nil
}

func (client *EthereumBlockchainClient) GetERC20Balance(
	ctx context.Context,
	address string,
	tokenAddress string,
) (*big.Int, error) {
	tAddress := common.HexToAddress(tokenAddress)
	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token ABI: %v", err)
	}

	instance := bind.NewBoundContract(
		tAddress,
		parsedABI,
		client.eClient,
		client.eClient,
		client.eClient,
	)

	balance := new(big.Int)
	args := []interface{}{&balance}
	err = instance.Call(&bind.CallOpts{Context: ctx}, &args, "balanceOf", address)
	return balance, nil
}
