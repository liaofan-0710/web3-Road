package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/")
	if err != nil {
		log.Fatal(err)
	}

	// 假设我们已经连接了客户端，下一步就是加载私钥。
	privateKey, err := crypto.HexToECDSA("siyao")
	if err != nil {
		log.Fatal(err)
	}

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

	// 代币传输不需要传输 ETH，因此将交易“值”设置为“0”。
	value := big.NewInt(0) // in wei (0 eth)
	// 然而，燃气价格总是根据市场需求和用户愿意支付的价格而波动的，因此对燃气价格进行硬编码有时并不理想。
	// go-ethereum 客户端提供 SuggestGasPrice 函数，用于根据'x'个先前块来获得平均燃气价格。
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 先将您要发送代币的地址存储在变量中。
	toAddress := common.HexToAddress("0x12c8c51Cc6a1c7E6b53E4b91e4c5C4D6D82418D0")
	// 让我们将代币合约地址分配给变量。
	tokenAddress := common.HexToAddress("0x23331964e6Df6Ca64d7EE406b22EF59D77f08D45")
	// 函数名将是传递函数的名称，即 ERC-20 规范中的 transfer 和参数类型。
	// 第一个参数类型是 address（令牌的接收者），第二个类型是 uint256（要发送的代币数量）。
	// 不需要没有空格和参数名称。 我们还需要用字节切片格式。
	transferFnSignature := []byte("transfer(address,uint256)")
	// 我们现在将从 go-ethereum 导入 crypto/sha3 包以生成函数签名的 Keccak256 哈希。
	// 然后我们只使用前 4 个字节来获取方法 ID。
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb
	// 接下来，我们需要将给我们发送代币的地址左填充到 32 字节。
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x00000000000000000000000023331964e6df6ca64d7ee406b22ef59d77f08d45
	// 接下来我们确定要发送多少个代币，在这个例子里是 1,000 个，并且我们需要在 big.Int 中格式化为 wei。
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens
	// 代币量也需要左填充到 32 个字节。
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000
	// 接下来我们只需将方法 ID，填充后的地址和填后的转账量，接到将成为我们数据字段的字节片。
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// 燃气上限制将取决于交易数据的大小和智能合约必须执行的计算步骤。 幸运的是，客户端提供了 EstimateGas 方法，它可以为我们估算所需的燃气量。
	// 这个函数从 ethereum 包中获取 CallMsg 结构，我们在其中指定数据和地址。 它将返回我们估算的完成交易所需的估计燃气上限。
	//gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	//	To:   &toAddress,
	//	Data: data,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(gasLimit) // 23256
	gasLimit := uint64(60000)
	// 接下来我们需要做的是构建交易类型，这类似于在 ETH 转账部分中看到的，除了_to_字段将是代币智能合约地址。
	// 这个常让人困惑。我们还必须在调用中包含 0 ETH 的值字段和刚刚生成的数据字节。
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// 下一步是使用发件人的私钥对事务进行签名。 SignTx 方法需要 EIP155igner，需要我们先从客户端拿到链 ID。
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 最后发送交易。
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc
}
