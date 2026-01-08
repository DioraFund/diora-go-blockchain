// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Pausable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";

contract DIO is ERC20, ERC20Burnable, ERC20Pausable, Ownable, ERC20Permit {
    // Anti-whale protection
    uint256 public maxWalletAmount;
    uint256 public maxTransactionAmount;
    
    // Tax mechanism
    uint256 public liquidityTax = 200; // 2%
    uint256 public stakingTax = 300;   // 3%
    uint256 public burnTax = 100;      // 1%
    uint256 public totalTax = 600;     // 6%
    
    // Addresses for tax distribution
    address public liquidityWallet;
    address public stakingWallet;
    address public burnWallet;
    
    // Excluded from tax
    mapping(address => bool) private _isExcludedFromTax;
    mapping(address => bool) private _isExcludedFromMaxAmount;
    
    // Staking mechanism
    struct Stake {
        uint256 amount;
        uint256 lockPeriod;
        uint256 startTime;
        uint256 rewardRate;
        bool active;
    }
    
    mapping(address => Stake[]) public stakes;
    mapping(address => uint256) public totalStaked;
    uint256 public totalStakedGlobally;
    uint256 public annualRewardRate = 1000; // 10% annual
    
    // Governance
    struct Proposal {
        uint256 id;
        address proposer;
        string title;
        string description;
        uint256 value;
        uint256 startTime;
        uint256 endTime;
        uint256 votesFor;
        uint256 votesAgainst;
        mapping(address => bool) hasVoted;
        bool executed;
    }
    
    mapping(uint256 => Proposal) public proposals;
    uint256 public proposalCount;
    mapping(address => uint256) public votingPower;
    
    // Reputation system
    mapping(address => uint256) public reputationScore;
    mapping(address => mapping(string => bool)) public achievements;
    
    // Events
    event Staked(address indexed user, uint256 amount, uint256 lockPeriod);
    event Unstaked(address indexed user, uint256 amount, uint256 reward);
    event ProposalCreated(uint256 indexed proposalId, address indexed proposer, string title);
    event Voted(uint256 indexed proposalId, address indexed voter, bool support, uint256 weight);
    event ProposalExecuted(uint256 indexed proposalId, bool passed);
    event ReputationUpdated(address indexed user, uint256 newScore, string reason);
    event AchievementUnlocked(address indexed user, string achievement);
    
    constructor(
        address _liquidityWallet,
        address _stakingWallet,
        address _burnWallet
    ) ERC20("Diora", "DIO") Ownable(msg.sender) ERC20Permit("Diora") {
        // Mint initial supply: 1 billion DIO
        uint256 initialSupply = 1_000_000_000 * 10**decimals();
        _mint(msg.sender, initialSupply);
        
        // Set tax wallets
        liquidityWallet = _liquidityWallet;
        stakingWallet = _stakingWallet;
        burnWallet = _burnWallet;
        
        // Set anti-whale limits (5% of total supply)
        maxWalletAmount = initialSupply * 5 / 100;
        maxTransactionAmount = initialSupply * 2 / 100;
        
        // Exclude owner and tax wallets from limits
        _isExcludedFromMaxAmount[msg.sender] = true;
        _isExcludedFromMaxAmount[liquidityWallet] = true;
        _isExcludedFromMaxAmount[stakingWallet] = true;
        _isExcludedFromMaxAmount[burnWallet] = true;
        
        // Exclude from tax
        _isExcludedFromTax[msg.sender] = true;
        _isExcludedFromTax[address(this)] = true;
        _isExcludedFromTax[liquidityWallet] = true;
        _isExcludedFromTax[stakingWallet] = true;
        _isExcludedFromTax[burnWallet] = true;
    }
    
    // Override transfer function to include tax and anti-whale
    function _transfer(
        address from,
        address to,
        uint256 amount
    ) internal override {
        require(!_isBlacklisted[from] && !_isBlacklisted[to], "Transfer from/to blacklisted address");
        
        if (!_isExcludedFromMaxAmount[to]) {
            require(balanceOf(to) + amount <= maxWalletAmount, "Max wallet amount exceeded");
        }
        
        if (!_isExcludedFromMaxAmount[from]) {
            require(amount <= maxTransactionAmount, "Max transaction amount exceeded");
        }
        
        uint256 taxAmount = 0;
        if (!_isExcludedFromTax[from] && !_isExcludedFromTax[to]) {
            taxAmount = (amount * totalTax) / 10000;
            
            // Distribute taxes
            if (liquidityTax > 0) {
                uint256 liquidityAmount = (amount * liquidityTax) / 10000;
                super._transfer(from, liquidityWallet, liquidityAmount);
            }
            
            if (stakingTax > 0) {
                uint256 stakingAmount = (amount * stakingTax) / 10000;
                super._transfer(from, stakingWallet, stakingAmount);
            }
            
            if (burnTax > 0) {
                uint256 burnAmount = (amount * burnTax) / 10000;
                super._transfer(from, burnWallet, burnAmount);
                _burn(burnWallet, burnAmount);
            }
        }
        
        uint256 transferAmount = amount - taxAmount;
        super._transfer(from, to, transferAmount);
        
        // Update reputation for active users
        _updateReputation(from, 1, "transfer");
        _updateReputation(to, 1, "receive");
    }
    
    // Staking functions
    function stake(uint256 amount, uint256 lockPeriod) external {
        require(amount > 0, "Amount must be greater than 0");
        require(lockPeriod >= 30 days, "Minimum lock period is 30 days");
        require(balanceOf(msg.sender) >= amount, "Insufficient balance");
        
        _transfer(msg.sender, address(this), amount);
        
        uint256 rewardRate = annualRewardRate + (lockPeriod / 90 days) * 200; // Bonus for longer locks
        
        stakes[msg.sender].push(Stake({
            amount: amount,
            lockPeriod: lockPeriod,
            startTime: block.timestamp,
            rewardRate: rewardRate,
            active: true
        }));
        
        totalStaked[msg.sender] += amount;
        totalStakedGlobally += amount;
        votingPower[msg.sender] += amount;
        
        emit Staked(msg.sender, amount, lockPeriod);
        _updateReputation(msg.sender, 10, "stake");
        
        // Check achievements
        if (totalStaked[msg.sender] >= 10000 * 10**decimals() && !achievements[msg.sender]["big_staker"]) {
            achievements[msg.sender]["big_staker"] = true;
            emit AchievementUnlocked(msg.sender, "Big Staker");
        }
    }
    
    function unstake(uint256 stakeIndex) external {
        require(stakeIndex < stakes[msg.sender].length, "Invalid stake index");
        
        Stake storage userStake = stakes[msg.sender][stakeIndex];
        require(userStake.active, "Stake already inactive");
        require(block.timestamp >= userStake.startTime + userStake.lockPeriod, "Stake still locked");
        
        uint256 reward = calculateReward(msg.sender, stakeIndex);
        uint256 totalAmount = userStake.amount + reward;
        
        userStake.active = false;
        totalStaked[msg.sender] -= userStake.amount;
        totalStakedGlobally -= userStake.amount;
        votingPower[msg.sender] -= userStake.amount;
        
        require(balanceOf(address(this)) >= totalAmount, "Insufficient contract balance");
        _transfer(address(this), msg.sender, totalAmount);
        
        emit Unstaked(msg.sender, userStake.amount, reward);
        _updateReputation(msg.sender, 5, "unstake");
    }
    
    function calculateReward(address user, uint256 stakeIndex) public view returns (uint256) {
        Stake storage userStake = stakes[user][stakeIndex];
        if (!userStake.active) return 0;
        
        uint256 stakingDuration = block.timestamp - userStake.startTime;
        uint256 maxDuration = userStake.lockPeriod;
        
        if (stakingDuration > maxDuration) {
            stakingDuration = maxDuration;
        }
        
        return (userStake.amount * userStake.rewardRate * stakingDuration) / (365 days * 10000);
    }
    
    // Governance functions
    function createProposal(
        string memory title,
        string memory description,
        uint256 value
    ) external returns (uint256) {
        require(votingPower[msg.sender] >= 1000 * 10**decimals(), "Insufficient voting power");
        require(bytes(title).length > 0, "Title cannot be empty");
        
        proposalCount++;
        Proposal storage proposal = proposals[proposalCount];
        
        proposal.id = proposalCount;
        proposal.proposer = msg.sender;
        proposal.title = title;
        proposal.description = description;
        proposal.value = value;
        proposal.startTime = block.timestamp;
        proposal.endTime = block.timestamp + 7 days;
        
        emit ProposalCreated(proposalCount, msg.sender, title);
        _updateReputation(msg.sender, 15, "create_proposal");
        
        return proposalCount;
    }
    
    function vote(uint256 proposalId, bool support) external {
        require(proposalId > 0 && proposalId <= proposalCount, "Invalid proposal");
        require(votingPower[msg.sender] > 0, "No voting power");
        
        Proposal storage proposal = proposals[proposalId];
        require(block.timestamp >= proposal.startTime && block.timestamp <= proposal.endTime, "Voting period ended");
        require(!proposal.hasVoted[msg.sender], "Already voted");
        
        proposal.hasVoted[msg.sender] = true;
        
        if (support) {
            proposal.votesFor += votingPower[msg.sender];
        } else {
            proposal.votesAgainst += votingPower[msg.sender];
        }
        
        emit Voted(proposalId, msg.sender, support, votingPower[msg.sender]);
        _updateReputation(msg.sender, 3, "vote");
    }
    
    function executeProposal(uint256 proposalId) external {
        require(proposalId > 0 && proposalId <= proposalCount, "Invalid proposal");
        
        Proposal storage proposal = proposals[proposalId];
        require(block.timestamp > proposal.endTime, "Voting period not ended");
        require(!proposal.executed, "Proposal already executed");
        
        bool passed = proposal.votesFor > proposal.votesAgainst;
        proposal.executed = true;
        
        if (passed) {
            // Execute proposal logic here
            // This could be parameter changes, fund transfers, etc.
        }
        
        emit ProposalExecuted(proposalId, passed);
        _updateReputation(proposal.proposer, passed ? 25 : 5, passed ? "proposal_passed" : "proposal_rejected");
    }
    
    // Reputation system
    function _updateReputation(address user, uint256 points, string memory reason) internal {
        reputationScore[user] += points;
        emit ReputationUpdated(user, reputationScore[user], reason);
        
        // Check for reputation-based achievements
        if (reputationScore[user] >= 1000 && !achievements[user]["reputed"]) {
            achievements[user]["reputed"] = true;
            emit AchievementUnlocked(user, "Reputed User");
        }
    }
    
    // Admin functions
    function setTaxes(
        uint256 _liquidityTax,
        uint256 _stakingTax,
        uint256 _burnTax
    ) external onlyOwner {
        require(_liquidityTax + _stakingTax + _burnTax <= 1000, "Total tax cannot exceed 10%");
        
        liquidityTax = _liquidityTax;
        stakingTax = _stakingTax;
        burnTax = _burnTax;
        totalTax = _liquidityTax + _stakingTax + _burnTax;
    }
    
    function setMaxAmounts(
        uint256 _maxWalletAmount,
        uint256 _maxTransactionAmount
    ) external onlyOwner {
        maxWalletAmount = _maxWalletAmount;
        maxTransactionAmount = _maxTransactionAmount;
    }
    
    function setTaxWallets(
        address _liquidityWallet,
        address _stakingWallet,
        address _burnWallet
    ) external onlyOwner {
        liquidityWallet = _liquidityWallet;
        stakingWallet = _stakingWallet;
        burnWallet = _burnWallet;
    }
    
    function excludeFromTax(address account, bool excluded) external onlyOwner {
        _isExcludedFromTax[account] = excluded;
    }
    
    function excludeFromMaxAmount(address account, bool excluded) external onlyOwner {
        _isExcludedFromMaxAmount[account] = excluded;
    }
    
    function pause() external onlyOwner {
        _pause();
    }
    
    function unpause() external onlyOwner {
        _unpause();
    }
    
    // Blacklist functionality
    mapping(address => bool) private _isBlacklisted;
    
    function blacklist(address account, bool blacklisted) external onlyOwner {
        _isBlacklisted[account] = blacklisted;
    }
    
    // View functions
    function getUserStakes(address user) external view returns (Stake[] memory) {
        return stakes[user];
    }
    
    function getProposal(uint256 proposalId) external view returns (
        uint256 id,
        address proposer,
        string memory title,
        string memory description,
        uint256 startTime,
        uint256 endTime,
        uint256 votesFor,
        uint256 votesAgainst,
        bool executed
    ) {
        Proposal storage proposal = proposals[proposalId];
        return (
            proposal.id,
            proposal.proposer,
            proposal.title,
            proposal.description,
            proposal.startTime,
            proposal.endTime,
            proposal.votesFor,
            proposal.votesAgainst,
            proposal.executed
        );
    }
    
    function isExcludedFromTax(address account) external view returns (bool) {
        return _isExcludedFromTax[account];
    }
    
    function isExcludedFromMaxAmount(address account) external view returns (bool) {
        return _isExcludedFromMaxAmount[account];
    }
    
    function isBlacklisted(address account) external view returns (bool) {
        return _isBlacklisted[account];
    }
    
    // The following functions are overrides required by Solidity
    function _update(address from, address to, uint256 value)
        internal
        override(ERC20, ERC20Pausable)
    {
        super._update(from, to, value);
    }
}
