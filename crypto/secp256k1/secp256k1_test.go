package secp256k1

import (
	"fmt"
	"testing"

	"github.com/nervina-labs/joyid-sdk-go/utils"
)

func TestGetPubkey(t *testing.T) {
	key := ImportKey("ccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")
	_, pubkey := key.Pubkey()
	want := "a0a7a7597b019828a1dda6ed52ab25181073ec3a9825d28b9abbb932fe1ec83dd117a8eef7649c25be5a591d08f80ffe7e9c14100ad1b58ac78afa606a576453"
	if got := fmt.Sprintf("%x", pubkey); got != want {
		t.Errorf("GetPubkey() = %q, want %q", got, want)
	}
}

func TestVerifySignature1(t *testing.T) {
	key := ImportKey("ccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")
	message, _ := utils.HexToBytes("0xacba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")
	got := utils.BytesTo0xHex(key.Sign(message))

	want := "0x692dc94fdaf9d9dded7cad66755da9cb79ec918f7bb69b4939a9ce1ac41c6589750d48ace8bc766531312c4ee36d9ec2a94921adb9f391ddde47a44baae6e8f000"
	if got != want {
		t.Errorf("VerifiSignature() = %s, want %s", got, want)
	}
}
