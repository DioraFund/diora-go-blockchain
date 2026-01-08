package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type Wallet struct {
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	PrivateKey string    `json:"private_key"`
	PublicKey  string    `json:"public_key"`
	Balance    *big.Int  `json:"balance"`
	Nonce      uint64    `json:"nonce"`
	CreatedAt  time.Time `json:"created_at"`
}

type WalletManager struct {
	wallets map[string]*Wallet
	mutex   sync.RWMutex
}

var walletManager = &WalletManager{
	wallets: make(map[string]*Wallet),
}

func NewWallet(name string) (*Wallet, error) {
	walletManager.mutex.Lock()
	defer walletManager.mutex.Unlock()

	// Generate private key
	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	// Generate public key (simplified)
	publicKey := generatePublicKey(privateKey)

	// Generate address
	address := generateAddress(publicKey)

	wallet := &Wallet{
		Name:       name,
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Balance:    big.NewInt(0),
		Nonce:      0,
		CreatedAt:  time.Now(),
	}

	walletManager.wallets[address] = wallet
	return wallet, nil
}

func GetWallet(address string) (*Wallet, bool) {
	walletManager.mutex.RLock()
	defer walletManager.mutex.RUnlock()

	wallet, exists := walletManager.wallets[address]
	return wallet, exists
}

func GetAllWallets() []*Wallet {
	walletManager.mutex.RLock()
	defer walletManager.mutex.RUnlock()

	wallets := make([]*Wallet, 0, len(walletManager.wallets))
	for _, wallet := range walletManager.wallets {
		wallets = append(wallets, wallet)
	}
	return wallets
}

func UpdateBalance(address string, amount *big.Int) {
	walletManager.mutex.Lock()
	defer walletManager.mutex.Unlock()

	if wallet, exists := walletManager.wallets[address]; exists {
		wallet.Balance.Add(wallet.Balance, amount)
	}
}

func generatePrivateKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(bytes), nil
}

func generatePublicKey(privateKey string) string {
	// Simplified public key generation
	// In real implementation, this would use elliptic curve cryptography
	return "0x" + privateKey[4:] + "public"
}

func generateAddress(publicKey string) string {
	// Simplified address generation
	// In real implementation, this would hash the public key
	hash := publicKey[2:] // Remove 0x prefix
	if len(hash) > 40 {
		hash = hash[:40]
	} else {
		// Pad with zeros if too short
		for len(hash) < 40 {
			hash += "0"
		}
	}
	return "0x" + hash
}

func (w *Wallet) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"name":        w.Name,
		"address":     w.Address,
		"public_key":  w.PublicKey,
		"balance":     w.Balance.String(),
		"nonce":       w.Nonce,
		"created_at":  w.CreatedAt,
		"private_key": w.PrivateKey, // Only for demo - never expose private keys in production
	}
}

func (w *Wallet) Send(to string, amount *big.Int) error {
	// Check balance
	if w.Balance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance")
	}

	// Create transaction
	tx, err := NewTransaction(w.Address, to, amount.String())
	if err != nil {
		return err
	}

	// Add to transactions pool
	transactions = append(transactions, *tx)

	// Update balances
	w.Balance.Sub(w.Balance, amount)
	w.Nonce++

	// Add to block
	blockData := fmt.Sprintf("Transaction: %s -> %s (%s)", w.Address, to, amount.String())
	blockchain.AddBlock(blockData)

	// Update recipient balance
	UpdateBalance(to, amount)

	return nil
}
