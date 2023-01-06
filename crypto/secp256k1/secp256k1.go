package secp256k1

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/keccak"
)

type Key struct {
	PrivateKey *ecdsa.PrivateKey
}

func ImportKey(privKey string) *Key {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = secp256k1.S256()
	privateKey.D, _ = new(big.Int).SetString(privKey, 16)
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
	return keccak.Keccak160(([]byte(pubkey)))
}

func (key *Key) Sign(message []byte) string {
	r, s, err := ecdsa.Sign(rand.Reader, key.PrivateKey, message)
	if err != nil {
		return ""
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sigBytes := make([]byte, 64)
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return fmt.Sprintf("%x", sigBytes)
}

func (key *Key) VerifySignature(message []byte) bool {
	sig := key.Sign(message)
	r, s := new(big.Int), new(big.Int)
	r, _ = r.SetString(sig[:64], 16)
	s, _ = s.SetString(sig[64:], 16)
	pubkey, _ := key.Pubkey()
	return ecdsa.Verify(pubkey, []byte(message), r, s)
}
