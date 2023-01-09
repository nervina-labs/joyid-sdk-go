package address

import (
	"testing"

	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/utils"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

func TestFromR1PrivKey(t *testing.T) {
	address, _ := FromPrivKey("4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761", alg.Secp256r1, types.NetworkTest).Encode()
	want := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqq9sfrkfah2cj79nyp7e6p283ualq8779rsgww3jf"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}

func TestFromK1PrivKey(t *testing.T) {
	address, _ := FromPrivKey("4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761", alg.Secp256k1, types.NetworkTest).Encode()
	want := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvmef595"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}

func TestFromR1PubkeyHash(t *testing.T) {
	pubkeyHash, _ := utils.HexToBytes("0x6091d93dbab12f16640fb3a0a8f1e77e03fbc51c")
	address, _ := FromPubkeyHash(pubkeyHash, alg.Secp256r1, types.NetworkTest).Encode()
	want := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqq9sfrkfah2cj79nyp7e6p283ualq8779rsgww3jf"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}

func TestFromK1PubkeyHash(t *testing.T) {
	pubkeyHash, _ := utils.HexToBytes("0x6500fc0e86fd49ef7dfc4b25dfd654eacaad53fb")
	address, _ := FromPubkeyHash(pubkeyHash, alg.Secp256k1, types.NetworkTest).Encode()
	want := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvmef595"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}
