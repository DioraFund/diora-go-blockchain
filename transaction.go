package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

type Transaction struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    *big.Int  `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Signature string    `json:"signature"`
	Hash      string    `json:"hash"`
}

type TransactionRequest struct {
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

func NewTransaction(from, to, amount string) (*Transaction, error) {
	// Convert amount to big.Int
	amountInt, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount")
	}

	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amountInt,
		Timestamp: time.Now(),
	}

	// Calculate hash
	tx.Hash = tx.CalculateHash()

	// Sign transaction (simplified - in real blockchain this would use private key)
	tx.Signature = "SIGNATURE_" + tx.Hash

	return tx, nil
}

func (tx *Transaction) CalculateHash() string {
	data := fmt.Sprintf("%s%s%s%s",
		tx.From,
		tx.To,
		tx.Amount.String(),
		tx.Timestamp.String())

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (tx *Transaction) IsValid() bool {
	// Check if hash matches
	return tx.Hash == tx.CalculateHash()
}

func (tx *Transaction) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"from":      tx.From,
		"to":        tx.To,
		"amount":    tx.Amount.String(),
		"timestamp": tx.Timestamp,
		"signature": tx.Signature,
		"hash":      tx.Hash,
	}
}
