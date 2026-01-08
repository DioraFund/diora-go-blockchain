// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract DioraNFT is ERC721, ERC721URIStorage, ERC721Burnable, Ownable, ReentrancyGuard {
    using Counters for Counters.Counter;
    
    Counters.Counter private _tokenIdCounter;
    
    // NFT metadata
    struct NFTMetadata {
        string name;
        string description;
        string image;
        string category;
        uint256 createdAt;
        address creator;
        bool verified;
    }
    
    mapping(uint256 => NFTMetadata) public nftMetadata;
    mapping(string => bool) private _uriExists;
    
    // Royalty system
    struct Royalty {
        address creator;
        uint256 percentage; // Basis points (100 = 1%)
    }
    
    mapping(uint256 => Royalty) public royalties;
    uint256 public defaultRoyalty = 250; // 2.5%
    
    // Marketplace
    struct Listing {
        uint256 tokenId;
        address seller;
        uint256 price;
        bool active;
        uint256 createdAt;
    }
    
    mapping(uint256 => Listing) public listings;
    uint256 public marketplaceFee = 250; // 2.5%
    address public marketplaceFeeRecipient;
    
    // Collections
    struct Collection {
        string name;
        string description;
        address creator;
        uint256 createdAt;
        bool verified;
        uint256 totalSupply;
        uint256 maxSupply;
    }
    
    mapping(string => Collection) public collections;
    mapping(uint256 => string) public tokenToCollection;
    string[] public collectionNames;
    
    // Events
    event NFTMinted(uint256 indexed tokenId, address indexed creator, string collection, string uri);
    event NFTListed(uint256 indexed tokenId, address indexed seller, uint256 price);
    event NFTSold(uint256 indexed tokenId, address indexed seller, address indexed buyer, uint256 price);
    event NFTBid(uint256 indexed tokenId, address indexed bidder, uint256 amount);
    event CollectionCreated(string indexed name, address indexed creator);
    event RoyaltySet(uint256 indexed tokenId, address creator, uint256 percentage);
    
    constructor(address _marketplaceFeeRecipient) ERC721("DioraNFT", "DNFT") Ownable(msg.sender) {
        marketplaceFeeRecipient = _marketplaceFeeRecipient;
    }
    
    // Minting functions
    function mint(
        address to,
        string memory uri,
        string memory name,
        string memory description,
        string memory category,
        string memory collection
    ) external returns (uint256) {
        require(bytes(uri).length > 0, "URI cannot be empty");
        require(!_uriExists[uri], "URI already exists");
        
        uint256 tokenId = _tokenIdCounter.current();
        _tokenIdCounter.increment();
        
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
        
        nftMetadata[tokenId] = NFTMetadata({
            name: name,
            description: description,
            image: uri,
            category: category,
            createdAt: block.timestamp,
            creator: to,
            verified: false
        });
        
        _uriExists[uri] = true;
        tokenToCollection[tokenId] = collection;
        
        // Update collection supply
        if (bytes(collection).length > 0) {
            collections[collection].totalSupply++;
        }
        
        // Set default royalty
        royalties[tokenId] = Royalty({
            creator: to,
            percentage: defaultRoyalty
        });
        
        emit NFTMinted(tokenId, to, collection, uri);
        return tokenId;
    }
    
    function mintCollection(
        string memory collectionName,
        string memory collectionDescription,
        uint256 maxSupply,
        address to,
        string memory uri,
        string memory name,
        string memory description,
        string memory category
    ) external returns (uint256) {
        require(bytes(collectionName).length > 0, "Collection name cannot be empty");
        require(collections[collectionName].creator == address(0), "Collection already exists");
        
        // Create collection
        collections[collectionName] = Collection({
            name: collectionName,
            description: collectionDescription,
            creator: msg.sender,
            createdAt: block.timestamp,
            verified: false,
            totalSupply: 0,
            maxSupply: maxSupply
        });
        
        collectionNames.push(collectionName);
        emit CollectionCreated(collectionName, msg.sender);
        
        // Mint NFT in collection
        return mint(to, uri, name, description, category, collectionName);
    }
    
    // Marketplace functions
    function listItem(uint256 tokenId, uint256 price) external {
        require(ownerOf(tokenId) == msg.sender, "Not token owner");
        require(price > 0, "Price must be greater than 0");
        require(!listings[tokenId].active, "Already listed");
        
        listings[tokenId] = Listing({
            tokenId: tokenId,
            seller: msg.sender,
            price: price,
            active: true,
            createdAt: block.timestamp
        });
        
        emit NFTListed(tokenId, msg.sender, price);
    }
    
    function buyItem(uint256 tokenId) external payable nonReentrant {
        Listing storage listing = listings[tokenId];
        require(listing.active, "Not for sale");
        require(msg.value >= listing.price, "Insufficient payment");
        
        address seller = listing.seller;
        uint256 price = listing.price;
        
        // Calculate fees
        uint256 marketplaceAmount = (price * marketplaceFee) / 10000;
        uint256 royaltyAmount = 0;
        
        // Pay royalty if set
        if (royalties[tokenId].creator != address(0)) {
            royaltyAmount = (price * royalties[tokenId].percentage) / 10000;
            payable(royalties[tokenId].creator).transfer(royaltyAmount);
        }
        
        // Pay marketplace fee
        if (marketplaceAmount > 0) {
            payable(marketplaceFeeRecipient).transfer(marketplaceAmount);
        }
        
        // Pay seller
        uint256 sellerAmount = price - marketplaceAmount - royaltyAmount;
        payable(seller).transfer(sellerAmount);
        
        // Transfer NFT
        _transfer(seller, msg.sender, tokenId);
        
        // Update listing
        listing.active = false;
        
        emit NFTSold(tokenId, seller, msg.sender, price);
    }
    
    function cancelListing(uint256 tokenId) external {
        require(listings[tokenId].seller == msg.sender, "Not seller");
        require(listings[tokenId].active, "Not active");
        
        listings[tokenId].active = false;
    }
    
    // Royalty functions
    function setRoyalty(uint256 tokenId, uint256 percentage) external {
        require(ownerOf(tokenId) == msg.sender, "Not token owner");
        require(percentage <= 1000, "Royalty cannot exceed 10%");
        
        royalties[tokenId] = Royalty({
            creator: msg.sender,
            percentage: percentage
        });
        
        emit RoyaltySet(tokenId, msg.sender, percentage);
    }
    
    function getRoyaltyInfo(uint256 tokenId, uint256 salePrice) external view returns (address, uint256) {
        Royalty memory royalty = royalties[tokenId];
        uint256 royaltyAmount = (salePrice * royalty.percentage) / 10000;
        return (royalty.creator, royaltyAmount);
    }
    
    // Collection management
    function verifyCollection(string memory collectionName, bool verified) external onlyOwner {
        require(collections[collectionName].creator != address(0), "Collection does not exist");
        collections[collectionName].verified = verified;
    }
    
    function verifyNFT(uint256 tokenId, bool verified) external onlyOwner {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        nftMetadata[tokenId].verified = verified;
    }
    
    // Admin functions
    function setMarketplaceFee(uint256 fee) external onlyOwner {
        require(fee <= 1000, "Fee cannot exceed 10%");
        marketplaceFee = fee;
    }
    
    function setMarketplaceFeeRecipient(address recipient) external onlyOwner {
        marketplaceFeeRecipient = recipient;
    }
    
    function setDefaultRoyalty(uint256 royalty) external onlyOwner {
        require(royalty <= 1000, "Royalty cannot exceed 10%");
        defaultRoyalty = royalty;
    }
    
    function emergencyWithdraw() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }
    
    // View functions
    function getTokenMetadata(uint256 tokenId) external view returns (
        string memory name,
        string memory description,
        string memory image,
        string memory category,
        uint256 createdAt,
        address creator,
        bool verified
    ) {
        NFTMetadata memory metadata = nftMetadata[tokenId];
        return (
            metadata.name,
            metadata.description,
            metadata.image,
            metadata.category,
            metadata.createdAt,
            metadata.creator,
            metadata.verified
        );
    }
    
    function getListing(uint256 tokenId) external view returns (
        address seller,
        uint256 price,
        bool active,
        uint256 createdAt
    ) {
        Listing memory listing = listings[tokenId];
        return (
            listing.seller,
            listing.price,
            listing.active,
            listing.createdAt
        );
    }
    
    function getCollection(string memory collectionName) external view returns (
        string memory name,
        string memory description,
        address creator,
        uint256 createdAt,
        bool verified,
        uint256 totalSupply,
        uint256 maxSupply
    ) {
        Collection memory collection = collections[collectionName];
        return (
            collection.name,
            collection.description,
            collection.creator,
            collection.createdAt,
            collection.verified,
            collection.totalSupply,
            collection.maxSupply
        );
    }
    
    function getCollections() external view returns (string[] memory) {
        return collectionNames;
    }
    
    function getTokensByCollection(string memory collectionName) external view returns (uint256[] memory) {
        uint256[] memory tokens = new uint256[](collections[collectionName].totalSupply);
        uint256 index = 0;
        
        for (uint256 i = 1; i <= _tokenIdCounter.current(); i++) {
            if (keccak256(bytes(tokenToCollection[i])) == keccak256(bytes(collectionName))) {
                tokens[index] = i;
                index++;
            }
        }
        
        return tokens;
    }
    
    function getTokensByOwner(address owner) external view returns (uint256[] memory) {
        uint256 balance = balanceOf(owner);
        uint256[] memory tokens = new uint256[](balance);
        uint256 index = 0;
        
        for (uint256 i = 1; i <= _tokenIdCounter.current(); i++) {
            if (ownerOf(i) == owner) {
                tokens[index] = i;
                index++;
            }
        }
        
        return tokens;
    }
    
    // The following functions are overrides required by Solidity
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
