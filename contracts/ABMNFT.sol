// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/**
 * @title ABMNFT
 * @dev Official ABM Foundation NFT Collection
 * @notice This contract creates a collection of 15 unique NFTs representing ABM Diora ecosystem
 * @author ABM Foundation
 * @version 1.0.0
 */
contract ABMNFT is ERC721URIStorage, Ownable, ReentrancyGuard {
    using Counters for Counters.Counter;
    using Strings for uint256;

    // Counter for token IDs
    Counters.Counter private _tokenIdCounter;

    // Maximum supply of NFTs
    uint256 public constant MAX_SUPPLY = 15;

    // Minting price (in DIO tokens)
    uint256 public mintPrice = 100 * 10**18; // 100 DIO

    // Maximum mints per transaction
    uint256 public constant MAX_MINT_PER_TX = 3;

    // Maximum mints per wallet
    uint256 public constant MAX_MINT_PER_WALLET = 5;

    // Mapping to track mints per wallet
    mapping(address => uint256) public walletMintCount;

    // Base URI for metadata
    string private _baseTokenURI;

    // Contract URI for collection metadata
    string private _contractURI;

    // Minting status
    bool public isMintingActive = false;

    // Presale status
    bool public isPresaleActive = false;

    // Presale price
    uint256 public presalePrice = 50 * 10**18; // 50 DIO

    // Presale whitelist
    mapping(address => bool) public presaleWhitelist;

    // Team wallet for withdrawals
    address payable public teamWallet;

    // Events
    event Minted(address indexed to, uint256 indexed tokenId, string tokenURI);
    event BatchMinted(address indexed to, uint256[] tokenIds);
    event MintingStatusChanged(bool isActive);
    event PresaleStatusChanged(bool isActive);
    event PriceUpdated(uint256 newPrice);
    event BaseURIUpdated(string newBaseURI);
    event ContractURIUpdated(string newContractURI);

    /**
     * @dev Constructor
     * @param name_ Name of the NFT collection
     * @param symbol_ Symbol of the NFT collection
     * @param baseTokenURI_ Base URI for token metadata
     * @param contractURI_ URI for collection metadata
     * @param teamWallet_ Address for team withdrawals
     */
    constructor(
        string memory name_,
        string memory symbol_,
        string memory baseTokenURI_,
        string memory contractURI_,
        address payable teamWallet_
    ) ERC721(name_, symbol_) {
        _baseTokenURI = baseTokenURI_;
        _contractURI = contractURI_;
        teamWallet = teamWallet_;
    }

    /**
     * @dev Modifier to check if minting is active
     */
    modifier whenMintingActive() {
        require(isMintingActive, "ABMNFT: Minting is not active");
        _;
    }

    /**
     * @dev Modifier to check if presale is active
     */
    modifier whenPresaleActive() {
        require(isPresaleActive, "ABMNFT: Presale is not active");
        _;
    }

    /**
     * @dev Modifier to check if address is whitelisted for presale
     */
    modifier onlyWhitelisted() {
        require(presaleWhitelist[msg.sender], "ABMNFT: Address is not whitelisted for presale");
        _;
    }

    /**
     * @dev Modifier to check supply constraints
     */
    modifier supplyCheck(uint256 quantity) {
        require(_tokenIdCounter.current() + quantity <= MAX_SUPPLY, "ABMNFT: Exceeds maximum supply");
        _;
    }

    /**
     * @dev Modifier to check mint limits
     */
    modifier mintLimitCheck(uint256 quantity) {
        require(quantity <= MAX_MINT_PER_TX, "ABMNFT: Exceeds max mints per transaction");
        require(walletMintCount[msg.sender] + quantity <= MAX_MINT_PER_WALLET, "ABMNFT: Exceeds max mints per wallet");
        _;
    }

    /**
     * @dev Public mint function
     * @param quantity Number of NFTs to mint
     */
    function publicMint(uint256 quantity) 
        external 
        payable 
        nonReentrant 
        whenMintingActive 
        supplyCheck(quantity) 
        mintLimitCheck(quantity) 
    {
        require(msg.value >= mintPrice * quantity, "ABMNFT: Insufficient payment");
        
        _mintBatch(msg.sender, quantity);
    }

    /**
     * @dev Presale mint function for whitelisted addresses
     * @param quantity Number of NFTs to mint
     */
    function presaleMint(uint256 quantity) 
        external 
        payable 
        nonReentrant 
        whenPresaleActive 
        onlyWhitelisted 
        supplyCheck(quantity) 
        mintLimitCheck(quantity) 
    {
        require(msg.value >= presalePrice * quantity, "ABMNFT: Insufficient presale payment");
        
        _mintBatch(msg.sender, quantity);
    }

    /**
     * @dev Owner mint function
     * @param to Address to mint to
     * @param quantity Number of NFTs to mint
     */
    function ownerMint(address to, uint256 quantity) 
        external 
        onlyOwner 
        supplyCheck(quantity) 
    {
        _mintBatch(to, quantity);
    }

    /**
     * @dev Internal batch mint function
     * @param to Address to mint to
     * @param quantity Number of NFTs to mint
     */
    function _mintBatch(address to, uint256 quantity) internal {
        uint256[] memory tokenIds = new uint256[](quantity);
        
        for (uint256 i = 0; i < quantity; i++) {
            _tokenIdCounter.increment();
            uint256 tokenId = _tokenIdCounter.current();
            _safeMint(to, tokenId);
            tokenIds[i] = tokenId;
        }
        
        walletMintCount[to] += quantity;
        
        emit BatchMinted(to, tokenIds);
    }

    /**
     * @dev Toggle minting status
     * @param isActive New minting status
     */
    function toggleMinting(bool isActive) external onlyOwner {
        isMintingActive = isActive;
        emit MintingStatusChanged(isActive);
    }

    /**
     * @dev Toggle presale status
     * @param isActive New presale status
     */
    function togglePresale(bool isActive) external onlyOwner {
        isPresaleActive = isActive;
        emit PresaleStatusChanged(isActive);
    }

    /**
     * @dev Add addresses to presale whitelist
     * @param addresses Array of addresses to whitelist
     */
    function addToPresaleWhitelist(address[] calldata addresses) external onlyOwner {
        for (uint256 i = 0; i < addresses.length; i++) {
            presaleWhitelist[addresses[i]] = true;
        }
    }

    /**
     * @dev Remove addresses from presale whitelist
     * @param addresses Array of addresses to remove from whitelist
     */
    function removeFromPresaleWhitelist(address[] calldata addresses) external onlyOwner {
        for (uint256 i = 0; i < addresses.length; i++) {
            presaleWhitelist[addresses[i]] = false;
        }
    }

    /**
     * @dev Update minting price
     * @param newPrice New minting price in wei
     */
    function updateMintPrice(uint256 newPrice) external onlyOwner {
        mintPrice = newPrice;
        emit PriceUpdated(newPrice);
    }

    /**
     * @dev Update presale price
     * @param newPrice New presale price in wei
     */
    function updatePresalePrice(uint256 newPrice) external onlyOwner {
        presalePrice = newPrice;
        emit PriceUpdated(newPrice);
    }

    /**
     * @dev Update base URI
     * @param baseTokenURI_ New base URI
     */
    function setBaseURI(string memory baseTokenURI_) external onlyOwner {
        _baseTokenURI = baseTokenURI_;
        emit BaseURIUpdated(baseTokenURI_);
    }

    /**
     * @dev Update contract URI
     * @param contractURI_ New contract URI
     */
    function setContractURI(string memory contractURI_) external onlyOwner {
        _contractURI = contractURI_;
        emit ContractURIUpdated(contractURI_);
    }

    /**
     * @dev Update team wallet
     * @param newTeamWallet New team wallet address
     */
    function updateTeamWallet(address payable newTeamWallet) external onlyOwner {
        teamWallet = newTeamWallet;
    }

    /**
     * @dev Withdraw funds to team wallet
     */
    function withdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "ABMNFT: No funds to withdraw");
        
        (bool success, ) = teamWallet.call{value: balance}("");
        require(success, "ABMNFT: Withdrawal failed");
    }

    /**
     * @dev Emergency withdraw function
     */
    function emergencyWithdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "ABMNFT: No funds to withdraw");
        
        (bool success, ) = teamWallet.call{value: balance}("");
        require(success, "ABMNFT: Emergency withdrawal failed");
    }

    /**
     * @dev Get current token ID
     */
    function getCurrentTokenId() external view returns (uint256) {
        return _tokenIdCounter.current();
    }

    /**
     * @dev Get total supply
     */
    function totalSupply() external view returns (uint256) {
        return _tokenIdCounter.current();
    }

    /**
     * @dev Get remaining supply
     */
    function remainingSupply() external view returns (uint256) {
        return MAX_SUPPLY - _tokenIdCounter.current();
    }

    /**
     * @dev Get mint price for human readable format
     */
    function getMintPrice() external view returns (uint256) {
        return mintPrice;
    }

    /**
     * @dev Get presale price for human readable format
     */
    function getPresalePrice() external view returns (uint256) {
        return presalePrice;
    }

    /**
     * @dev Check if address is whitelisted
     */
    function isWhitelisted(address user) external view returns (bool) {
        return presaleWhitelist[user];
    }

    /**
     * @dev Get tokens owned by address
     */
    function tokensOfOwner(address owner) external view returns (uint256[] memory) {
        uint256 balance = balanceOf(owner);
        uint256[] memory tokens = new uint256[](balance);
        
        for (uint256 i = 0; i < balance; i++) {
            tokens[i] = tokenOfOwnerByIndex(owner, i);
        }
        
        return tokens;
    }

    /**
     * @dev Internal function to return base URI
     */
    function _baseURI() internal view override returns (string memory) {
        return _baseTokenURI;
    }

    /**
     * @dev Override tokenURI function
     */
    function tokenURI(uint256 tokenId) 
        public 
        view 
        override(ERC721URIStorage) 
        returns (string memory) 
    {
        require(_exists(tokenId), "ABMNFT: URI query for nonexistent token");
        
        string memory baseURI = _baseURI();
        return bytes(baseURI).length > 0 
            ? string(abi.encodePacked(baseURI, tokenId.toString(), ".json"))
            : "";
    }

    /**
     * @dev Return contract URI for collection metadata
     */
    function contractURI() external view returns (string memory) {
        return _contractURI;
    }

    /**
     * @dev Override supportsInterface
     */
    function supportsInterface(bytes4 interfaceId) 
        public 
        view 
        override(ERC721URIStorage) 
        returns (bool) 
    {
        return super.supportsInterface(interfaceId);
    }

    /**
     * @dev Receive function to handle direct ETH transfers
     */
    receive() external payable {
        revert("ABMNFT: Direct transfers not allowed");
    }

    /**
     * @dev Fallback function
     */
    fallback() external payable {
        revert("ABMNFT: Direct transfers not allowed");
    }
}
