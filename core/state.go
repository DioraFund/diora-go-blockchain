package core

import (
	"crypto/sha256"
	"math/big"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type State struct {
	db    *leveldb.DB
	trie  *Trie
	cache map[Address]*Account
	mu    sync.RWMutex
}

func NewState(db *leveldb.DB) *State {
	return &State{
		db:    db,
		trie:  NewTrie(db),
		cache: make(map[Address]*Account),
	}
}

func (s *State) GetAccount(addr Address) *Account {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if account, exists := s.cache[addr]; exists {
		return account
	}

	account := &Account{
		Nonce:   0,
		Balance: big.NewInt(0),
		Storage: make(map[Hash]Hash),
	}

	// Load from database
	data, err := s.db.Get(addr.Bytes(), nil)
	if err == nil {
		// Deserialize account (simplified)
		account.Balance.SetBytes(data)
	}

	s.cache[addr] = account
	return account
}

func (s *State) SetBalance(addr Address, balance *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account := s.GetAccount(addr)
	account.Balance.Set(balance)
}

func (s *State) GetBalance(addr Address) *big.Int {
	account := s.GetAccount(addr)
	return new(big.Int).Set(account.Balance)
}

func (s *State) SetNonce(addr Address, nonce uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account := s.GetAccount(addr)
	account.Nonce = nonce
}

func (s *State) GetNonce(addr Address) uint64 {
	account := s.GetAccount(addr)
	return account.Nonce
}

func (s *State) SetCode(addr Address, code []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account := s.GetAccount(addr)
	account.Code = code
	account.CodeHash = Keccak256Hash(code).Bytes()
}

func (s *State) GetCode(addr Address) []byte {
	account := s.GetAccount(addr)
	return account.Code
}

func (s *State) GetCodeHash(addr Address) []byte {
	account := s.GetAccount(addr)
	return account.CodeHash
}

func (s *State) SetState(addr Address, key, value Hash) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account := s.GetAccount(addr)
	account.Storage[key] = value
}

func (s *State) GetState(addr Address, key Hash) Hash {
	account := s.GetAccount(addr)
	if value, exists := account.Storage[key]; exists {
		return value
	}
	return Hash{}
}

func (s *State) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Save all cached accounts to database
	for addr, account := range s.cache {
		data := account.Balance.Bytes()
		if err := s.db.Put(addr.Bytes(), data, nil); err != nil {
			return err
		}

		// Save code if exists
		if len(account.Code) > 0 {
			codeHash := Keccak256(account.Code)
			if err := s.db.Put(codeHash, account.Code, nil); err != nil {
				return err
			}
		}

		// Save storage
		for key, value := range account.Storage {
			storageKey := append(addr.Bytes(), key.Bytes()...)
			if err := s.db.Put(storageKey, value.Bytes(), nil); err != nil {
				return err
			}
		}
	}

	// Clear cache
	s.cache = make(map[Address]*Account)
	return nil
}

func Keccak256Hash(data []byte) Hash {
	hasher := sha256.New()
	hasher.Write(data)
	return BytesToHash(hasher.Sum(nil))
}

func (s *State) Root() Hash {
	return Hash{}
}

func (s *State) Copy() *State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newState := &State{
		db:    s.db,
		trie:  s.trie,
		cache: make(map[Address]*Account),
	}

	// Copy cache
	for addr, account := range s.cache {
		newAccount := &Account{
			Nonce:    account.Nonce,
			Balance:  new(big.Int).Set(account.Balance),
			CodeHash: append([]byte(nil), account.CodeHash...),
			Code:     append([]byte(nil), account.Code...),
			Storage:  make(map[Hash]Hash),
		}

		for key, value := range account.Storage {
			newAccount.Storage[key] = value
		}

		newState.cache[addr] = newAccount
	}

	return newState
}
