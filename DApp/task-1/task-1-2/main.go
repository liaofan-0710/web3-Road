package main

// ## 任务 2：合约代码生成 任务目标
// 使用 abigen 工具自动生成 Go 绑定代码，用于与 Sepolia 测试网络上的智能合约进行交互。
//  具体任务
// 1. 编写智能合约
//    - 使用 Solidity 编写一个简单的智能合约，例如一个计数器合约。
//    - 编译智能合约，生成 ABI 和字节码文件。
// 2. 使用 abigen 生成 Go 绑定代码
//    - 安装 abigen 工具。
//    - 使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码。
// 3. 使用生成的 Go 绑定代码与合约交互
//    - 编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约。
//    - 调用合约的方法，例如增加计数器的值。
//    - 输出调用结果。

import (
	"Project/count"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

const (
	contractAddr = "0x2d977c3C8ecE0A5740A0a95ADB1B30eC2E9a3cA2"
)

func main() {
	// 1.1 使用 Solidity 编写一个简单的智能合约，例如一个计数器合约 -- ./count/Count.sol

	// 1.2 编译智能合约，生成 ABI 和字节码文件
	// 1.2.1 solcjs --bin Count.sol
	// 1.2.2 solcjs --abi Count.sol
	// 1.2.3 使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码
	// go install github.com/ethereum/go-ethereum/...@latest
	// abigen --abi=Count_sol_SimpleCounter.abi --pkg=count --out=count.go

	// 3.1 使用生成的 Go 绑定代码与合约交互
	// 编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/t-lb3dGnfhk7pHbNlHHi1")
	if err != nil {
		log.Fatal(err)
	}
	// 创建合约实例
	countContract, err := count.NewCount(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	// 根据 hex 创建私钥实例
	privateKey, err := crypto.HexToECDSA("195e42819e963554f63765fe4ed4ff472c38b426c8ca6dfe2e2b1eac8ee982f8")
	if err != nil {
		log.Fatal(err)
	}
	// 调用合约方法
	// 准备数据
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("demo_save_key"))
	copy(value[:], []byte("demo_save_value11111"))

	// // 初始化交易opt实例
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal("NewKeyedTransactorWithChainID err:", err)
	}

	// 调用合约的方法，例如增加计数器的值
	increment, err := countContract.Increment(opt)
	if err != nil {
		log.Fatal("increment err:", err)
	}

	callOpt := &bind.CallOpts{Context: context.Background()}
	valueInContract, err := countContract.GetCount(callOpt)
	if err != nil {
		log.Fatal("GetCount err:", err)
	}

	// 输出调用结果。
	fmt.Println("increment:", increment)
	fmt.Println("valueInContract:", valueInContract)

}
