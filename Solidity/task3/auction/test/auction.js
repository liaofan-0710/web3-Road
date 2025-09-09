const { ethers, deployments } = require('hardhat')
const { expect } = require("chai")

describe("Test auction", function () {
    it("Should be ok", async function () {
        await main()
    });
});

async function main() {
    const [ signer, buyer ] = await ethers.getSigners()
    await deployments.fixture(['deployNftAuction'])
    const nftAuctionProxy = await deployments.get('NftAuctionProxy')

    console.log("nftAuctionProxy::", nftAuctionProxy);

    // 1. 部署 ERC721
    const TestERC721 = await ethers.getContractFactory("TestERC721")
    const testERC721 = await TestERC721.deploy();
    await testERC721.waitForDeployment();
    const testERC721Address = await testERC721.getAddress();

    // mint 10个 NFT
    for (let i = 0; i < 10; i++) {
        await testERC721.mint(signer.address, i + 1);
    }

    const tokenId = 1;
    // 2. 调用 createAuction 方法创建
    const nftAuction = await ethers.getContractAt(
        "NftAuction", 
        nftAuctionProxy.address
    );

    // 给代理合约授权
    await testERC721.connect(signer).setApprovalForAll(nftAuctionProxy.address, true);

    await nftAuction.createAuction(
        10, 
        ethers.parseEther("0.01"), 
        testERC721Address,
        tokenId
    );

    const auction = await nftAuction.auctions(0);

    console.log("创建拍卖成功::", auction);

    // 3. 购买者参与拍卖
    // await testERC721.connect(signer).approve(nftAuction.address, tokenId);
    // ETH参与竞价

    await nftAuction.connect(buyer).placeBid(0, { value: ethers.parseEther("0.01") });

    // 4. 结束拍卖
    // 等待 10 s
    await new Promise((resolve) => setTimeout(resolve, 10 * 1000));

    await nftAuction.connect(signer).endAuction(0);

    // 验证结果
    const auctionResult = await nftAuction.auctions(0);
    console.log("结束拍卖后读取拍卖成功::", auctionResult);
    expect(auctionResult.highestBidder).to.equal(buyer.address);
    expect(auctionResult.highestBid).to.equal(ethers.parseEther("0.01"));

    // 验证 NFT 所有权
    const owner = await testERC721.ownerOf(tokenId);
    console.log("owner::", owner);
    expect(owner).to.equal(buyer.address);
}

main()