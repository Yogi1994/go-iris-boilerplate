package utils

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//RandASCIIBytes - A helper function create and fill a slice of length n with characters from a-zA-Z0-9_-. It panics if there are any problems getting random bytes.
func RandASCIIBytes(n int) string {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(letterBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = letterBytes[randomPos]
	}

	return string(output)
}
