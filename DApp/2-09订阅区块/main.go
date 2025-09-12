package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// 订阅区块需要 websocket RPC URL
	//client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/t-lb3dGnfhk7pHbNlHHi1")
	if err != nil {
		log.Fatal(err)
	}

	// 接下来，我们将创建一个新的通道，用于接收最新的区块头
	headers := make(chan *types.Header)
	// 现在我们调用客户端的 SubscribeNewHead 方法，它接收我们刚创建的区块头通道，该方法将返回一个订阅对象
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	// 订阅将推送新的区块头事件到我们的通道，因此我们可以使用一个 select 语句来监听新消息。
	// 订阅对象还包括一个 error 通道，该通道将在订阅失败时发送消息。
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			// 要获得该区块的完整内容，我们可以将区块头的摘要传递给客户端的 BlockByHash 函数
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(block.Hash().Hex())      // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64()) // 3477413
			//fmt.Println(block.Time().Uint64())     // 1529525947
			fmt.Println(block.Time())              // 1529525947
			fmt.Println(block.Nonce())             // 130524141876765836
			fmt.Println(len(block.Transactions())) // 7
		}
	}
}
