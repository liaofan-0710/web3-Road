package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/")
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress("0x12c8c51Cc6a1c7E6b53E4b91e4c5C4D6D82418D0")
	// 调用 ethclient 的 BalanceAt 方法，给它传递账户地址和可选的区块号。将区块号设置为 nil 将返回最新的余额
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("balance", balance) // 2217034205429474656
	// 传区块高度能读取指定区块时的账户余额，区块高度必须是 big.Int 类型
	blockNumber := big.NewInt(5532993)
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("balanceAt", balanceAt) // 0
	// 以太坊中的数字是使用尽可能小的单位来处理的，因为它们是定点精度，在 ETH 中它是_wei_。要读取 ETH 值，
	// 您必须做计算 wei/10^18。因为我们正在处理大数，我们得导入原生的 Go math 和 math/big 包。这是您做的转换
	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("ethValue", ethValue) // 0
	// 有时您想知道待处理的账户余额是多少，例如，在提交或等待交易确认后。
	// 客户端提供了类似 BalanceAt 的方法，名为 PendingBalanceAt，它接收账户地址作为参数。
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println("pendingBalance", pendingBalance) // 2217034205429474656
}
