package secp256k1

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/keccak"
	"github.com/nervina-labs/joyid-sdk-go/utils"
)

type Key struct {
	PrivateKey *ecdsa.PrivateKey
}

func (k *Key) Bytes() []byte {
	return math.PaddedBigBytes(k.PrivateKey.D, k.PrivateKey.Params().BitSize/8)
}

func ImportKey(privKey string) *Key {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = secp256k1.S256()
	privateKey.D, _ = new(big.Int).SetString(utils.Trim0x(privKey), 16)
	return &Key{PrivateKey: privateKey}
}

func GenerateKey() (*Key, error) {
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Key{PrivateKey: privateKey}, nil
}

func (key *Key) Pubkey() (*ecdsa.PublicKey, []byte) {
	pubkey := key.PrivateKey.PublicKey
	pubkey.Curve = secp256k1.S256()
	pubkey.X, pubkey.Y = pubkey.Curve.ScalarBaseMult(key.PrivateKey.D.Bytes())
	pubkeyBytes := pubkey.X.Bytes()
	pubkeyBytes = append(pubkeyBytes, pubkey.Y.Bytes()...)
	return &pubkey, pubkeyBytes
}

func (key *Key) PubkeyHash() []byte {
	_, pubkey := key.Pubkey()
	return keccak.Keccak160((pubkey))
}

func (key *Key) Sign(message []byte) []byte {
	sig, err := secp256k1.Sign(message, key.Bytes())
	if err != nil {
		return []byte{}
	}
	return sig
}
