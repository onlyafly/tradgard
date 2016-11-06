package extension

import "github.com/gorilla/securecookie"

//FIXME

// NewSecureCookie creates a new secureCookie object for encoding them
func NewSecureCookie() *securecookie.SecureCookie {
	// Hash keys should be at least 32 bytes long
	var hashKey = []byte("iMXwykcZbaf7rPJNMxSSvGbd20uSq0Mi")
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = []byte("te2Ltr1dMH0UlyOc")
	var secureCookie = securecookie.New(hashKey, blockKey)
	return secureCookie
}
