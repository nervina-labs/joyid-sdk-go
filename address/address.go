package address

import (
	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256k1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256r1"
	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

const (
	testnetJoyidCodeHash = "0x07595bf3a836bdb7f93a91a66cb2c8cbaf0d807e8adecddc9d910aac6f0a5d0f"
	mainnetJoyidCodeHash = "0x07595bf3a836bdb7f93a91a66cb2c8cbaf0d807e8adecddc9d910aac6f0a5d0f"
)

func FromPubkeyHash(pubkeyHash []byte, algIndex alg.AlgIndex, network types.Network) *address.Address {
	var args []byte
	if algIndex == alg.Secp256r1 {
		args = []byte{0x00, 0x01}
	} else {
		args = []byte{0x00, 0x02}
	}
	args = append(args, pubkeyHash...)
	var codeHash string
	if network == types.NetworkTest {
		codeHash = testnetJoyidCodeHash
	} else {
		codeHash = mainnetJoyidCodeHash
	}
	lockScript := &types.Script{
		CodeHash: types.HexToHash(codeHash),
		HashType: types.HashTypeType,
		Args:     args,
	}
	return &address.Address{
		Script:  lockScript,
		Network: network,
	}
}

func FromPrivKey(key string, algIndex alg.AlgIndex, network types.Network) *address.Address {
	var pubkeyHash []byte
	if algIndex == alg.Secp256k1 {
		pubkeyHash = secp256k1.ImportKey(key).PubkeyHash()
	} else {
		pubkeyHash = secp256r1.ImportKey(key).PubkeyHash()
	}
	return FromPubkeyHash(pubkeyHash, algIndex, network)
}
