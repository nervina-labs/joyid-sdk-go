package utils

import (
	"testing"
)

func TestBytesToHex(t *testing.T) {
	got := BytesToHex([]byte{204, 176, 131, 179, 122, 163, 70, 197, 206, 46, 31, 153, 166, 135, 161, 83, 186, 160, 64, 82, 242, 109, 182, 171, 60, 38, 214, 164, 204, 21, 197, 241})
	want := "ccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1"

	if got != want {
		t.Errorf("BytesToHex() = %s, want %s", got, want)
	}
}

func TestBytesTo0xHex(t *testing.T) {
	got := BytesTo0xHex([]byte{204, 176, 131, 179, 122, 163, 70, 197, 206, 46, 31, 153, 166, 135, 161, 83, 186, 160, 64, 82, 242, 109, 182, 171, 60, 38, 214, 164, 204, 21, 197, 241})
	want := "0xccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1"

	if got != want {
		t.Errorf("BytesTo0xHex() = %s, want %s", got, want)
	}
}

func TestHexToBytes(t *testing.T) {
	got1, _ := HexToBytes("0xccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")
	got2, _ := HexToBytes("ccb083b37aa346c5ce2e1f99a687a153baa04052f26db6ab3c26d6a4cc15c5f1")
	want := []byte{204, 176, 131, 179, 122, 163, 70, 197, 206, 46, 31, 153, 166, 135, 161, 83, 186, 160, 64, 82, 242, 109, 182, 171, 60, 38, 214, 164, 204, 21, 197, 241}
	if len(got1) != len(want) || len(got2) != len(want) || got1[0] != want[0] || got2[20] != want[20] {
		t.Errorf("HexToBytes() = %v, want %v", got1, want)
	}
}
