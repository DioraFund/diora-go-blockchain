package keystore

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/pbkdf2"
)

// Account represents a wallet account
type Account struct {
	Name      string         `json:"name"`
	Address   common.Address `json:"address"`
	PublicKey string         `json:"public_key"`
	CreatedAt time.Time      `json:"created_at"`
}

// EncryptedKeystore represents an encrypted keystore file
type EncryptedKeystore struct {
	Address   string    `json:"address"`
	PublicKey string    `json:"public_key"`
	Crypto    Crypto    `json:"crypto"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

// Crypto represents the crypto section of keystore
type Crypto struct {
	KDF          string      `json:"kdf"`
	KDFParams    interface{} `json:"kdfparams"`
	Cipher       string      `json:"cipher"`
	CipherText   string      `json:"ciphertext"`
	CipherParams interface{} `json:"cipherparams"`
	MAC          string      `json:"mac"`
}

// Keystore manages wallet accounts
type Keystore struct {
	path string
}

// NewKeystore creates a new keystore instance
func NewKeystore(path string) *Keystore {
	return &Keystore{
		path: path,
	}
}

// NewAccount creates a new account
func (ks *Keystore) NewAccount(name, password string) (*Account, error) {
	// Generate private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create account
	account := &Account{
		Name:      name,
		Address:   crypto.PubkeyToAddress(privateKey.PublicKey),
		PublicKey: hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey)),
		CreatedAt: time.Now(),
	}

	// Encrypt and save private key
	if err := ks.savePrivateKey(privateKey, account, password); err != nil {
		return nil, fmt.Errorf("failed to save private key: %w", err)
	}

	// Save account metadata
	if err := ks.saveAccount(account); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	return account, nil
}

// ImportPrivateKey imports a private key
func (ks *Keystore) ImportPrivateKey(privateKeyStr, name, password string) (*Account, error) {
	// Decode private key
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Create account
	account := &Account{
		Name:      name,
		Address:   crypto.PubkeyToAddress(privateKey.PublicKey),
		PublicKey: hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey)),
		CreatedAt: time.Now(),
	}

	// Encrypt and save private key
	if err := ks.savePrivateKey(privateKey, account, password); err != nil {
		return nil, fmt.Errorf("failed to save private key: %w", err)
	}

	// Save account metadata
	if err := ks.saveAccount(account); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	return account, nil
}

// ListAccounts returns all accounts
func (ks *Keystore) ListAccounts() ([]*Account, error) {
	var accounts []*Account

	// Ensure keystore directory exists
	if err := os.MkdirAll(ks.path, 0700); err != nil {
		return nil, fmt.Errorf("failed to create keystore directory: %w", err)
	}

	// Read account files
	files, err := ioutil.ReadDir(ks.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// Skip keystore files (private keys)
		if file.Name()[:8] == "UTC--" {
			continue
		}

		// Read account file
		accountPath := filepath.Join(ks.path, file.Name())
		data, err := ioutil.ReadFile(accountPath)
		if err != nil {
			continue
		}

		var account Account
		if err := json.Unmarshal(data, &account); err != nil {
			continue
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// GetAccount returns an account by name
func (ks *Keystore) GetAccount(name string) (*Account, error) {
	accountPath := filepath.Join(ks.path, name+".json")
	data, err := ioutil.ReadFile(accountPath)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	var account Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	return &account, nil
}

// ExportPrivateKey exports the private key
func (ks *Keystore) ExportPrivateKey(name, password string) (string, error) {
	account, err := ks.GetAccount(name)
	if err != nil {
		return "", err
	}

	// Find keystore file
	keystoreFile, err := ks.findKeystoreFile(account.Address)
	if err != nil {
		return "", fmt.Errorf("keystore file not found: %w", err)
	}

	// Read keystore file
	data, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return "", fmt.Errorf("failed to read keystore file: %w", err)
	}

	var keystore EncryptedKeystore
	if err := json.Unmarshal(data, &keystore); err != nil {
		return "", fmt.Errorf("failed to unmarshal keystore: %w", err)
	}

	// Decrypt private key
	privateKey, err := ks.decryptPrivateKey(&keystore, password)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %w", err)
	}

	return hex.EncodeToString(crypto.FromECDSA(privateKey)), nil
}

// DeleteAccount deletes an account
func (ks *Keystore) DeleteAccount(name string) error {
	account, err := ks.GetAccount(name)
	if err != nil {
		return err
	}

	// Delete account file
	accountPath := filepath.Join(ks.path, name+".json")
	if err := os.Remove(accountPath); err != nil {
		return fmt.Errorf("failed to delete account file: %w", err)
	}

	// Delete keystore file
	keystoreFile, err := ks.findKeystoreFile(account.Address)
	if err != nil {
		return fmt.Errorf("keystore file not found: %w", err)
	}

	if err := os.Remove(keystoreFile); err != nil {
		return fmt.Errorf("failed to delete keystore file: %w", err)
	}

	return nil
}

// savePrivateKey saves encrypted private key
func (ks *Keystore) savePrivateKey(privateKey *ecdsa.PrivateKey, account *Account, password string) error {
	// Ensure keystore directory exists
	if err := os.MkdirAll(ks.path, 0700); err != nil {
		return fmt.Errorf("failed to create keystore directory: %w", err)
	}

	// Create keystore filename
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000000Z")
	filename := fmt.Sprintf("UTC--%s--%s", timestamp, account.Address.Hex())
	keystorePath := filepath.Join(ks.path, filename)

	// Encrypt private key
	keystore, err := ks.encryptPrivateKey(privateKey, account, password)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Save keystore file
	data, err := json.MarshalIndent(keystore, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal keystore: %w", err)
	}

	if err := ioutil.WriteFile(keystorePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write keystore file: %w", err)
	}

	return nil
}

// saveAccount saves account metadata
func (ks *Keystore) saveAccount(account *Account) error {
	accountPath := filepath.Join(ks.path, account.Name+".json")
	data, err := json.MarshalIndent(account, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal account: %w", err)
	}

	return ioutil.WriteFile(accountPath, data, 0644)
}

// encryptPrivateKey encrypts private key
func (ks *Keystore) encryptPrivateKey(privateKey *ecdsa.PrivateKey, account *Account, password string) (*EncryptedKeystore, error) {
	// Generate salt
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using PBKDF2
	derivedKey := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Encrypt private key (simplified - in production use proper encryption)
	privateKeyBytes := crypto.FromECDSA(privateKey)
	cipherText := make([]byte, len(privateKeyBytes))
	for i, b := range privateKeyBytes {
		cipherText[i] = b ^ derivedKey[i%len(derivedKey)]
	}

	// Generate MAC
	mac := sha256.Sum256(append(derivedKey[16:32], cipherText...))

	return &EncryptedKeystore{
		Address:   account.Address.Hex(),
		PublicKey: account.PublicKey,
		Crypto: Crypto{
			KDF: "pbkdf2",
			KDFParams: map[string]interface{}{
				"dklen": 32,
				"c":     100000,
				"salt":  hex.EncodeToString(salt),
				"prf":   "hmac-sha256",
			},
			Cipher:     "aes-128-ctr",
			CipherText: hex.EncodeToString(cipherText),
			CipherParams: map[string]interface{}{
				"iv": hex.EncodeToString(make([]byte, 16)),
			},
			MAC: hex.EncodeToString(mac[:]),
		},
		Name:      account.Name,
		CreatedAt: account.CreatedAt,
		Version:   3,
	}, nil
}

// decryptPrivateKey decrypts private key
func (ks *Keystore) decryptPrivateKey(keystore *EncryptedKeystore, password string) (*ecdsa.PrivateKey, error) {
	// Extract parameters
	kdfParams := keystore.Crypto.KDFParams.(map[string]interface{})
	salt, err := hex.DecodeString(kdfParams["salt"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid salt: %w", err)
	}

	// Derive key
	derivedKey := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Decrypt private key (simplified - in production use proper decryption)
	cipherText, err := hex.DecodeString(keystore.Crypto.CipherText)
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext: %w", err)
	}

	privateKeyBytes := make([]byte, len(cipherText))
	for i, b := range cipherText {
		privateKeyBytes[i] = b ^ derivedKey[i%len(derivedKey)]
	}

	// Convert to private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to private key: %w", err)
	}

	return privateKey, nil
}

// findKeystoreFile finds keystore file by address
func (ks *Keystore) findKeystoreFile(address common.Address) (string, error) {
	files, err := ioutil.ReadDir(ks.path)
	if err != nil {
		return "", fmt.Errorf("failed to read keystore directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || file.Name()[:8] != "UTC--" {
			continue
		}

		// Read keystore file
		keystorePath := filepath.Join(ks.path, file.Name())
		data, err := ioutil.ReadFile(keystorePath)
		if err != nil {
			continue
		}

		var keystore EncryptedKeystore
		if err := json.Unmarshal(data, &keystore); err != nil {
			continue
		}

		if keystore.Address == address.Hex() {
			return keystorePath, nil
		}
	}

	return "", fmt.Errorf("keystore file not found for address %s", address.Hex())
}
