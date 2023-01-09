package address

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

func TestFromR1PrivKey(t *testing.T) {
	address, _ := FromPrivKey("4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761", alg.Secp256r1, types.NetworkTest).Encode()
	want := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqq9sfrkfah2cj79nyp7e6p283ualq8779rscnjmrj"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}

func TestFromK1PrivKey(t *testing.T) {
	address, _ := FromPrivKey("4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761", alg.Secp256k1, types.NetworkTest).Encode()
	want := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvty4750"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}

func TestFromR1PubkeyHash(t *testing.T) {
	pubkeyHash, _ := hexutil.Decode("0x6091d93dbab12f16640fb3a0a8f1e77e03fbc51c")
	address, _ := FromPubkeyHash(pubkeyHash, alg.Secp256r1, types.NetworkTest).Encode()
	want := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqq9sfrkfah2cj79nyp7e6p283ualq8779rscnjmrj"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}

func TestFromK1PubkeyHash(t *testing.T) {
	pubkeyHash, _ := hexutil.Decode("0x6500fc0e86fd49ef7dfc4b25dfd654eacaad53fb")
	address, _ := FromPubkeyHash(pubkeyHash, alg.Secp256k1, types.NetworkTest).Encode()
	want := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvty4750"
	if got := address; got != want {
		t.Errorf("FromPrivKey() = %q, want %q", got, want)
	}
}
