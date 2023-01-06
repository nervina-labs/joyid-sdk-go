package signer

import (
	"encoding/hex"
	"errors"

	"github.com/nervina-labs/joyid-sdk-go/aggregator"
	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256k1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256r1"
	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

type unlockMode uint8

const (
	native unlockMode = 1
	subkey unlockMode = 2

	testnetAggregatorUrl = "https://cota.nervina.dev/aggregator"
	mainnetAggreagtorUrl = "https://cota.nervina.dev/mainnet-aggregator"
)

type AlgPrivKey struct {
	PrivKey string
	Alg     alg.AlgIndex
}

func SignNativeUnlockTx(tx *types.Transaction, algKey AlgPrivKey, webAuthn *WebAuthnMsg) error {
	if algKey.Alg == alg.Secp256r1 {
		key := secp256r1.ImportKey(algKey.PrivKey)
		return signSecp256r1Tx(tx, key, native, webAuthn)
	}
	key := secp256k1.ImportKey(algKey.PrivKey)
	return SignSecp25k1Tx(tx, key, native)
}

func SignSubkeyUnlockTx(tx *types.Transaction, algKey AlgPrivKey, addr *address.Address, webAuthn *WebAuthnMsg) error {
	pubkeyHash := secp256r1.ImportKey(algKey.PrivKey).PubkeyHash()
	rpc := aggregator.NewRPCClient(testnetAggregatorUrl)
	if addr.Network == types.NetworkMain {
		rpc = aggregator.NewRPCClient(mainnetAggreagtorUrl)
	}
	unlockSmt, err := rpc.GetSubkeyUnlockSmt(addr, pubkeyHash, algKey.Alg)
	if err != nil {
		return err
	}
	witnesses := tx.Witnesses
	if len(witnesses) < 1 {
		return errors.New("first witness cannot be empty")
	}
	firstWitnessArgs, err := types.DeserializeWitnessArgs(witnesses[0])
	if err != nil {
		return errors.New("first witness must be WitnessArgs")
	}
	unlockBytes, err := hex.DecodeString(unlockSmt)
	if err != nil {
		return err
	}
	firstWitnessArgs.OutputType = unlockBytes
	tx.Witnesses[0] = firstWitnessArgs.Serialize()

	if algKey.Alg == alg.Secp256r1 {
		key := secp256r1.ImportKey(algKey.PrivKey)
		return signSecp256r1Tx(tx, key, subkey, webAuthn)
	}
	key := secp256k1.ImportKey(algKey.PrivKey)
	return SignSecp25k1Tx(tx, key, subkey)
}
