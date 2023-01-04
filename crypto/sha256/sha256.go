package sha256

import (
	"crypto/sha256"
)

func Sha256(data []byte) []byte {
	d := sha256.New()
	d.Write(data)
	return d.Sum(nil)
}
