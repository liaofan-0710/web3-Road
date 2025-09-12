package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/<API_KEY>")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x2958d15bc5b64b11Ec65e623Ac50C198519f8742")
	// 智能合约可以可选地释放“事件”，其作为交易收据的一部分存储日志。读取这些事件相当简单。
	// 首先我们需要构造一个过滤查询。我们从 go-ethereum 包中导入 FilterQuery 结构体并用过滤选项初始化它。
	// 我们告诉它我们想过滤的区块范围并指定从中读取此日志的合约地址。
	// 在示例中，我们将从在智能合约章节创建的智能合约中读取特定区块所有日志。
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6920583),
		// ToBlock:   big.NewInt(2394201),
		Addresses: []common.Address{
			contractAddress,
		},
		// Topics: [][]common.Hash{
		//  {},
		//  {},
		// },
	}

	// 下一步是调用 ethclient 的 FilterLogs，它接收我们的查询并将返回所有的匹配事件日志
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// 返回的所有日志将是 ABI 编码，因此它们本身不会非常易读。为了解码日志，我们需要导入我们智能合约的 ABI。
	// 为此，我们导入编译好的智能合约 Go 包，它将包含名称格式为 <Contract>ABI 的外部属性。
	// 之后，我们使用 go-ethereum 中的 accounts/abi 包的 abi.JSON 函数返回一个我们可以在 Go 应用程序中使用的解析过的 ABI 接口。
	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}

	// 现在我们可以通过日志进行迭代并将它们解码为我么可以使用的类型。
	// 若您回忆起我们的样例合约释放的日志在 Solidity 中是类型为 bytes32，那么 Go 中的等价物将是 [32]byte。
	// 我们可以使用这些类型创建一个匿名结构体，并将指针作为第一个参数传递给解析后的 ABI 接口的 Unpack 函数，以解码原始的日志数据。
	// 第二个参数是我们尝试解码的事件名称，最后一个参数是编码的日志数据。
	for _, vLog := range logs {
		// 此外，日志结构体包含附加信息，例如，区块摘要，区块号和交易摘要。
		fmt.Println(vLog.BlockHash.Hex())
		fmt.Println(vLog.BlockNumber)
		fmt.Println(vLog.TxHash.Hex())
		// 在 solidity 中声明事件时，在类型与参数名称之间添加 indexed 关键字，来标记可索引 top
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(common.Bytes2Hex(event.Key[:]))
		fmt.Println(common.Bytes2Hex(event.Value[:]))
		var topics []string
		for i := range vLog.Topics {
			topics = append(topics, vLog.Topics[i].Hex())
		}

		// 第一个主题总是事件的签名。我们的示例合约不包含可索引的事件，但如果它确实包含，这是如何读取事件主题。
		fmt.Println("topics[0]=", topics[0])
		if len(topics) > 1 {
			fmt.Println("indexed topics:", topics[1:])
		}
	}

	// 首个主题只是被哈希过的事件签名。
	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println("signature topics=", hash.Hex())
}
