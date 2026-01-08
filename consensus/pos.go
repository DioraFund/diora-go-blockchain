package consensus

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/diora-blockchain/diora/core"
	"github.com/diora-blockchain/diora/crypto"
)

type PoS struct {
	config      *Config
	validators  map[core.Address]*Validator
	delegators  map[core.Address][]*Delegation
	stakeAmount *big.Int
	totalStake  *big.Int
	currentSlot uint64
	validator   core.Address
	privateKey  *ecdsa.PrivateKey
	mu          sync.RWMutex
}

type Config struct {
	MinStakeAmount    *big.Int
	MaxValidators     int
	UnbondingPeriod  time.Duration
	SlotDuration      time.Duration
	EpochLength       uint64
	RewardRate        *big.Int
	SlashRate         *big.Int
	CommissionRate    uint64
}

type Validator struct {
	Address       core.Address
	Stake         *big.Int
	TotalDelegated *big.Int
	Commission    uint64
	Status        ValidatorStatus
	LastActive    time.Time
	TotalBlocks   uint64
	Rewards       *big.Int
	PublicKey     *ecdsa.PublicKey
}

type Delegation struct {
	Delegator   core.Address
	Validator   core.Address
	Amount      *big.Int
	Rewards     *big.Int
	StartTime   time.Time
	UnbondTime  time.Time
	Status      DelegationStatus
}

type ValidatorStatus int

const (
	ValidatorStatusInactive ValidatorStatus = iota
	ValidatorStatusActive
	ValidatorStatusSlashed
)

type DelegationStatus int

const (
	DelegationStatusActive DelegationStatus = iota
	DelegationStatusUnbonding
	DelegationStatusCompleted
)

func NewPoS(stakeAmount *big.Int, maxValidators int) *PoS {
	return &PoS{
		stakeAmount: stakeAmount,
		validators:  make(map[core.Address]*Validator),
		delegators:  make(map[core.Address][]*Delegation),
		totalStake:  big.NewInt(0),
		config: &Config{
			MinStakeAmount:   stakeAmount,
			MaxValidators:    maxValidators,
			UnbondingPeriod:  7 * 24 * time.Hour, // 7 days
			SlotDuration:     6 * time.Second,
			EpochLength:      100,
			RewardRate:       big.NewInt(1000000000000000000), // 1 DIO per epoch
			SlashRate:        big.NewInt(500000000000000000),  // 0.5 DIO
			CommissionRate:   1000, // 10%
		},
	}
}

func (pos *PoS) SetValidator(address core.Address, privateKey *ecdsa.PrivateKey) {
	pos.mu.Lock()
	defer pos.mu.Unlock()
	
	pos.validator = address
	pos.privateKey = privateKey
}

func (pos *PoS) IsValidator() bool {
	pos.mu.RLock()
	defer pos.mu.RUnlock()
	
	return pos.validator != (core.Address{}) && pos.isCurrentValidator(pos.validator)
}

func (pos *PoS) GetValidatorAddress() core.Address {
	pos.mu.RLock()
	defer pos.mu.RUnlock()
	
	return pos.validator
}

func (pos *PoS) RegisterValidator(address core.Address, stake *big.Int, commission uint64, publicKey *ecdsa.PublicKey) error {
	pos.mu.Lock()
	defer pos.mu.Unlock()

	if stake.Cmp(pos.config.MinStakeAmount) < 0 {
		return fmt.Errorf("stake amount too low: minimum %s", pos.config.MinStakeAmount.String())
	}

	if len(pos.validators) >= pos.config.MaxValidators {
		return fmt.Errorf("maximum validators reached: %d", pos.config.MaxValidators)
	}

	if _, exists := pos.validators[address]; exists {
		return fmt.Errorf("validator already registered")
	}

	validator := &Validator{
		Address:       address,
		Stake:         new(big.Int).Set(stake),
		TotalDelegated: big.NewInt(0),
		Commission:    commission,
		Status:        ValidatorStatusActive,
		LastActive:    time.Now(),
		TotalBlocks:   0,
		Rewards:       big.NewInt(0),
		PublicKey:     publicKey,
	}

	pos.validators[address] = validator
	pos.totalStake.Add(pos.totalStake, stake)

	return nil
}

func (pos *PoS) Delegate(delegator, validator core.Address, amount *big.Int) error {
	pos.mu.Lock()
	defer pos.mu.Unlock()

	v, exists := pos.validators[validator]
	if !exists {
		return fmt.Errorf("validator not found")
	}

	if v.Status != ValidatorStatusActive {
		return fmt.Errorf("validator not active")
	}

	delegation := &Delegation{
		Delegator:  delegator,
		Validator:  validator,
		Amount:     new(big.Int).Set(amount),
		Rewards:    big.NewInt(0),
		StartTime:  time.Now(),
		UnbondTime: time.Time{},
		Status:     DelegationStatusActive,
	}

	pos.delegators[delegator] = append(pos.delegators[delegator], delegation)
	v.TotalDelegated.Add(v.TotalDelegated, amount)
	pos.totalStake.Add(pos.totalStake, amount)

	return nil
}

func (pos *PoS) Unbond(delegator, validator core.Address, amount *big.Int) error {
	pos.mu.Lock()
	defer pos.mu.Unlock()

	delegations, exists := pos.delegators[delegator]
	if !exists {
		return fmt.Errorf("no delegations found")
	}

	var targetDelegation *Delegation
	var unbondAmount *big.Int

	for _, delegation := range delegations {
		if delegation.Validator == validator && delegation.Status == DelegationStatusActive {
			targetDelegation = delegation
			unbondAmount = new(big.Int).Set(delegation.Amount)
			if amount.Cmp(unbondAmount) < 0 {
				unbondAmount = amount
			}
			break
		}
	}

	if targetDelegation == nil {
		return fmt.Errorf("active delegation not found")
	}

	targetDelegation.Status = DelegationStatusUnbonding
	targetDelegation.UnbondTime = time.Now().Add(pos.config.UnbondingPeriod)

	v := pos.validators[validator]
	v.TotalDelegated.Sub(v.TotalDelegated, unbondAmount)
	pos.totalStake.Sub(pos.totalStake, unbondAmount)

	return nil
}

func (pos *PoS) CompleteUnbonding(delegator, validator core.Address) (*big.Int, error) {
	pos.mu.Lock()
	defer pos.mu.Unlock()

	delegations, exists := pos.delegators[delegator]
	if !exists {
		return nil, fmt.Errorf("no delegations found")
	}

	for _, delegation := range delegations {
		if delegation.Validator == validator && delegation.Status == DelegationStatusUnbonding {
			if time.Now().Before(delegation.UnbondTime) {
				return nil, fmt.Errorf("unbonding period not completed")
			}

			delegation.Status = DelegationStatusCompleted
			returnAmount := new(big.Int).Add(delegation.Amount, delegation.Rewards)
			return returnAmount, nil
		}
	}

	return nil, fmt.Errorf("unbonding delegation not found")
}

func (pos *PoS) SelectValidator(slot uint64) core.Address {
	pos.mu.RLock()
	defer pos.mu.RUnlock()

	activeValidators := pos.getActiveValidators()
	if len(activeValidators) == 0 {
		return core.Address{}
	}

	// Use weighted random selection based on stake
	totalWeight := big.NewInt(0)
	weights := make([]*big.Int, len(activeValidators))

	for i, validator := range activeValidators {
		weight := new(big.Int).Add(validator.Stake, validator.TotalDelegated)
		weights[i] = weight
		totalWeight.Add(totalWeight, weight)
	}

	if totalWeight.Cmp(big.NewInt(0)) == 0 {
		return core.Address{}
	}

	// Generate random number
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	randNum := new(big.Int).SetBytes(randBytes)
	randNum.Mod(randNum, totalWeight)

	// Select validator
	cumulative := big.NewInt(0)
	for i, validator := range activeValidators {
		cumulative.Add(cumulative, weights[i])
		if randNum.Cmp(cumulative) < 0 {
			return validator.Address
		}
	}

	return activeValidators[0].Address
}

func (pos *PoS) ValidateBlock(block *core.Block) error {
	pos.mu.RLock()
	defer pos.mu.RUnlock()

	// Check if block is signed by current validator
	expectedValidator := pos.SelectValidator(block.Header.Number.Uint64())
	if !bytes.Equal(block.Header.Validator.Bytes(), expectedValidator.Bytes()) {
		return fmt.Errorf("invalid validator: expected %s, got %s", 
			expectedValidator.Hex(), block.Header.Validator.Hex())
	}

	// Verify signature
	if len(block.Header.Signature) == 0 {
		return fmt.Errorf("block signature missing")
	}

	validator, exists := pos.validators[block.Header.Validator]
	if !exists {
		return fmt.Errorf("validator not found")
	}

	hash := block.ComputeHash()
	if !ecdsa.VerifyASN1(validator.PublicKey, hash.Bytes(), block.Header.Signature) {
		return fmt.Errorf("invalid block signature")
	}

	return nil
}

func (pos *PoS) SignBlock(block *core.Block) ([]byte, error) {
	if pos.privateKey == nil {
		return nil, fmt.Errorf("no private key available")
	}

	hash := block.ComputeHash()
	signature, err := ecdsa.SignASN1(rand.Reader, pos.privateKey, hash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to sign block: %w", err)
	}

	return signature, nil
}

func (pos *PoS) UpdateBlock(block *core.Block) {
	pos.mu.Lock()
	defer pos.mu.Unlock()

	validator, exists := pos.validators[block.Header.Validator]
	if !exists {
		return
	}

	validator.LastActive = time.Now()
	validator.TotalBlocks++

	// Calculate and distribute rewards
	pos.distributeRewards(block)
}

func (pos *PoS) distributeRewards(block *core.Block) {
	blockReward := new(big.Int).Set(pos.config.RewardRate)
	validator := pos.validators[block.Header.Validator]
	
	// Calculate commission
	commission := new(big.Int).Mul(blockReward, big.NewInt(int64(validator.Commission)))
	commission.Div(commission, big.NewInt(10000))
	
	// Validator reward
	validatorReward := new(big.Int).Add(commission, 
		new(big.Int).Mul(blockReward, big.NewInt(int64(validator.Stake.Uint64()))))
	validatorReward.Div(validatorReward, new(big.Int).Add(validator.Stake, validator.TotalDelegated))
	
	validator.Rewards.Add(validator.Rewards, validatorReward)
	
	// Delegator rewards
	if validator.TotalDelegated.Cmp(big.NewInt(0)) > 0 {
		delegatorReward := new(big.Int).Sub(blockReward, commission)
		delegatorReward.Sub(delegatorReward, new(big.Int).Mul(validatorReward, 
			new(big.Int).Sub(new(big.Int).Add(validator.Stake, validator.TotalDelegated), validator.Stake)))
		delegatorReward.Div(delegatorReward, validator.TotalDelegated)
		
		pos.distributeDelegatorRewards(validator.Address, delegatorReward)
	}
}

func (pos *PoS) distributeDelegatorRewards(validator core.Address, reward *big.Int) {
	for delegator, delegations := range pos.delegators {
		for _, delegation := range delegations {
			if delegation.Validator == validator && delegation.Status == DelegationStatusActive {
				delegationReward := new(big.Int).Mul(reward, delegation.Amount)
				delegationReward.Div(delegationReward, pos.validators[validator].TotalDelegated)
				delegation.Rewards.Add(delegation.Rewards, delegationReward)
			}
		}
	}
}

func (pos *PoS) Slash(validator core.Address, reason string) error {
	pos.mu.Lock()
	defer pos.mu.Unlock()

	v, exists := pos.validators[validator]
	if !exists {
		return fmt.Errorf("validator not found")
	}

	if v.Status == ValidatorStatusSlashed {
		return fmt.Errorf("validator already slashed")
	}

	// Calculate slash amount
	slashAmount := new(big.Int).Mul(pos.config.SlashRate, new(big.Int).Add(v.Stake, v.TotalDelegated))
	slashAmount.Div(slashAmount, big.NewInt(1000)) // Divide by 1000 to get percentage

	// Remove stake
	v.Stake.Sub(v.Stake, slashAmount)
	pos.totalStake.Sub(pos.totalStake, slashAmount)

	// Set status
	v.Status = ValidatorStatusSlashed

	return nil
}

func (pos *PoS) GetValidators() []*Validator {
	pos.mu.RLock()
	defer pos.mu.RUnlock()

	var validators []*Validator
	for _, v := range pos.validators {
		validators = append(validators, v)
	}

	// Sort by total stake (self + delegated)
	sort.Slice(validators, func(i, j int) bool {
		totalI := new(big.Int).Add(validators[i].Stake, validators[i].TotalDelegated)
		totalJ := new(big.Int).Add(validators[j].Stake, validators[j].TotalDelegated)
		return totalI.Cmp(totalJ) > 0
	})

	return validators
}

func (pos *PoS) GetDelegations(delegator core.Address) []*Delegation {
	pos.mu.RLock()
	defer pos.mu.RUnlock()

	return pos.delegators[delegator]
}

func (pos *PoS) GetTotalStake() *big.Int {
	pos.mu.RLock()
	defer pos.mu.RUnlock()

	return new(big.Int).Set(pos.totalStake)
}

func (pos *PoS) getActiveValidators() []*Validator {
	var active []*Validator
	for _, v := range pos.validators {
		if v.Status == ValidatorStatusActive {
			active = append(active, v)
		}
	}
	return active
}

func (pos *PoS) isCurrentValidator(address core.Address) bool {
	activeValidators := pos.getActiveValidators()
	for _, v := range activeValidators {
		if bytes.Equal(v.Address.Bytes(), address.Bytes()) {
			return true
		}
	}
	return false
}

func (pos *PoS) UpdateSlot(slot uint64) {
	pos.mu.Lock()
	defer pos.mu.Unlock()
	
	pos.currentSlot = slot
}

func (pos *PoS) GetCurrentSlot() uint64 {
	pos.mu.RLock()
	defer pos.mu.RUnlock()
	
	return pos.currentSlot
}

func (pos *PoS) GetEpoch() uint64 {
	return pos.GetCurrentSlot() / pos.config.EpochLength
}

func (pos *PoS) GetSlotTime(slot uint64) time.Time {
	genesisTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return genesisTime.Add(time.Duration(slot) * pos.config.SlotDuration)
}
