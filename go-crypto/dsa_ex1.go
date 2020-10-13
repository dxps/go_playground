// Example 1 of using Digital Signature Algorithm (DSA) package functions.
//
// It does the followings:
// 1. It generates a keypair of public and private keys
// 2. It signs a text (string) with the private key.
// 3. It verifies the signature with the public key.

package main

import (
	"crypto/dsa"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func main() {

	params := new(dsa.Parameters)

	if err := dsa.GenerateParameters(params, rand.Reader, dsa.L2048N256); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKey := new(dsa.PrivateKey)
	privateKey.PublicKey.Parameters = *params
	_ = dsa.GenerateKey(privateKey, rand.Reader) // it generates a public & private key pair

	publicKey := privateKey.PublicKey

	fmt.Printf(">>> Private key : %x\n", privateKey)
	fmt.Printf(">>> Public key  : %x\n", publicKey)

	h := md5.New()

	_, _ = io.WriteString(h, "The content to sign and verify")
	signHash := h.Sum(nil)

	r, s, err := dsa.Sign(rand.Reader, privateKey, signHash)
	if err != nil {
		fmt.Printf(">>> Error while signing: %s\n", err)
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	fmt.Printf(">>> Signature: %x\n", signature)

	verifyState := dsa.Verify(&publicKey, signHash, r, s)
	fmt.Printf(">>> Verify state: %t\n", verifyState)

	// Test that validation fails if changing the sign hash.
	_, _ = io.WriteString(h, "Additional content")
	signHash = h.Sum(nil)

	verifyState = dsa.Verify(&publicKey, signHash, r, s)
	fmt.Printf(">>> Verify state: %t\n", verifyState)
}
