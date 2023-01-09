package signer

import (
	"errors"

	"github.com/nervina-labs/joyid-sdk-go/aggregator"
	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256k1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256r1"
	"github.com/nervina-labs/joyid-sdk-go/utils"
	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

const (
	native byte = 1
	subkey byte = 2

	testnetAggregatorUrl = "http://127.0.0.1:3030"
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
	return signSecp25k1Tx(tx, key, native)
}

func SignSubkeyUnlockTx(tx *types.Transaction, algKey AlgPrivKey, webAuthn *WebAuthnMsg) error {
	if algKey.Alg == alg.Secp256r1 {
		key := secp256r1.ImportKey(algKey.PrivKey)
		return signSecp256r1Tx(tx, key, subkey, webAuthn)
	}
	key := secp256k1.ImportKey(algKey.PrivKey)
	return signSecp25k1Tx(tx, key, subkey)
}

func BuildOutputTypeWithSubkeySmt(tx *types.Transaction, algKey AlgPrivKey, addr *address.Address) error {
	var pubkeyHash []byte
	if algKey.Alg == alg.Secp256k1 {
		pubkeyHash = secp256k1.ImportKey(algKey.PrivKey).PubkeyHash()
	} else {
		pubkeyHash = secp256r1.ImportKey(algKey.PrivKey).PubkeyHash()
	}

	var rpc *aggregator.RPCClient
	if addr.Network == types.NetworkMain {
		rpc = aggregator.NewRPCClient(mainnetAggreagtorUrl)
	} else {
		rpc = aggregator.NewRPCClient(testnetAggregatorUrl)
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
	unlockBytes, err := utils.HexToBytes(unlockSmt)

	if err != nil {
		return errors.New("hex convert error")
	}
	firstWitnessArgs.OutputType = unlockBytes
	tx.Witnesses[0] = firstWitnessArgs.Serialize()
	return nil
}
