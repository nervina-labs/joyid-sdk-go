package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

type Secp256r1Key struct {
	PrivateKey *ecdsa.PrivateKey
}

func ImportKey(hex string) *Secp256r1Key {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P256()
	privateKey.D, _ = new(big.Int).SetString(hex, 16)
	return &Secp256r1Key{PrivateKey: privateKey}
}

func GenerateKey() (*Secp256r1Key, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Secp256r1Key{PrivateKey: privateKey}, nil
}

func (key *Secp256r1Key) Pubkey() (*ecdsa.PublicKey, string) {
	pubkey := key.PrivateKey.PublicKey
	pubkey.Curve = elliptic.P256()
	pubkey.X, pubkey.Y = pubkey.Curve.ScalarBaseMult(key.PrivateKey.D.Bytes())
	pubkeyHex := pubkey.X.Text(16) + pubkey.Y.Text(16)
	return &pubkey, pubkeyHex
}

func (key *Secp256r1Key) Sign(message string) string {
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

func (key *Secp256r1Key) VerifySignature(message string) bool {
	sig := key.Sign(message)
	r, s := new(big.Int), new(big.Int)
	r, _ = r.SetString(sig[:64], 16)
	s, _ = s.SetString(sig[64:], 16)
	pubkey, _ := key.Pubkey()
	return ecdsa.Verify(pubkey, []byte(message), r, s)
}
