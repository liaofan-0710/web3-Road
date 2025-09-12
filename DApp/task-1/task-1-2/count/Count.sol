// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title 简单计数器合约
 * @dev 这是一个基础的计数器合约，允许增加、减少和重置计数
 */
contract SimpleCounter {
    // 存储当前计数值
    uint256 private count;

    // 记录计数最后一次被修改的区块号
    uint256 public lastModifiedBlock;

    // 合约部署者地址
    address public owner;

    // 事件：当计数发生变化时触发
    event CountChanged(uint256 newValue, address changedBy, uint256 blockNumber);

    // 修饰器：限制只有合约所有者可以调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only the contract owner can call this function");
        _;
    }

    /**
     * @dev 构造函数，初始化计数器为0，设置合约所有者
     */
    constructor() {
        count = 0;
        owner = msg.sender;
        lastModifiedBlock = block.number;
        emit CountChanged(count, msg.sender, block.number);
    }

    /**
     * @dev 增加计数器值（任何人都可以调用）
     */
    function increment() public {
        count += 1;
        lastModifiedBlock = block.number;
        emit CountChanged(count, msg.sender, block.number);
    }

    /**
     * @dev 减少计数器值（任何人都可以调用）
     */
    function decrement() public {
        // 使用require确保计数器不会变成负数
        require(count > 0, "Counter cannot be negative");
        count -= 1;
        lastModifiedBlock = block.number;
        emit CountChanged(count, msg.sender, block.number);
    }

    /**
     * @dev 重置计数器为0（只有合约所有者可以调用）
     */
    function reset() public onlyOwner {
        count = 0;
        lastModifiedBlock = block.number;
        emit CountChanged(count, msg.sender, block.number);
    }

    /**
     * @dev 获取当前计数器的值
     * @return 当前计数值
     */
    function getCount() public view returns (uint256) {
        return count;
    }

    /**
     * @dev 获取合约信息
     * @return 当前计数值、合约所有者地址和最后修改的区块号
     */
    function getCounterInfo() public view returns (uint256, address, uint256) {
        return (count, owner, lastModifiedBlock);
    }
}