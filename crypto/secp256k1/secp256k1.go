package secp256k1

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type Secp256k1Key struct {
	PrivateKey *ecdsa.PrivateKey
}

func ImportKey(hex string) *Secp256k1Key {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = secp256k1.S256()
	privateKey.D, _ = new(big.Int).SetString(hex, 16)
	return &Secp256k1Key{PrivateKey: privateKey}
}

func GenerateKey() (*Secp256k1Key, error) {
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Secp256k1Key{PrivateKey: privateKey}, nil
}

func (key *Secp256k1Key) Pubkey() (*ecdsa.PublicKey, string) {
	pubkey := key.PrivateKey.PublicKey
	pubkey.Curve = secp256k1.S256()
	pubkey.X, pubkey.Y = pubkey.Curve.ScalarBaseMult(key.PrivateKey.D.Bytes())
	pubkeyHex := pubkey.X.Text(16) + pubkey.Y.Text(16)
	return &pubkey, pubkeyHex
}

func (key *Secp256k1Key) Sign(message string) string {
	r, s, err := ecdsa.Sign(rand.Reader, key.PrivateKey, []byte(message))
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

func (key *Secp256k1Key) VerifiSignature(message string) bool {
	sig := key.Sign(message)
	r, s := new(big.Int), new(big.Int)
	r, _ = r.SetString(sig[:64], 16)
	s, _ = s.SetString(sig[64:], 16)
	pubkey, _ := key.Pubkey()
	return ecdsa.Verify(pubkey, []byte(message), r, s)
}
