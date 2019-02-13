package main

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type Token struct {
	Subject string `json:"subject"`
	Admin   bool   `json:"admin"`
}

func (t *Token) Sign(algor string, key []byte) (string, error) {
	r := encodeJWTPart(map[string]interface{}{
		"alg": algor,
		"typ": "JWT",
	})
	r += "." + encodeJWTPart(t)

	var sig []byte
	switch algor {
	case "none":
		// do nothing
	case "HS256":
		h := hmac.New(sha256.New, key)
		h.Write([]byte(r))
		sig = h.Sum(nil)
	case "RS256":
		priv, err := x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return "", err
		}
		h := crypto.SHA256.New()
		h.Write([]byte(r))
		sig, err = rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, h.Sum(nil))
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("algor '%s' not support", algor)
	}

	r += "." + base64.RawURLEncoding.EncodeToString(sig)
	return r, nil
}

func ParseToken(token string, key []byte) (*Token, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token")
	}

	var meta map[string]interface{}
	err := decodeJWTPart(parts[0], &meta)
	if err != nil {
		return nil, err
	}

	if meta["typ"] != "JWT" {
		return nil, fmt.Errorf("invalid token")
	}

	k := parts[0] + "." + parts[1]
	sig, _ := base64.RawURLEncoding.DecodeString(parts[2])

	switch meta["alg"] {
	case "none":
		if len(sig) != 0 {
			return nil, fmt.Errorf("invalid signature")
		}
	case "HS256":
		h := hmac.New(sha256.New, key)
		h.Write([]byte(k))
		if !hmac.Equal(sig, h.Sum(nil)) {
			return nil, fmt.Errorf("invalid signature")
		}
	case "RS256":
		pub, err := x509.ParsePKIXPublicKey(key)
		if err != nil {
			return nil, err
		}
		h := crypto.SHA256.New()
		h.Write([]byte(k))
		err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, h.Sum(nil), sig)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("algor '%s' not support", meta["alg"])
	}

	var t Token
	err = decodeJWTPart(parts[1], &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func encodeJWTPart(data interface{}) string {
	b, _ := json.Marshal(data)
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeJWTPart(data string, v interface{}) error {
	b, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
