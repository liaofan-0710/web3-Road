// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 现在轮到有趣的部分。 我们需要弄清楚交易的 data 部分。 这意味着我们需要找出我们将要调用的智能合约函数名，以及函数将接收的输入。
// 然后我们使用函数名的 keccak-256 哈希来检索 方法 ID，它是前 8 个字符（4 个字节）。
// 然后，我们附加我们发送的地址，并附加我们打算转账的代币数量。 这些输入需要 256 位长（32 字节）并填充左侧。 方法 ID 不需填充。
// 为了演示，我创造了一个新的代币(RCCDemoToken, RDT)，可以使用 Remix 在线工具编译合约之后，部署到 Sepolia 网络：

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyToken is ERC20, Ownable {
    uint256 public constant RATE = 100000000; // 100000000 MyToken per 1 ETH
    uint256 public constant MIN_ETH = 0.001 ether;

    constructor(address initialOwner) ERC20("RCCDemoToken", "RDT") Ownable(msg.sender) {
    }

    function mint() public payable {
        require(msg.value >= MIN_ETH, "Not enough ETH sent");
        uint256 tokensToMint = (msg.value * RATE);
        _mint(msg.sender, tokensToMint);
    }

    function withdrawETH() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No ETH to withdraw");
        payable(owner()).transfer(balance);
    }

    receive() external payable {
        mint();
    }
}