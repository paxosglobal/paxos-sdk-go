package paxos

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseSigningKey(key, name string) (crypto.Signer, error) {
	if strings.HasPrefix(strings.TrimSpace(key), "-----BEGIN") {
		return parseSigningKeyPEM(key)
	}
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid %s hex: expected ASCII hexadecimal characters", name)
	}
	if len(keyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("%s must be %d bytes (got %d)", name, ed25519.PrivateKeySize, len(keyBytes))
	}
	if err := validateSigningPublicKey(keyBytes); err != nil {
		return nil, err
	}
	return ed25519.PrivateKey(keyBytes), nil
}

func parseSigningKeyPEM(key string) (crypto.Signer, error) {
	block, rest := pem.Decode([]byte(strings.TrimSpace(key)))
	if block == nil {
		return nil, errors.New("invalid signing key PEM: malformed PEM block")
	}
	if len(bytes.TrimSpace(rest)) != 0 {
		return nil, errors.New("invalid signing key PEM: expected a single PEM block")
	}
	var parsed any
	var err error
	switch block.Type {
	case "PRIVATE KEY":
		parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		parsed, err = x509.ParseECPrivateKey(block.Bytes)
	case "ENCRYPTED PRIVATE KEY":
		return nil, errors.New("encrypted PEM signing keys are not supported; decrypt the key first")
	default:
		return nil, errors.New("unsupported PEM block type; expected an unencrypted PKCS#8 \"PRIVATE KEY\" or SEC1 \"EC PRIVATE KEY\" containing an Ed25519 (EdDSA) or P-256 (ES256) key")
	}
	if err != nil {
		return nil, errors.New("invalid signing key PEM: not a valid PKCS#8 or SEC1 private key; the SDK supports Ed25519 (EdDSA) and P-256 (ES256) keys only")
	}
	switch k := parsed.(type) {
	case ed25519.PrivateKey:
		return k, nil
	case *ecdsa.PrivateKey:
		if k.Curve == elliptic.P256() {
			if !hasFullLengthScalar(block) {
				return nil, errors.New("invalid signing key PEM: not a valid PKCS#8 or SEC1 private key; the SDK supports Ed25519 (EdDSA) and P-256 (ES256) keys only")
			}
			return k, nil
		}
	}
	return nil, errors.New("unsupported signing key type in PEM; the SDK supports Ed25519 (EdDSA) and P-256 (ES256) keys only")
}

type sec1PrivateKey struct {
	Version    int
	PrivateKey []byte
	Curve      asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey  asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

type pkcs8PrivateKey struct {
	Version    int
	Algorithm  asn1.RawValue
	PrivateKey []byte
	Attributes asn1.RawValue `asn1:"optional,tag:0"`
}

func hasFullLengthScalar(block *pem.Block) bool {
	der := block.Bytes
	if block.Type == "PRIVATE KEY" {
		var wrapper pkcs8PrivateKey
		if rest, err := asn1.Unmarshal(der, &wrapper); err != nil || len(rest) != 0 {
			return false
		}
		der = wrapper.PrivateKey
	}
	var inner sec1PrivateKey
	if rest, err := asn1.Unmarshal(der, &inner); err != nil || len(rest) != 0 {
		return false
	}
	return len(inner.PrivateKey) == 32
}

// validateSigningPublicKey rejects inconsistent seed+public-key pairs before
// Go's Ed25519 implementation can produce signatures with the wrong public half.
// Callers must first check the 64-byte key length.
func validateSigningPublicKey(key []byte) error {
	derived := ed25519.NewKeyFromSeed(key[:ed25519.SeedSize])
	if !bytes.Equal(derived[ed25519.SeedSize:], key[ed25519.SeedSize:]) {
		return fmt.Errorf("signing key public half does not match its seed")
	}
	return nil
}

// signRequestWithKey creates a JWS using an already-parsed Ed25519 private key.
func signRequestWithKey(privateKey crypto.Signer, keyID, method, path string, body []byte) (string, error) {
	return signRequestWithKeyAndTime(privateKey, keyID, method, path, body, time.Now())
}

// signRequestWithKeyAndTime creates a JWS using an already-parsed Ed25519
// private key and a specific timestamp. This is primarily used for testing.
func signRequestWithKeyAndTime(privateKey crypto.Signer, keyID, method, path string, body []byte, now time.Time) (string, error) {
	if body == nil {
		body = []byte{}
	}

	// Build the protected header with explicit key order matching all SDK languages.
	alg := "EdDSA"
	if _, ok := privateKey.(*ecdsa.PrivateKey); ok {
		alg = "ES256"
	}
	timestamp := strconv.FormatInt(now.Unix(), 10)
	headerJSON := []byte(`{"alg":` + jsonString(alg) + `,"kid":` + jsonString(keyID) +
		`,"paxos.com/request-method":` + jsonString(method) +
		`,"paxos.com/request-path":` + jsonString(path) +
		`,"paxos.com/timestamp":` + jsonString(timestamp) + `}`)

	// Base64url encode header and payload (no padding).
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(body)

	// Sign base64url(header).base64url(payload).
	signingInput := encodedHeader + "." + encodedPayload
	var signature []byte
	switch k := privateKey.(type) {
	case *ecdsa.PrivateKey:
		digest := sha256.Sum256([]byte(signingInput))
		r, sv, err := ecdsa.Sign(rand.Reader, k, digest[:])
		if err != nil {
			return "", fmt.Errorf("signing failed: %w", err)
		}
		signature = make([]byte, 64)
		r.FillBytes(signature[:32])
		sv.FillBytes(signature[32:])
	case ed25519.PrivateKey:
		signature = ed25519.Sign(k, []byte(signingInput))
	default:
		return "", errors.New("unsupported signing key type")
	}

	// Base64url encode the signature.
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	// Return JWS compact serialization.
	return encodedHeader + "." + encodedPayload + "." + encodedSignature, nil
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
