package internal

import (
	"encoding/binary"
	"math/rand/v2"
)

const defaultRandomAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// RandomString generates a cryptographically random string with the specified length.
//
// The generated string matches [A-Za-z0-9]+ and it's transparent to URL-encoding.
func RandomString(length int) string {
	return RandomStringWithAlphabet(length, defaultRandomAlphabet)
}

// RandomStringWithAlphabet generates a cryptographically random string
// with the specified length and characters set.
func RandomStringWithAlphabet(length int, alphabet string) string {
	b := make([]byte, length)
	max := len(alphabet)
	for i := range b {
		b[i] = alphabet[rand.IntN(max)]
	}
	return string(b)
}

func NewChaCha8() *rand.ChaCha8 {
	var b [32]byte
	binary.LittleEndian.PutUint64(b[:], rand.Uint64())
	binary.LittleEndian.AppendUint64(b[:], rand.Uint64())
	binary.LittleEndian.AppendUint64(b[:], rand.Uint64())
	binary.LittleEndian.AppendUint64(b[:], rand.Uint64())
	return rand.NewChaCha8(b)
}
