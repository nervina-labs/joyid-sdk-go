package secp256k1

import (
	"testing"
)

func TestGetPubkey(t *testing.T) {
	key := ImportKey("ccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")
	_, pubkey := key.Pubkey()
	want := "a0a7a7597b019828a1dda6ed52ab25181073ec3a9825d28b9abbb932fe1ec83dd117a8eef7649c25be5a591d08f80ffe7e9c14100ad1b58ac78afa606a576453"
	if got := pubkey; got != want {
		t.Errorf("GetPubkey() = %q, want %q", got, want)
	}
}

func TestVerifiSignature1(t *testing.T) {
	key := ImportKey("ccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")
	got := key.VerifiSignature("acba4329945ecb0e4f1db924e48a7ab27db75f36346f6b2b88e70d49a9cadeb2")

	want := true
	if got != want {
		t.Errorf("VerifiSignature() = %t, want %t", got, want)
	}
}

func TestVerifiSignature2(t *testing.T) {
	key, _ := GenerateKey()
	got := key.VerifiSignature("ccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")

	want := true
	if got != want {
		t.Errorf("VerifiSignature() = %t, want %t", got, want)
	}
}
