package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/nervina-labs/joyid-sdk-go/utils"
	"github.com/nervosnetwork/ckb-sdk-go/v2/crypto/blake2b"
)

type Key struct {
	PrivateKey *ecdsa.PrivateKey
}

func (k *Key) Bytes() []byte {
	return math.PaddedBigBytes(k.PrivateKey.D, k.PrivateKey.Params().BitSize/8)
}

func ImportKey(privKey string) *Key {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = elliptic.P256()
	privateKey.D, _ = new(big.Int).SetString(utils.Trim0x(privKey), 16)
	return &Key{PrivateKey: privateKey}
}

func GenerateKey() (*Key, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Key{PrivateKey: privateKey}, nil
}

func (key *Key) Pubkey() (*ecdsa.PublicKey, []byte) {
	pubkey := key.PrivateKey.PublicKey
	pubkey.Curve = elliptic.P256()
	pubkey.X, pubkey.Y = pubkey.Curve.ScalarBaseMult(key.PrivateKey.D.Bytes())
	pubkeyBytes := make([]byte, 64)
	copy(pubkeyBytes[32-len(pubkey.X.Bytes()):32], pubkey.X.Bytes())
	copy(pubkeyBytes[64-len(pubkey.Y.Bytes()):64], pubkey.Y.Bytes())
	return &pubkey, pubkeyBytes
}

func (key *Key) PubkeyHash() []byte {
	_, pubkey := key.Pubkey()
	return blake2b.Blake160(pubkey)
}

func (key *Key) Sign(message []byte) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, key.PrivateKey, message)
	if err != nil {
		return []byte{}
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sigBytes := make([]byte, 64)
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}

func (key *Key) VerifySignature(message []byte) bool {
	sig := key.Sign(message)
	r, s := new(big.Int), new(big.Int)
	r = r.SetBytes(sig[:32])
	s = s.SetBytes(sig[32:])
	pubkey, _ := key.Pubkey()
	return ecdsa.Verify(pubkey, message, r, s)
}
