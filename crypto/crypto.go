package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/internal/fips140/sha3"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	secp256k1N     *big.Int
	secp256k1HalfN *big.Int
)

func init() {
	secp256k1N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1HalfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
}

// Keccak256Hash calculates the Keccak-256 hash of the input data
func Keccak256Hash(data []byte) common.Hash {
	return common.BytesToHash(Keccak256(data))
}

// Keccak256 calculates the Keccak-256 hash of the input data
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// GenerateKey generates a new ECDSA private key
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// Sign calculates an ECDSA signature
func Sign(hash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes")
	}
	seckey := math.PaddedBigBytes(prv.D, 32)
	defer zeroBytes(seckey)
	return secp256k1.Sign(hash, seckey)
}

// VerifySignature checks that the given public key created signature over hash
func VerifySignature(pubkey, hash, signature []byte) bool {
	return secp256k1.VerifySignature(pubkey, hash, signature)
}

// DecompressPubkey parses a public key in the 33-byte compressed format
func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	x, y := secp256k1.DecompressPubkey(pubkey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	return secp256k1.CompressPubkey(pubkey.X, pubkey.Y)
}

// PubkeyToAddress returns the Ethereum address of a public key
func PubkeyToAddress(pubkey *ecdsa.PublicKey) common.Address {
	pubBytes := FromECDSAPub(pubkey)
	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
}

// FromECDSAPub creates a public key from a raw ECDSA public key
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)
}

// ToECDSA creates a private key with the given D value
func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, true)
}

func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The private key must be non-zero and less than the order of the curve
	if priv.D.Sign() <= 0 || priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, X=%d", priv.D)
	}

	// The priv.D must not be 0 or >= N
	priv.PublicKey.X, priv.PublicKey.Y = elliptic.P256().ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return priv, nil
}

// FromECDSA exports a private key into a binary dump
func FromECDSA(prv *ecdsa.PrivateKey) []byte {
	if prv == nil {
		return nil
	}
	return math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
}

// HexToECDSA parses a secp256k1 private key
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToECDSA(b)
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b common.Address, nonce uint64) common.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return common.BytesToAddress(Keccak256(data)[12:])
}

// CreateAddress2 creates an ethereum address given the address, salt, and init code hash
func CreateAddress2(b common.Address, salt [32]byte, inithash []byte) common.Address {
	return common.BytesToAddress(Keccak256([]byte{0xff}, b.Bytes(), salt[:], inithash)[12:])
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 27.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}
	// reject upper range of s values (ECDSA malleability)
	// see discussion in https://github.com/ethereum/go-ethereum/pull/2053
	if homestead && s.Cmp(secp256k1HalfN) > 0 {
		return false
	}
	// Frontier: reject s value above N/2
	if v != 0 && v != 27 {
		return false
	}
	if r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 {
		return true
	}
	return false
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}

// S256 returns an instance of the secp256k1 curve
func S256() elliptic.Curve {
	return elliptic.P256()
}
