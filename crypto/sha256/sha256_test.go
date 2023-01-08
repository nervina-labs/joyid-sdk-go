package sha256

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestSha256(t *testing.T) {
	msg := []byte("abc")
	exp, _ := hex.DecodeString("ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := Sha256(in); return h[:] }, msg, exp)
}

func checkhash(t *testing.T, name string, f func([]byte) []byte, msg, exp []byte) {
	sum := f(msg)
	if !bytes.Equal(exp, sum) {
		t.Fatalf("hash %s mismatch: want: %x have: %x", name, exp, sum)
	}
}
