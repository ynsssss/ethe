package main

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func generateWallet() (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func getAddress(privateKey *ecdsa.PrivateKey) (string, error) {
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("couldnt cast public key")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex(), nil
}
