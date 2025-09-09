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

func main() {
	client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/P0AtFolBHAoBDO5Sq73Nn")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("fd9228e4d98f0375cecce4bf9c1460b2bd91740d6d03d18bd3513eaee2605c84")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 要转移的 ETH 数量
	value := big.NewInt(0) // in wei (1 eth)
	// ETH 转账的燃气应设上限为“21000”单位。
	gasLimit := uint64(21000)
	// gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	gasPrice := big.NewInt(0) // in wei (30 gwei)

	// 根据'x'个先前块来获得平均燃气价格
	//gasPrice, err = client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(gasPrice, err)
	// 将 ETH 发送给谁
	toAddress := common.HexToAddress("0xf57599B35863f7421E9E06195fB5b4AAf4E9108d")

	// 生成我们的未签名以太坊事务
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 下一步是使用发件人的私钥对事务进行签名。 为此，我们调用 SignTx 方法，该方法接受一个未签名的事务和我们之前构造的私钥。 SignTx 方法需要 EIP155 签名者，这个也需要我们先从客户端拿到链 ID
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 通过在 client 实例调用 SendTransaction 来将已签名的事务广播到整个网络
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0x6abe4a5dcb8e28d58825935eddb9d5c8c9a885be3968aa846ff062d016a9f1d3

}
