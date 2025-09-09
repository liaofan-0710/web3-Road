# NFT 拍卖项目

## 项目介绍

技术栈：
```shell
开发语言：solidity JavaScript
开发框架：Hardhat
工具：npm npx nvm
```

安装依赖
```shell
npm install hardhat@2.26.3
npm install @openzeppelin/contracts @chainlink/contracts @nomiclabs/hardhat-ethers hardhat-deploy
```
本地开发网络
```shell
1. 启动本地Hardhat网络：
npx hardhat node

2. 在新终端中部署合约：
npx hardhat deploy --network localhost
npx hardhat deploy --tags xxxx --network localhost
```

测试网部署
```shell
创建.env
SEPOLIA_RPC_URL=<your-sepolia-rpc-url>
PRIVATE_KEY=<your-private-key>

部署到Sepolia测试网
npx hardhat deploy --network sepolia
```
测试
```shell
运行测试套件：
npx hardhat test
npx hardhat test test/xx.js

运行带gas报告的测试
REPORT_GAS=true npx hardhat test
```

## 项目结构

```
.
├── contracts/
│   ├── NftAuction.sol                   # 拍卖合约
│   ├── NftAuctionFactory.sol            # 工厂合约
│   ├── NftAuctionV2.sol                 # 升级版拍卖合约
│   └── NftAuctionTwo.sol                # 改造版拍卖合约
│
├── deploy/
│   ├── 01_deploy_nft_auction.js         # 部署NftAuction合约
│   ├── 02_upgrade_nft_auction.js        # 升级NftAuction合约
│   └── 03_deploy_nft_auction_factory.js # 部署NftAuctionFactory工厂合约
│
├── test/
│   ├── 02_upgrade_nft_auction.js        # 升级NftAuction合约
│   └── auction.test.js                  # 拍卖功能测试脚本
│
└── hardhat.config.js                    # Hardhat配置文件
```

## 功能介绍
  ✅ ERC721 标准（NFT 转移 + 铸造）

  ✅ 拍卖合约(创建拍卖, 出价, 结束拍卖)

  ✅ 工厂模式（每个拍卖独立合约）
  
  ✅ Chainlink 预言机（价格换算）

  ✅ 可升级（透明代理模式）

  ❌ 跨链拍卖 (NFT 跨链拍卖, 不同链上参与拍卖)

  ✅ 合约升级 (使用 UUPS 或透明代理模式实现合约升级)

  ✅ 测试与部署 (编写单元测试和集成测试，覆盖所有功能, 使用 Hardhat 部署脚本)

## 未来TODO
1. 拍卖合约 (负责处理拍卖逻辑、出价、以及通过 CCIP 接收和发送消息与资产)
2. 工厂合约 (负责上架/铸造nft)
3. CCIP 消息集成 (发送端： 源链的合约会通过 CCIP 发送一条消息到目标链的拍卖合约,包含出价者的地址、出价金额。 接收端： 执行相应的拍卖逻辑，更新最高出价、记录出价者)
