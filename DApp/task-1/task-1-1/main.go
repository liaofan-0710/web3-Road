package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

// ## 任务 1：区块链读写 任务目标
// 使用 Sepolia 测试网络实现基础的区块链交互，包括查询区块和发送交易。
//  具体任务
// 1. 环境搭建
//    - 安装必要的开发工具，如 Go 语言环境、 go-ethereum 库。
//    - 注册 Infura 账户，获取 Sepolia 测试网络的 API Key。
// 2. 查询区块
//    - 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//    - 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
//    - 输出查询结果到控制台。
// 3. 发送交易
//    - 准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
//    - 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//    - 构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
//    - 对交易进行签名，并将签名后的交易发送到网络。
//    - 输出交易的哈希值。

func main() {
	// 查询区块
	// 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/t-lb3dGnfhk7pHbNlHHi1")
	if err != nil {
		log.Fatal(err)
	}
	//blockIndex := big.NewInt(0)
	var blockIndex string
	fmt.Println("输入区块号：") //  5671744
	_, err = fmt.Scanln(&blockIndex)
	if err != nil {
		fmt.Println("blockIndex get fail", err)
	}
	blockValue := new(big.Int)
	_, ok := blockValue.SetString(blockIndex, 10) // 10 表示十进制
	if !ok {
		// 处理转换失败的情况
		fmt.Println("bigInt error:", ok)
	}
	block, err := client.BlockByNumber(context.Background(), blockValue)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块的哈希：", block.Hash().Hex())
	fmt.Println("时间戳：", block.Time())
	fmt.Println("交易数量", len(block.Transactions()))
	client.Close()

	// 发送交易
	// 准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥
	privateKey, err := crypto.HexToECDSA("195e42819e963554f63765fe4ed4ff472c38b426c8ca6dfe2e2b1eac8ee982f8")
	if err != nil {
		log.Fatal(err)
	}

	// 编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络
	client, err = ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/t-lb3dGnfhk7pHbNlHHi1")
	if err != nil {
		log.Fatal(err)
	}
	// 构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress("0xf57599B35863f7421E9E06195fB5b4AAf4E9108d")

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(1000000000000000) // in wei (1 eth) 以太网支持最多 18 个小数位，因此 1 个 ETH 为 1 加 18 个零。 0.001
	gasLimit := uint64(21000)             // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	// 对交易进行签名，并将签名后的交易发送到网络
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 对交易进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 发送到网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 输出交易的哈希值
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // 0x04aaea92f1e500bec5f13d5f075bfb05c7a834eb31914af56fe2b189a729b23d
}
