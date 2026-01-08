#!/usr/bin/env python3
"""
ABM Diora NFT Metadata Generator
Generates metadata for 15 unique NFTs representing ABM Diora ecosystem
"""

import json
import os

# NFT Categories and their properties
NFT_CATEGORIES = {
    1: {
        "name": "Genesis Block",
        "category": "Foundation",
        "rarity": "Legendary",
        "component": "Genesis",
        "significance": "Network Birth",
        "color": "Blue Gradient",
        "animation": "Pulse Effect"
    },
    2: {
        "name": "Hybrid PoS",
        "category": "Consensus",
        "rarity": "Epic",
        "component": "Consensus Mechanism",
        "significance": "Security & Performance",
        "color": "Purple Gradient",
        "animation": "Rotate Effect"
    },
    3: {
        "name": "EVM Compatibility",
        "category": "Virtual Machine",
        "rarity": "Epic",
        "component": "Smart Contracts",
        "significance": "Developer Experience",
        "color": "Green Gradient",
        "animation": "Code Flow"
    },
    4: {
        "name": "Enterprise Security",
        "category": "Security",
        "rarity": "Legendary",
        "component": "Security Framework",
        "significance": "Institutional Trust",
        "color": "Red Gradient",
        "animation": "Shield Effect"
    },
    5: {
        "name": "1000+ TPS",
        "category": "Performance",
        "rarity": "Rare",
        "component": "Throughput",
        "significance": "High Performance",
        "color": "Orange Gradient",
        "animation": "Speed Lines"
    },
    6: {
        "name": "DIO Token",
        "category": "Economics",
        "rarity": "Epic",
        "component": "Native Token",
        "significance": "Economic Engine",
        "color": "Gold Gradient",
        "animation": "Coin Spin"
    },
    7: {
        "name": "6-Second Blocks",
        "category": "Performance",
        "rarity": "Rare",
        "component": "Block Time",
        "significance": "Fast Finality",
        "color": "Cyan Gradient",
        "animation": "Clock Pulse"
    },
    8: {
        "name": "42 Validators",
        "category": "Network",
        "rarity": "Rare",
        "component": "Validator Network",
        "significance": "Decentralization",
        "color": "Indigo Gradient",
        "animation": "Network Nodes"
    },
    9: {
        "name": "Developer Tools",
        "category": "Ecosystem",
        "rarity": "Common",
        "component": "Development Suite",
        "significance": "Developer Experience",
        "color": "Teal Gradient",
        "animation": "Tool Animation"
    },
    10: {
        "name": "Layer 2 Ready",
        "category": "Scalability",
        "rarity": "Rare",
        "component": "Scaling Solution",
        "significance": "Future Growth",
        "color": "Pink Gradient",
        "animation": "Layer Effect"
    },
    11: {
        "name": "Cross-Chain Bridge",
        "category": "Interoperability",
        "rarity": "Epic",
        "component": "Bridge Protocol",
        "significance": "Network Connection",
        "color": "Yellow Gradient",
        "animation": "Bridge Animation"
    },
    12: {
        "name": "DeFi Integration",
        "category": "Finance",
        "rarity": "Rare",
        "component": "DeFi Protocol",
        "significance": "Financial Services",
        "color": "Emerald Gradient",
        "animation": "Chart Animation"
    },
    13: {
        "name": "Governance Model",
        "category": "Governance",
        "rarity": "Epic",
        "component": "DAO Structure",
        "significance": "Community Control",
        "color": "Violet Gradient",
        "animation": "Vote Animation"
    },
    14: {
        "name": "API Gateway",
        "category": "Infrastructure",
        "rarity": "Common",
        "component": "API Layer",
        "significance": "Data Access",
        "color": "Gray Gradient",
        "animation": "Data Flow"
    },
    15: {
        "name": "Future Vision",
        "category": "Roadmap",
        "rarity": "Legendary",
        "component": "Future State",
        "significance": "2027+ Roadmap",
        "color": "Rainbow Gradient",
        "animation": "Future Vision"
    }
}

def generate_metadata(token_id, properties):
    """Generate metadata for a specific token ID"""
    return {
        "name": f"ABM Diora #{token_id} - {properties['name']}",
        "description": f"Unique NFT representing {properties['name']} from the ABM Diora blockchain ecosystem. This {properties['rarity'].lower()} NFT symbolizes {properties['significance'].lower()} in our next-generation Layer 1 solution.",
        "image": f"https://diorafund.github.io/diora-blockchain/images/nft/{token_id}.png",
        "external_url": "https://diorafund.github.io/diora-whitepaper/",
        "attributes": [
            {
                "trait_type": "Type",
                "value": properties["name"]
            },
            {
                "trait_type": "Category",
                "value": properties["category"]
            },
            {
                "trait_type": "Rarity",
                "value": properties["rarity"]
            },
            {
                "trait_type": "Blockchain Component",
                "value": properties["component"]
            },
            {
                "trait_type": "Significance",
                "value": properties["significance"]
            },
            {
                "trait_type": "Color Scheme",
                "value": properties["color"]
            },
            {
                "trait_type": "Animation",
                "value": properties["animation"]
            }
        ],
        "background_color": "1e40af",
        "animation_url": f"https://diorafund.github.io/diora-blockchain/animations/nft/{token_id}.mp4"
    }

def main():
    """Generate all metadata files"""
    # Create metadata directory if it doesn't exist
    os.makedirs("metadata", exist_ok=True)
    
    # Generate metadata for all 15 NFTs
    for token_id in range(1, 16):
        if token_id in NFT_CATEGORIES:
            metadata = generate_metadata(token_id, NFT_CATEGORIES[token_id])
            
            # Write metadata to file
            filename = f"metadata/{token_id}.json"
            with open(filename, 'w') as f:
                json.dump(metadata, f, indent=2)
            
            print(f"Generated metadata for token #{token_id}: {NFT_CATEGORIES[token_id]['name']}")
    
    print(f"\nGenerated {len(NFT_CATEGORIES)} metadata files successfully!")
    
    # Generate summary
    rarity_count = {}
    for props in NFT_CATEGORIES.values():
        rarity = props["rarity"]
        rarity_count[rarity] = rarity_count.get(rarity, 0) + 1
    
    print("\nRarity Distribution:")
    for rarity, count in rarity_count.items():
        print(f"  {rarity}: {count}")

if __name__ == "__main__":
    main()
