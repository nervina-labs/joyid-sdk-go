package utils

import (
	"encoding/hex"
	"fmt"
)

func BytesToHex(bytes []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(bytes))
}
