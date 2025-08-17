// ### ✅ 作业3：编写一个讨饭合约
// 任务目标
// 1. 使用 Solidity 编写一个合约，允许用户向合约地址发送以太币。
// 2. 记录每个捐赠者的地址和捐赠金额。
// 3. 允许合约所有者提取所有捐赠的资金。
// 任务步骤
// 1. 编写合约
//   - 创建一个名为 BeggingContract 的合约。
//   - 合约应包含以下功能：
//   - 一个 mapping 来记录每个捐赠者的捐赠金额。
//   - 一个 donate 函数，允许用户向合约发送以太币，并记录捐赠信息。
//   - 一个 withdraw 函数，允许合约所有者提取所有资金。
//   - 一个 getDonation 函数，允许查询某个地址的捐赠金额。
//   - 使用 payable 修饰符和 address.transfer 实现支付和提款。
// 2. 部署合约
//   - 在 Remix IDE 中编译合约。
//   - 部署合约到 Goerli 或 Sepolia 测试网。
// 3. 测试合约
//   - 使用 MetaMask 向合约发送以太币，测试 donate 功能。
//   - 调用 withdraw 函数，测试合约所有者是否可以提取资金。
//   - 调用 getDonation 函数，查询某个地址的捐赠金额。
// 任务要求
// 1. 合约代码：
//   - 使用 mapping 记录捐赠者的地址和金额。
//   - 使用 payable 修饰符实现 donate 和 withdraw 函数。
//   - 使用 onlyOwner 修饰符限制 withdraw 函数只能由合约所有者调用。
// 2. 测试网部署：
//   - 合约必须部署到 Goerli 或 Sepolia 测试网。
// 3. 功能测试：
//   - 确保 donate、withdraw 和 getDonation 函数正常工作。
// 提交内容
// 1. 合约代码：提交 Solidity 合约文件（如 BeggingContract.sol）。
// 2. 合约地址：提交部署到测试网的合约地址。
// 3. 测试截图：提交在 Remix 或 Etherscan 上测试合约的截图。
// 额外挑战（可选）
// 1. 捐赠事件：添加 Donation 事件，记录每次捐赠的地址和金额。
// 2. 捐赠排行榜：实现一个功能，显示捐赠金额最多的前 3 个地址。
// 3. 时间限制：添加一个时间限制，只有在特定时间段内才能捐赠。

// 2. 合约地址: 0xFc2C2dC0d8b3f389062639bcd75E460722ea87e6
// 3.


// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BeggingContract{
    // 记录每个捐赠者的捐赠金额
    mapping(address => uint) private donorRecording;
    // 存储所有捐赠者地址
    address[] private donorAddresses;
    // 合约所有者
    address public immutable owner;

    // 捐赠时间限制
    uint256 public donationStartTime;
    uint256 public donationEndTime;

    // 排行榜结构
    struct DonorRank {
        address donor;
        uint amount;
    }
    DonorRank[3] public topDonors;

    // 事件
    event Donated(address indexed donor, uint amount);
    event Withdrawn(address indexed owner, uint amount);
    event DonationPeriodSet(uint start, uint end);

    // 修饰器：仅所有者可调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this");
        _;
    }

    // 修饰器：检查捐赠时间
    modifier onlyDuringDonationPeriod() {
        require(block.timestamp >= donationStartTime && block.timestamp <= donationEndTime,
            "Donations not allowed at this time");
        _;
    }

    constructor() {
        owner = msg.sender; // 记录部署者为所有者
        donationStartTime = block.timestamp;
        donationEndTime = block.timestamp + 30 days; // 默认30天捐赠期
    }

    // 接收以太币的fallback函数
    receive() external payable onlyDuringDonationPeriod {
        recordDonation(msg.sender, msg.value);
    }

    // 捐赠函数（支持直接发送ETH）
    function donate() external payable onlyDuringDonationPeriod {
        require(msg.value > 0, "Donation amount must be greater than 0");
        recordDonation(msg.sender, msg.value);
    }

    // 内部函数：记录捐赠
    function recordDonation(address donor, uint amount) private {
        // 如果是第一次捐赠，记录地址
        if (donorRecording[donor] == 0) {
            donorAddresses.push(donor);
        }

        // 更新捐赠记录
        donorRecording[donor] += amount;

        // 更新排行榜
        updateTopDonors(donor, donorRecording[donor]);

        emit Donated(donor, amount);
    }

    function updateTopDonors(address donor, uint amount) private {
        // 检查捐赠者是否已在榜上
        int existingIndex = -1;
        for (uint i = 0; i < 3; i++) {
            if (topDonors[i].donor == donor) {
                existingIndex = int(i);
                break;
            }
        }

        // 移除旧记录（如果存在）
        if (existingIndex >= 0) {
            for (uint i = uint(existingIndex); i < 2; i++) {
                topDonors[i] = topDonors[i+1];
            }
            topDonors[2] = DonorRank(address(0), 0); // 清空末尾
        }

        // 确定新记录插入位置
        uint8 newIndex = 3;
        for (uint8 i = 0; i < 3; i++) {
            if (amount > topDonors[i].amount) {
                newIndex = i;
                break;
            }
        }

        // 插入新记录并后移
        if (newIndex < 3) {
            for (uint j = 2; j > newIndex; j--) {
                topDonors[j] = topDonors[j-1];
            }
            topDonors[newIndex] = DonorRank(donor, amount);
        }
    }

    // 提取合约资金（仅所有者）
    function withdraw() external onlyOwner {
        uint contractBalance = address(this).balance;
        require(contractBalance > 0, "No funds to withdraw");

        payable(owner).transfer(contractBalance);
        emit Withdrawn(owner, contractBalance);
    }

    // 查询捐赠金额
    function getDonation(address donor) public view returns (uint) {
        return donorRecording[donor];
    }

    // 获取捐赠者总数
    function getDonorCount() public view returns (uint) {
        return donorAddresses.length;
    }

    // 通过索引获取捐赠者地址
    function getDonorByIndex(uint index) public view returns (address) {
        require(index < donorAddresses.length, "Index out of bounds");
        return donorAddresses[index];
    }

    // 设置捐赠时间段（仅所有者）
    function setDonationPeriod(uint start, uint end) external onlyOwner {
        require(start < end, "Invalid time range");
        donationStartTime = start;
        donationEndTime = end;
        emit DonationPeriodSet(start, end);
    }

    // 获取合约余额
    function getContractBalance() public view returns (uint) {
        return address(this).balance;
    }
}

