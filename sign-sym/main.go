package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func main() {
	msg := "hello, golang"
	key := "super-secret"
	fmt.Printf("message: %s\n", msg)
	fmt.Printf("key: %s\n", key)

	// sign
	signature := sign([]byte(key), []byte(msg))
	fmt.Printf("signature (HS256): %x\n", signature)

	// verify
	newSignature := sign([]byte(key), []byte(msg))
	fmt.Printf("verify result: %v\n", hmac.Equal(signature, newSignature))
}

func sign(key []byte, msg []byte) []byte {
	signer := hmac.New(sha256.New, []byte(key))
	signer.Write([]byte(msg))
	return signer.Sum(nil)
}
