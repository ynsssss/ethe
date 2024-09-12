package signature

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// validateSignature validates the Ethereum signature for the given address
func ValidateSignature(address string, signature string, message string) (bool, error) {
	data := []byte(message)
	hash := crypto.Keccak256Hash(data)

	sign, err := hexutil.Decode(signature) 
	if err != nil {
		return false, err
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), sign)
	if err != nil {
		return false, err
	}

	ad := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	if ad.String() == address {
		return true, nil
	}
	return false, nil
}

func hashMessage(message string) []byte {
	// a standard prefix is added to the message
	// https://ethereum.stackexchange.com/a/84960
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	prefixedMessage := []byte(prefix + message)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(prefixedMessage)
	return hash.Sum(nil)
}
