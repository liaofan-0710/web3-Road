require("@nomicfoundation/hardhat-toolbox");
require("hardhat-deploy");
require("@openzeppelin/hardhat-upgrades");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  namedAccounts: {
    deployer: 0,
    seller: {
      default: 1,
    },
    nftContract: {
      default: 1,
    },
    user: 1,
    user: 2,
  }
};
