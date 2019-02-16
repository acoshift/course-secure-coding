package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

func main() {
	pass := "supersecret"
	invalid := "keyboardcat"
	fmt.Println("password:", pass)

	{
		h := newBcryptHasher(12)
		hashed, _ := h.Hash(pass)
		fmt.Println("bcrypt:", hashed)
		fmt.Println("       test-ok:", boolToPass(h.Compare(hashed, pass)))
		fmt.Println("  test-invalid:", boolToPass(!h.Compare(hashed, invalid)))
	}

	{
		h := newBcryptSHA256Hasher(12)
		hashed, _ := h.Hash(pass)
		fmt.Println("bcrypt-sha256:", hashed)
		fmt.Println("       test-ok:", boolToPass(h.Compare(hashed, pass)))
		fmt.Println("  test-invalid:", boolToPass(!h.Compare(hashed, invalid)))
	}

	{
		h := newScryptHasher(0, 0, 0, 0, 0)
		hashed, _ := h.Hash(pass)
		fmt.Println("scrypt:", hashed)
		fmt.Println("       test-ok:", boolToPass(h.Compare(hashed, pass)))
		fmt.Println("  test-invalid:", boolToPass(!h.Compare(hashed, invalid)))
	}
}

func boolToPass(b bool) string {
	if b {
		return "passed"
	}
	return "not passed"
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, password string) bool
}

type bcryptHasher struct {
	c int
}

func newBcryptHasher(cost int) PasswordHasher {
	return &bcryptHasher{cost}
}

func (h bcryptHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), h.c)
	return string(hashed), err
}

func (bcryptHasher) Compare(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

type bcryptSHA256Hasher struct {
	bcryptHasher
}

func newBcryptSHA256Hasher(cost int) PasswordHasher {
	return &bcryptSHA256Hasher{bcryptHasher{cost}}
}

func (h bcryptSHA256Hasher) Hash(password string) (string, error) {
	s := sha256.Sum256([]byte(password))
	return h.bcryptHasher.Hash(string(s[:]))
}

func (h bcryptSHA256Hasher) Compare(hashed, password string) bool {
	s := sha256.Sum256([]byte(password))
	return h.bcryptHasher.Compare(hashed, string(s[:]))
}

type scryptHasher struct {
	n, r, p, s, k int
}

func newScryptHasher(n, r, p, s, k int) PasswordHasher {
	if n <= 0 {
		n = 32768
	}
	if r <= 0 {
		r = 8
	}
	if p <= 0 {
		p = 1
	}
	if s <= 0 {
		s = 16
	}
	if k <= 0 {
		k = 32
	}
	return &scryptHasher{n, r, p, s, k}
}

func (h scryptHasher) generateSalt() ([]byte, error) {
	p := make([]byte, h.s)
	_, err := rand.Read(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (h scryptHasher) Hash(password string) (string, error) {
	salt, err := h.generateSalt()
	if err != nil {
		return "", err
	}
	dk, err := scrypt.Key([]byte(password), salt, h.n, h.r, h.p, h.k)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d$%d$%d$%s$%s", h.n, h.r, h.p, h.encodeBase64(salt), h.encodeBase64(dk)), nil
}

func (h scryptHasher) Compare(hashed, password string) bool {
	n, r, p, salt, dk := h.decode(hashed)
	if len(dk) == 0 {
		return false
	}

	pk, err := scrypt.Key([]byte(password), salt, n, r, p, len(dk))
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(dk, pk) == 1
}

func (h scryptHasher) decode(hashed string) (n, r, p int, salt, dk []byte) {
	xs := strings.Split(hashed, "$")
	if len(xs) != 5 {
		return
	}

	var err error
	n, err = strconv.Atoi(xs[0])
	if err != nil {
		return
	}
	r, err = strconv.Atoi(xs[1])
	if err != nil {
		return
	}
	p, err = strconv.Atoi(xs[2])
	if err != nil {
		return
	}
	salt = h.decodeBase64(xs[3])
	dk = h.decodeBase64(xs[4])
	return
}

func (h scryptHasher) encodeBase64(p []byte) string {
	return base64.RawStdEncoding.EncodeToString(p)
}

func (h scryptHasher) decodeBase64(s string) []byte {
	p, _ := base64.RawStdEncoding.DecodeString(s)
	return p
}
