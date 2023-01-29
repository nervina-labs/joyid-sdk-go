package secp256r1

import (
	"testing"

	"github.com/nervina-labs/joyid-sdk-go/utils"
)

func TestGenerateKey(t *testing.T) {
	testcount := 10
	want := 32
	for testcount > 0 {
		gotkey, _ := GenerateKey()
		testcount = testcount - 1
		if got := len(gotkey.Bytes()); got != want {
			t.Errorf("GenerateKey length = %d, want %d", got, want)
		}
	}
}

func TestImportKey(t *testing.T) {
	testcases := []struct {
		key, wantkey string
	}{
		{"0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761", "0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"},
		{"0x2262cd6c965d0065f93fb1fce03444e7f2a354b215b16dc44fe88a7246b6213b", "0x2262cd6c965d0065f93fb1fce03444e7f2a354b215b16dc44fe88a7246b6213b"},
	}

	for _, tc := range testcases {
		gotkey := ImportKey(tc.key).Bytes()
		if got := utils.BytesTo0xHex(gotkey); got != tc.wantkey {
			t.Errorf("ImportKey() = %q, want %q", got, tc.wantkey)
		}
	}
}

func TestGetPubkey(t *testing.T) {
	key := ImportKey("0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761")
	_, pubkey := key.Pubkey()
	want := "0x4599a5795423d54ab8e1f44f5c6ef5be9b1829beddb787bc732e4469d25f8c93e94afa393617f905bf1765c35dc38501a862b4b2f794a88b4f9010da02411a85"
	if got := utils.BytesTo0xHex(pubkey); got != want {
		t.Errorf("GetPubkey() = %q, want %q", got, want)
	}
}

func TestGetPubkeyHash(t *testing.T) {
	testcases := []struct {
		key, wantpkhash string
	}{
		{"0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761", "0x6091d93dbab12f16640fb3a0a8f1e77e03fbc51c"},
		{"0x2262cd6c965d0065f93fb1fce03444e7f2a354b215b16dc44fe88a7246b6213b", "0x724acb22b1dead86b78f1375c3d176a4a2943653"},
	}

	for _, tc := range testcases {
		gotpkhash := ImportKey(tc.key).PubkeyHash()
		if got := utils.BytesTo0xHex(gotpkhash); got != tc.wantpkhash {
			t.Errorf("GetPubkeyHash() = %q, want %q", got, tc.wantpkhash)
		}
	}
}

func TestVerifySignature1(t *testing.T) {
	key := ImportKey("0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761")
	message, _ := utils.HexToBytes("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	got := key.VerifySignature(message)

	want := true
	if got != want {
		t.Errorf("VerifiSignature() = %t, want %t", got, want)
	}
}

func TestVerifySignature2(t *testing.T) {
	key, _ := GenerateKey()
	message, _ := utils.HexToBytes("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	got := key.VerifySignature(message)

	want := true
	if got != want {
		t.Errorf("VerifiSignature() = %t, want %t", got, want)
	}
}
