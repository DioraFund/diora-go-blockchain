/**
 * ABM Diora NFT Deployment Script
 * Deploys ABMNFT contract to ABM Diora blockchain
 */

const { ethers } = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
    console.log("🚀 Deploying ABM Diora NFT Collection...");
    
    // Get deployer account
    const [deployer] = await ethers.getSigners();
    console.log("📝 Deploying contracts with account:", deployer.address);
    
    // Deployment parameters
    const name = "ABM Foundation Genesis Collection";
    const symbol = "ABMNFT";
    const baseTokenURI = "https://diorafund.github.io/diora-blockchain/metadata/";
    const contractURI = "https://diorafund.github.io/diora-blockchain/metadata/collection.json";
    const teamWallet = deployer.address; // Change to actual team wallet
    
    console.log("📋 Contract Parameters:");
    console.log("  Name:", name);
    console.log("  Symbol:", symbol);
    console.log("  Base URI:", baseTokenURI);
    console.log("  Contract URI:", contractURI);
    console.log("  Team Wallet:", teamWallet);
    
    try {
        // Deploy contract
        console.log("🔨 Deploying ABMNFT contract...");
        const ABMNFT = await ethers.getContractFactory("ABMNFT");
        const abmNFT = await ABMNFT.deploy(
            name,
            symbol,
            baseTokenURI,
            contractURI,
            teamWallet
        );
        
        console.log("⏳ Waiting for deployment confirmation...");
        await abmNFT.deployed();
        
        const contractAddress = abmNFT.address;
        console.log("✅ ABMNFT deployed successfully!");
        console.log("📍 Contract Address:", contractAddress);
        
        // Save deployment info
        const deploymentInfo = {
            network: "diora",
            contractName: "ABMNFT",
            contractAddress: contractAddress,
            deployer: deployer.address,
            deploymentTime: new Date().toISOString(),
            transactionHash: abmNFT.deployTransaction.hash,
            gasUsed: abmNFT.deployTransaction.gasLimit.toString(),
            gasPrice: abmNFT.deployTransaction.gasPrice.toString(),
            parameters: {
                name,
                symbol,
                baseTokenURI,
                contractURI,
                teamWallet
            }
        };
        
        // Save to file
        const deploymentsDir = path.join(__dirname, "deployments");
        if (!fs.existsSync(deploymentsDir)) {
            fs.mkdirSync(deploymentsDir, { recursive: true });
        }
        
        const deploymentFile = path.join(deploymentsDir, "ABMNFT.json");
        fs.writeFileSync(deploymentFile, JSON.stringify(deploymentInfo, null, 2));
        
        console.log("💾 Deployment info saved to:", deploymentFile);
        
        // Verify contract (if verification is supported)
        console.log("🔍 Verifying contract...");
        try {
            await hre.run("verify:verify", {
                address: contractAddress,
                constructorArguments: [
                    name,
                    symbol,
                    baseTokenURI,
                    contractURI,
                    teamWallet
                ],
            });
            console.log("✅ Contract verified successfully!");
        } catch (error) {
            console.log("⚠️ Contract verification failed:", error.message);
        }
        
        // Post-deployment setup
        console.log("⚙️ Setting up contract...");
        
        // Set minting price (100 DIO = 100 * 10^18 wei)
        const mintPrice = ethers.utils.parseEther("100");
        await abmNFT.updateMintPrice(mintPrice);
        console.log("💰 Mint price set to 100 DIO");
        
        // Set presale price (50 DIO)
        const presalePrice = ethers.utils.parseEther("50");
        await abmNFT.updatePresalePrice(presalePrice);
        console.log("💰 Presale price set to 50 DIO");
        
        // Add team wallet to presale whitelist
        await abmNFT.addToPresaleWhitelist([teamWallet]);
        console.log("👥 Team wallet added to presale whitelist");
        
        console.log("🎉 Contract setup complete!");
        
        // Display summary
        console.log("\n📊 Deployment Summary:");
        console.log("=".repeat(50));
        console.log("Contract Name: ABMNFT");
        console.log("Collection: ABM Foundation Genesis Collection");
        console.log("Total Supply: 15 NFTs");
        console.log("Mint Price: 100 DIO");
        console.log("Presale Price: 50 DIO");
        console.log("Max per Wallet: 5 NFTs");
        console.log("Max per Transaction: 3 NFTs");
        console.log("Contract Address:", contractAddress);
        console.log("Explorer: https://diorafund.github.io/diora-blockchain/explorer/address/" + contractAddress);
        console.log("=".repeat(50));
        
    } catch (error) {
        console.error("❌ Deployment failed:", error);
        process.exit(1);
    }
}

// Execute deployment
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error("❌ Error:", error);
        process.exit(1);
    });
