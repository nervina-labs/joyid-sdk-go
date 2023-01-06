package address

import (
	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256k1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256r1"
	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

const (
	joyidCodeHash = "d23761b364210735c19c60561d213fb3beae2fd6172743719eff6920e020baac"
)

func FromPubkeyHash(pubkeyHash []byte, network types.Network) *address.Address {
	lockScript := &types.Script{
		CodeHash: types.HexToHash(joyidCodeHash),
		HashType: types.HashTypeType,
		Args:     pubkeyHash,
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
	return FromPubkeyHash(pubkeyHash, network)
}
