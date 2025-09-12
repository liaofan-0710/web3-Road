package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// client, err := ethclient.Dial("https://rinkeby.infura.io")
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/xxx")
	if err != nil {
		log.Fatal(err)
	}

	// 转账交易包括打算转账的以太币数量，燃气限额，燃气价格，一个自增数(nonce)，接收地址以及可选择性的添加的数据。 在发送到以太坊网络之前，必须使用发送方的私钥对该交易进行签名。
	// 假设我们已经连接了客户端，下一步就是加载私钥。
	privateKey, err := crypto.HexToECDSA("siyao")
	if err != nil {
		log.Fatal(err)
	}

	// 之后我们需要获得帐户的随机数(nonce)。 每笔交易都需要一个 nonce。 根据定义，nonce 是仅使用一次的数字。
	// 如果是发送交易的新帐户，则该随机数将为“0”。 来自帐户的每个新事务都必须具有前一个 nonce 增加 1 的 nonce。
	// 很难对所有 nonce 进行手动跟踪，于是 ethereum 客户端提供一个帮助方法 PendingNonceAt，它将返回你应该使用的下一个 nonce。

	// 该函数需要我们发送的帐户的公共地址 - 这个我们可以从私钥派生
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 接下来我们可以读取我们应该用于帐户交易的随机数。
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 下一步是设置我们将要转移的 ETH 数量。 但是我们必须将 ETH 以太转换为 wei，因为这是以太坊区块链所使用的。
	// 以太网支持最多 18 个小数位，因此 1 个 ETH 为 1 加 18 个零。
	value := big.NewInt(10000000000000000) // in wei (1 eth)
	// ETH 转账的燃气应设上限为“21000”单位。
	gasLimit := uint64(21000) // in units
	// 然而，燃气价格总是根据市场需求和用户愿意支付的价格而波动的，因此对燃气价格进行硬编码有时并不理想。
	// go-ethereum 客户端提供 SuggestGasPrice 函数，用于根据'x'个先前块来获得平均燃气价格。
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 接下来我们弄清楚我们将 ETH 发送给谁。
	toAddress := common.HexToAddress("0xf57599B35863f7421E9E06195fB5b4AAf4E9108d")
	var data []byte
	// 现在我们最终可以通过导入 go-ethereum core/types 包并调用 NewTransaction 来生成我们的未签名以太坊事务，这个函数需要接收 nonce，地址，值，燃气上限值，燃气价格和可选发的数据。
	// 发送 ETH 的数据字段为“nil”。 在与智能合约进行交互时，我们将使用数据字段，仅仅转账以太币是不需要数据字段的。
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	// 下一步是使用发件人的私钥对事务进行签名。 为此，我们调用 SignTx 方法，该方法接受一个未签名的事务和我们之前构造的私钥。
	// SignTx 方法需要 EIP155 签名者，这个也需要我们先从客户端拿到链 ID。
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 现在通过在 client 实例调用 SendTransaction 来将已签名的事务广播到整个网络。
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
