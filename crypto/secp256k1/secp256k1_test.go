package secp256k1

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
		{"0xccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1", "0xccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1"},
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
	testcases := []struct {
		key, wantpk string
	}{
		{"0xccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1", "0xa0a7a7597b019828a1dda6ed52ab25181073ec3a9825d28b9abbb932fe1ec83dd117a8eef7649c25be5a591d08f80ffe7e9c14100ad1b58ac78afa606a576453"},
		{"0x2262cd6c965d0065f93fb1fce03444e7f2a354b215b16dc44fe88a7246b6213b", "0x0009455d20f00e6a944017377122412b927c23e85bd4da670ac619217a9de67b393c74fd78be4a30cec529505f11408cdc42a81a9bf8f08584b6c39cb9fc1783"},
	}

	for _, tc := range testcases {
		_, gotpk := ImportKey(tc.key).Pubkey()
		if got := utils.BytesTo0xHex(gotpk); got != tc.wantpk {
			t.Errorf("GetPubkey() = %q, want %q", got, tc.wantpk)
		}
	}
}

func TestGetPubkeyHash(t *testing.T) {
	testcases := []struct {
		key, wantpkhash string
	}{
		{"0xccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1", "0x426e2a8a44eac6d1573befaf9c2ba1ab17c683de"},
		{"0x2262cd6c965d0065f93fb1fce03444e7f2a354b215b16dc44fe88a7246b6213b", "0x4e580ac0fdc70417042ce5634b7b1f62fa52e170"},
	}

	for _, tc := range testcases {
		gotpkhash := ImportKey(tc.key).PubkeyHash()
		if got := utils.BytesTo0xHex(gotpkhash); got != tc.wantpkhash {
			t.Errorf("GetPubkeyHash() = %q, want %q", got, tc.wantpkhash)
		}
	}
}

func TestVerifySignature(t *testing.T) {
	key := ImportKey("0xccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")
	message, _ := utils.HexToBytes("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	got := utils.BytesTo0xHex(key.Sign(message))

	want := "0x692dc94fdaf9d9dded7cad66755da9cb79ec918f7bb69b4939a9ce1ac41c6589750d48ace8bc766531312c4ee36d9ec2a94921adb9f391ddde47a44baae6e8f000"
	if got != want {
		t.Errorf("Sign() = %s, want %s", got, want)
	}
}

func TestRecoverPubkey1(t *testing.T) {
	key, _ := GenerateKey()
	_, pubkey := key.Pubkey()
	message, _ := utils.HexToBytes("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	sig := key.Sign(message)
	recoveryPubkey := key.RecoverPubkey(message, sig)
	got := utils.BytesToHex(pubkey)

	want := utils.BytesToHex(recoveryPubkey[1:])
	if got != want {
		t.Errorf("RecoverPubkey() = %s, want %s", got, want)
	}
}

func TestRecoverPubkey2(t *testing.T) {
	key := ImportKey("0x2262cd6c965d0065f93fb1fce03444e7f2a354b215b16dc44fe88a7246b6213b")
	_, pubkey := key.Pubkey()
	got := utils.BytesToHex(pubkey)

	message, _ := utils.HexToBytes("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	sig, _ := utils.HexToBytes("0x5e999cf4bf6798a154204d03aca6de2fec11e7b8367c18d9d277b09f4617728f74c6c1b8b15ca8407bc4b51f4e38bb5f1d5434d61ec1c1045b9b1e7cf4db533500")
	recoveryPubkey := key.RecoverPubkey(message, sig)

	want := utils.BytesToHex(recoveryPubkey[1:])
	if got != want {
		t.Errorf("RecoverPubkey() = %s, want %s", got, want)
	}
}
