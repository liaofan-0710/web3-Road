const {ethers, upgrades} = require("hardhat");
const fs = require("fs")
const path = require("path")

module.exports = async ({ getNamedAccounts, deployments }) => {
  const { save } = deployments;
  const { deployer } = await getNamedAccounts();

  console.log("部署用户地址:", deployer);
  
  // 读取 .cache/proxyNftAuction.json 文件
  const storePath = path.resolve(__dirname, "./.cache/proxyNftAuction.json");
  const storeData = fs.readFileSync(storePath, "utf-8");
  const { proxyAddress, implAddress, abi } = JSON.parse(storeData);

  // 升级版的业务合约
  const NftAuctionV2 = await ethers.getContractFactory("NftAuctionV2");

  // 升级合约
  const nftAuctionProxy = await upgrades.upgradeProxy(proxyAddress, NftAuctionV2);

  await nftAuctionProxy.waitForDeployment();

  const proxyAddressV2 = await nftAuctionProxy.getAddress();

  // 保存代理合约地址
//   fs.writeFileSync(
//     storePath, JSON.stringify({
//       proxyAddress: proxyAddressV2,
//       implAddress: implAddress,
//       abi: abi
//     }, null, 2)
//   );

  await save("NftAuctionProxyV2", {
    address: proxyAddressV2,
    // implAddress: implAddress,
    abi: abi
  });
}

module.exports.tags = ["upgradeNftAuction"];