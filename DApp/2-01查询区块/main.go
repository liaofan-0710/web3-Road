package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/P0AtFolBHAoBDO5Sq73Nn")
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(5671744)

	// 区块头: 调用客户端的 HeaderByNumber 来返回有关一个区块的头信息。若传入 nil，它将返回最新的区块头
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	fmt.Println(header.Number.Uint64())     // 5671744
	fmt.Println(header.Time)                // 1712798400
	fmt.Println(header.Difficulty.Uint64()) // 0
	fmt.Println(header.Hash().Hex())        // 0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5
	if err != nil {
		log.Fatal(err)
	}

	// 完整区块：调用客户端的 BlockByNumber 方法来获得完整区块。您可以读取该区块的所有内容和元数据，例如，区块号，区块时间戳，区块摘要，区块难度以及交易列表等等。
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Time())                // 1712798400
	fmt.Println(block.Difficulty().Uint64()) // 0
	fmt.Println(block.Hash().Hex())          // 0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5
	fmt.Println(len(block.Transactions()))   // 70

	// 调用 TransactionCount 只返回一个区块的交易数目
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count) // 70
}
