package next_utils

import "crypto/sha256"

func Sha256(b []byte) []byte {
	b32 := sha256.Sum256(b)
	return b32[:]
}

func DoubleSha256(b []byte) []byte {
	return Sha256(Sha256(b))
}
