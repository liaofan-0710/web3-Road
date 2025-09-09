const { deployments, upgrades, ethers } = require("hardhat");
const path = require('path');
const fs = require('fs');

module.exports = async ({ getNamedAccounts, deployments }) => {
  const { save } = deployments;
  const { deployer, seller, nftContract } = await getNamedAccounts();

  console.log("部署用户地址:", deployer, seller, nftContract);
  const NftAuction = await ethers.getContractFactory("NftAuctionFactory");

  // 通过代理合约部署
  const nftAuctionProxy = await upgrades.deployProxy(NftAuction, [seller,86400, ethers.parseEther("1.0"), nftContract, 1, deployer, 10], { initializer: "initialize" });

  await nftAuctionProxy.waitForDeployment();

  const proxyAddress = await nftAuctionProxy.getAddress();
  console.log("代理合约地址:", proxyAddress);
  const implAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);
  console.log("实现合约地址: ", implAddress);

  const storePath = path.resolve(__dirname, "./.cache/proxyNftAuction.json");

  fs.writeFileSync(
    storePath, 
    JSON.stringify({
        proxyAddress,
        implAddress,
        abi: NftAuction.interface.format("json"),
    })
  );

  await save("NftAuctionProxy", {
    abi: NftAuction.interface.format("json"),
    address: proxyAddress,
  })

};
// add tags and dependencies
module.exports.tags = ["deployNftAuctionFactory"];