package secp256r1

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestGetPubkey(t *testing.T) {
	key := ImportKey("4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761")
	_, pubkey := key.Pubkey()
	want := "4599a5795423d54ab8e1f44f5c6ef5be9b1829beddb787bc732e4469d25f8c93e94afa393617f905bf1765c35dc38501a862b4b2f794a88b4f9010da02411a85"
	if got := fmt.Sprintf("%x", pubkey); got != want {
		t.Errorf("GetPubkey() = %q, want %q", got, want)
	}
}

func TestVerifySignature1(t *testing.T) {
	key := ImportKey("4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761")
	message, _ := hexutil.Decode("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	got := key.VerifySignature(message)

	want := true
	if got != want {
		t.Errorf("VerifiSignature() = %t, want %t", got, want)
	}
}

func TestVerifySignature2(t *testing.T) {
	key, _ := GenerateKey()
	message, _ := hexutil.Decode("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	got := key.VerifySignature(message)

	want := true
	if got != want {
		t.Errorf("VerifiSignature() = %t, want %t", got, want)
	}
}
