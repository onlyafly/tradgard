package extension

import (
	"fmt"
	"os"

	"github.com/gorilla/securecookie"
)

// NewSecureCookie creates a new secureCookie object for encoding them
func NewSecureCookie() *securecookie.SecureCookie {
	// Hash keys should be at least 32 bytes long
	var hashKey = []byte(os.Getenv("COOKIE_32_BYTE_HASH_KEY"))

	fmt.Println(os.Getenv("COOKIE_32_BYTE_HASH_KEY"))
	fmt.Println(os.Getenv("COOKIE_32_BYTE_BLOCK_KEY"))

	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = []byte(os.Getenv("COOKIE_32_BYTE_BLOCK_KEY"))
	var secureCookie = securecookie.New(hashKey, blockKey)
	return secureCookie
}
