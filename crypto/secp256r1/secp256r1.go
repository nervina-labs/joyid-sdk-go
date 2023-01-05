package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

func ImportPrivateKey(key string) *ecdsa.PrivateKey {
	var privateKey ecdsa.PrivateKey
	privateKey.Curve = elliptic.P256()
	privateKey.D, _ = new(big.Int).SetString(key, 16)
	return &privateKey
}

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func GetPubKey(privKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, string) {
	pubkey := privKey.PublicKey
	pubkey.Curve = elliptic.P256()
	pubkey.X, pubkey.Y = pubkey.Curve.ScalarBaseMult(privKey.D.Bytes())
	pubkeyHex := pubkey.X.Text(16) + pubkey.Y.Text(16)
	return &pubkey, pubkeyHex
}

func Sign(privKey *ecdsa.PrivateKey, message string) string {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, []byte(message))
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

func VerifiSignature(privKey *ecdsa.PrivateKey, message string) bool {
	sig := Sign(privKey, message)
	r, s := new(big.Int), new(big.Int)
	r, _ = r.SetString(sig[:64], 16)
	s, _ = s.SetString(sig[64:], 16)
	pubkey, _ := GetPubKey(privKey)
	return ecdsa.Verify(pubkey, []byte(message), r, s)
}
