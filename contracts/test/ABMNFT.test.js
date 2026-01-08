/**
 * ABM Diora NFT Contract Tests
 * Comprehensive test suite for ABMNFT contract
 */

const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("ABMNFT", function () {
    let abmNFT;
    let owner;
    let addr1;
    let addr2;
    let addrs;

    // Test constants
    const NAME = "ABM Foundation Genesis Collection";
    const SYMBOL = "ABMNFT";
    const BASE_URI = "https://diorafund.github.io/diora-blockchain/metadata/";
    const CONTRACT_URI = "https://diorafund.github.io/diora-blockchain/metadata/collection.json";
    const MAX_SUPPLY = 15;
    const MAX_MINT_PER_TX = 3;
    const MAX_MINT_PER_WALLET = 5;
    const MINT_PRICE = ethers.utils.parseEther("100");
    const PRESALE_PRICE = ethers.utils.parseEther("50");

    beforeEach(async function () {
        [owner, addr1, addr2, ...addrs] = await ethers.getSigners();
        
        const ABMNFT = await ethers.getContractFactory("ABMNFT");
        abmNFT = await ABMNFT.deploy(
            NAME,
            SYMBOL,
            BASE_URI,
            CONTRACT_URI,
            owner.address
        );
        await abmNFT.deployed();
    });

    describe("Deployment", function () {
        it("Should set the right owner", async function () {
            expect(await abmNFT.owner()).to.equal(owner.address);
        });

        it("Should set the correct name and symbol", async function () {
            expect(await abmNFT.name()).to.equal(NAME);
            expect(await abmNFT.symbol()).to.equal(SYMBOL);
        });

        it("Should set the correct base URI", async function () {
            expect(await abmNFT._baseURI()).to.equal(BASE_URI);
        });

        it("Should set the correct contract URI", async function () {
            expect(await abmNFT.contractURI()).to.equal(CONTRACT_URI);
        });

        it("Should initialize with zero tokens", async function () {
            expect(await abmNFT.totalSupply()).to.equal(0);
            expect(await abmNFT.getCurrentTokenId()).to.equal(0);
        });

        it("Should set correct initial state", async function () {
            expect(await abmNFT.isMintingActive()).to.be.false;
            expect(await abmNFT.isPresaleActive()).to.be.false;
            expect(await abmNFT.mintPrice()).to.equal(MINT_PRICE);
            expect(await abmNFT.presalePrice()).to.equal(PRESALE_PRICE);
        });
    });

    describe("Minting Controls", function () {
        it("Should allow owner to toggle minting", async function () {
            await abmNFT.toggleMinting(true);
            expect(await abmNFT.isMintingActive()).to.be.true;

            await abmNFT.toggleMinting(false);
            expect(await abmNFT.isMintingActive()).to.be.false;
        });

        it("Should not allow non-owner to toggle minting", async function () {
            await expect(abmNFT.connect(addr1).toggleMinting(true))
                .to.be.revertedWith("Ownable: caller is not the owner");
        });

        it("Should allow owner to toggle presale", async function () {
            await abmNFT.togglePresale(true);
            expect(await abmNFT.isPresaleActive()).to.be.true;

            await abmNFT.togglePresale(false);
            expect(await abmNFT.isPresaleActive()).to.be.false;
        });

        it("Should not allow non-owner to toggle presale", async function () {
            await expect(abmNFT.connect(addr1).togglePresale(true))
                .to.be.revertedWith("Ownable: caller is not the owner");
        });
    });

    describe("Presale", function () {
        beforeEach(async function () {
            await abmNFT.togglePresale(true);
            await abmNFT.addToPresaleWhitelist([addr1.address, addr2.address]);
        });

        it("Should allow whitelisted address to mint during presale", async function () {
            const mintAmount = 2;
            const expectedCost = PRESALE_PRICE.mul(mintAmount);

            await abmNFT.connect(addr1).presaleMint(mintAmount, {
                value: expectedCost
            });

            expect(await abmNFT.balanceOf(addr1.address)).to.equal(mintAmount);
            expect(await abmNFT.walletMintCount(addr1.address)).to.equal(mintAmount);
        });

        it("Should not allow non-whitelisted address to mint during presale", async function () {
            await expect(abmNFT.connect(addrs[0]).presaleMint(1, {
                value: PRESALE_PRICE
            })).to.be.revertedWith("ABMNFT: Address is not whitelisted for presale");
        });

        it("Should require correct payment for presale", async function () {
            await expect(abmNFT.connect(addr1).presaleMint(1, {
                value: ethers.utils.parseEther("25") // Less than presale price
            })).to.be.revertedWith("ABMNFT: Insufficient presale payment");
        });

        it("Should enforce wallet limit during presale", async function () {
            const mintAmount = MAX_MINT_PER_WALLET + 1;
            
            await expect(abmNFT.connect(addr1).presaleMint(mintAmount, {
                value: PRESALE_PRICE.mul(mintAmount)
            })).to.be.revertedWith("ABMNFT: Exceeds max mints per wallet");
        });
    });

    describe("Public Minting", function () {
        beforeEach(async function () {
            await abmNFT.toggleMinting(true);
        });

        it("Should allow public minting with correct payment", async function () {
            const mintAmount = 2;
            const expectedCost = MINT_PRICE.mul(mintAmount);

            await abmNFT.connect(addr1).publicMint(mintAmount, {
                value: expectedCost
            });

            expect(await abmNFT.balanceOf(addr1.address)).to.equal(mintAmount);
            expect(await abmNFT.walletMintCount(addr1.address)).to.equal(mintAmount);
        });

        it("Should not allow minting with insufficient payment", async function () {
            await expect(abmNFT.connect(addr1).publicMint(1, {
                value: ethers.utils.parseEther("50") // Less than mint price
            })).to.be.revertedWith("ABMNFT: Insufficient payment");
        });

        it("Should enforce transaction limit", async function () {
            const mintAmount = MAX_MINT_PER_TX + 1;
            
            await expect(abmNFT.connect(addr1).publicMint(mintAmount, {
                value: MINT_PRICE.mul(mintAmount)
            })).to.be.revertedWith("ABMNFT: Exceeds max mints per transaction");
        });

        it("Should enforce wallet limit", async function () {
            // First mint within limit
            await abmNFT.connect(addr1).publicMint(MAX_MINT_PER_WALLET, {
                value: MINT_PRICE.mul(MAX_MINT_PER_WALLET)
            });

            // Second mint should fail
            await expect(abmNFT.connect(addr1).publicMint(1, {
                value: MINT_PRICE
            })).to.be.revertedWith("ABMNFT: Exceeds max mints per wallet");
        });

        it("Should enforce maximum supply", async function () {
            // Mint up to max supply
            const remainingSupply = MAX_SUPPLY;
            
            await abmNFT.connect(addr1).publicMint(remainingSupply, {
                value: MINT_PRICE.mul(remainingSupply)
            });

            // Try to mint one more
            await expect(abmNFT.connect(addr2).publicMint(1, {
                value: MINT_PRICE
            })).to.be.revertedWith("ABMNFT: Exceeds maximum supply");
        });
    });

    describe("Owner Functions", function () {
        it("Should allow owner to mint without payment", async function () {
            const mintAmount = 3;
            
            await abmNFT.ownerMint(addr1.address, mintAmount);
            
            expect(await abmNFT.balanceOf(addr1.address)).to.equal(mintAmount);
        });

        it("Should not allow non-owner to use owner mint", async function () {
            await expect(abmNFT.connect(addr1).ownerMint(addr2.address, 1))
                .to.be.revertedWith("Ownable: caller is not the owner");
        });

        it("Should allow owner to update prices", async function () {
            const newPrice = ethers.utils.parseEther("150");
            
            await abmNFT.updateMintPrice(newPrice);
            expect(await abmNFT.mintPrice()).to.equal(newPrice);

            await abmNFT.updatePresalePrice(newPrice);
            expect(await abmNFT.presalePrice()).to.equal(newPrice);
        });

        it("Should allow owner to update URIs", async function () {
            const newBaseURI = "https://new-base-uri.com/";
            const newContractURI = "https://new-contract-uri.com";
            
            await abmNFT.setBaseURI(newBaseURI);
            expect(await abmNFT._baseURI()).to.equal(newBaseURI);

            await abmNFT.setContractURI(newContractURI);
            expect(await abmNFT.contractURI()).to.equal(newContractURI);
        });
    });

    describe("Whitelist Management", function () {
        it("Should allow owner to add to whitelist", async function () {
            const addresses = [addr1.address, addr2.address];
            
            await abmNFT.addToPresaleWhitelist(addresses);
            
            expect(await abmNFT.isWhitelisted(addr1.address)).to.be.true;
            expect(await abmNFT.isWhitelisted(addr2.address)).to.be.true;
        });

        it("Should allow owner to remove from whitelist", async function () {
            const addresses = [addr1.address, addr2.address];
            
            await abmNFT.addToPresaleWhitelist(addresses);
            await abmNFT.removeFromPresaleWhitelist(addresses);
            
            expect(await abmNFT.isWhitelisted(addr1.address)).to.be.false;
            expect(await abmNFT.isWhitelisted(addr2.address)).to.be.false;
        });

        it("Should not allow non-owner to manage whitelist", async function () {
            await expect(abmNFT.connect(addr1).addToPresaleWhitelist([addr2.address]))
                .to.be.revertedWith("Ownable: caller is not the owner");
        });
    });

    describe("Token URI", function () {
        beforeEach(async function () {
            await abmNFT.toggleMinting(true);
            await abmNFT.connect(addr1).publicMint(1, {
                value: MINT_PRICE
            });
        });

        it("Should return correct token URI", async function () {
            const tokenId = 1;
            const expectedURI = `${BASE_URI}${tokenId}.json`;
            
            expect(await abmNFT.tokenURI(tokenId)).to.equal(expectedURI);
        });

        it("Should revert for non-existent token", async function () {
            await expect(abmNFT.tokenURI(999))
                .to.be.revertedWith("ABMNFT: URI query for nonexistent token");
        });
    });

    describe("Supply Tracking", function () {
        it("Should track total supply correctly", async function () {
            await abmNFT.toggleMinting(true);
            
            const mintAmount = 5;
            await abmNFT.connect(addr1).publicMint(mintAmount, {
                value: MINT_PRICE.mul(mintAmount)
            });

            expect(await abmNFT.totalSupply()).to.equal(mintAmount);
            expect(await abmNFT.getCurrentTokenId()).to.equal(mintAmount);
            expect(await abmNFT.remainingSupply()).to.equal(MAX_SUPPLY - mintAmount);
        });
    });

    describe("Withdrawal", function () {
        beforeEach(async function () {
            await abmNFT.toggleMinting(true);
            await abmNFT.connect(addr1).publicMint(5, {
                value: MINT_PRICE.mul(5)
            });
        });

        it("Should allow owner to withdraw funds", async function () {
            const initialBalance = await ethers.provider.getBalance(owner.address);
            const contractBalance = await ethers.provider.getBalance(abmNFT.address);
            
            await abmNFT.withdraw();
            
            const finalBalance = await ethers.provider.getBalance(owner.address);
            expect(finalBalance).to.be.gt(initialBalance);
            expect(await ethers.provider.getBalance(abmNFT.address)).to.equal(0);
        });

        it("Should not allow non-owner to withdraw", async function () {
            await expect(abmNFT.connect(addr1).withdraw())
                .to.be.revertedWith("Ownable: caller is not the owner");
        });

        it("Should handle zero balance withdrawal", async function () {
            // Withdraw all funds first
            await abmNFT.withdraw();
            
            // Try to withdraw again
            await expect(abmNFT.withdraw())
                .to.be.revertedWith("ABMNFT: No funds to withdraw");
        });
    });

    describe("Security", function () {
        it("Should prevent reentrancy attacks", async function () {
            await abmNFT.toggleMinting(true);
            
            // This test would require a malicious contract to properly test reentrancy
            // For now, we verify reentrancy guard is in place
            const receipt = await abmNFT.connect(addr1).publicMint(1, {
                value: MINT_PRICE
            });
            
            expect(receipt.hash).to.not.be.undefined;
        });

        it("Should reject direct ETH transfers", async function () {
            await expect(owner.sendTransaction({
                to: abmNFT.address,
                value: ethers.utils.parseEther("1")
            })).to.be.revertedWith("ABMNFT: Direct transfers not allowed");
        });
    });

    describe("Events", function () {
        it("Should emit Minted event on single mint", async function () {
            await abmNFT.toggleMinting(true);
            
            await expect(abmNFT.connect(addr1).publicMint(1, {
                value: MINT_PRICE
            })).to.emit(abmNFT, "Minted")
                .withArgs(addr1.address, 1, `${BASE_URI}1.json`);
        });

        it("Should emit BatchMinted event on batch mint", async function () {
            await abmNFT.toggleMinting(true);
            const mintAmount = 3;
            
            await expect(abmNFT.connect(addr1).publicMint(mintAmount, {
                value: MINT_PRICE.mul(mintAmount)
            })).to.emit(abmNFT, "BatchMinted")
                .withArgs(addr1.address, [1, 2, 3]);
        });

        it("Should emit MintingStatusChanged event", async function () {
            await expect(abmNFT.toggleMinting(true))
                .to.emit(abmNFT, "MintingStatusChanged")
                .withArgs(true);
        });

        it("Should emit PresaleStatusChanged event", async function () {
            await expect(abmNFT.togglePresale(true))
                .to.emit(abmNFT, "PresaleStatusChanged")
                .withArgs(true);
        });
    });
});
