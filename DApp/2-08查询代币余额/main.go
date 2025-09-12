package main

import (
	token "Project/contracts_erc20" // for demo
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
)

// npm install -g solc
// go install github.com/ethereum/go-ethereum/...@latest
// solcjs --abi IERC20Metadata.sol
// abigen --abi=erc20_sol_ERC20.abi --pkg=token --out=erc20.go

// solcjs --abi IERC20Metadata.sol
// abigen --abi=IERC20Metadata_sol_IERC20Metadata.abi --pkg=token --out=erc20.go

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/t-lb3dGnfhk7pHbNlHHi1")
	if err != nil {
		log.Fatal(err)
	}
	// Golem (GNT) Address
	// 假设我们已经创建了以太坊客户端实例，将新的 token 包导入我们的项目，并实例化它。这个例子里我们用 RCCDemoToken 代币的地址
	tokenAddress := common.HexToAddress("0xfadea654ea83c00e5003d2ea15c59830b65471c0")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	// 我们可以调用任何 ERC20 的方法。 例如，我们可以查询用户的代币余额
	address := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	// 我们还可以读 ERC20 智能合约的公共变量
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("name: %s\n", name)         // "name: RCCDemoToken"
	fmt.Printf("symbol: %s\n", symbol)     // "symbol: RDT"
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"
	fmt.Printf("wei: %s\n", bal)           // "wei: 9996999999000000000000000"
	// 我们可以做一些简单的数学运算将余额转换为可读的十进制格式
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance: %f", value) // "balance: 9996999.999000"
}
