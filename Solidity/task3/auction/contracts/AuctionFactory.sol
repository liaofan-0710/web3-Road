// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./SingleAuction.sol"; // 导入单个拍卖合约

contract AuctionFactory is Initializable, UUPSUpgradeable {
    // 管理员地址
    address public admin;
    
    // 所有拍卖合约地址列表
    address[] public allAuctions;
    
    // 映射：拍卖ID => 拍卖合约地址
    mapping(uint256 => address) public auctionAddresses;
    
    // 映射：用户地址 => 创建的拍卖合约地址列表
    mapping(address => address[]) public userAuctions;
    
    // 事件：新拍卖创建
    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed auctionAddress,
        address indexed seller,
        address nftContract,
        uint256 tokenId
    );
    
    // 初始化函数
    function initialize() public initializer {
        admin = msg.sender;
    }
    
    /**
     * @dev 创建新的拍卖合约
     * @param _duration 拍卖持续时间
     * @param _startPrice 起始价格（以USD为单位，使用8位小数）
     * @param _nftAddress NFT合约地址
     * @param _tokenId NFT Token ID
     * @param _tokenAddress 支付代币地址（address(0)表示ETH）
     * @return auctionAddress 新创建的拍卖合约地址
     */
    function createAuction(
        uint256 _duration,
        uint256 _startPrice,
        address _nftAddress,
        uint256 _tokenId,
        address _tokenAddress
    ) external returns (address auctionAddress) {
        require(_duration >= 10 minutes, "Duration must be at least 10 minutes");
        require(_startPrice > 0, "Start price must be greater than 0");
        
        // 部署新的拍卖合约
        SingleAuction newAuction = new SingleAuction(
            msg.sender,
            _duration,
            _startPrice,
            _nftAddress,
            _tokenId,
            _tokenAddress,
            address(this) // 传入工厂合约地址以便回调
        );
        
        // 获取新拍卖合约地址
        auctionAddress = address(newAuction);
        
        // 记录拍卖合约
        uint256 auctionId = allAuctions.length;
        auctionAddresses[auctionId] = auctionAddress;
        allAuctions.push(auctionAddress);
        userAuctions[msg.sender].push(auctionAddress);
        
        // 转移NFT到拍卖合约
        IERC721(_nftAddress).safeTransferFrom(msg.sender, auctionAddress, _tokenId);
        
        // 触发事件
        emit AuctionCreated(
            auctionId,
            auctionAddress,
            msg.sender,
            _nftAddress,
            _tokenId
        );
        
        return auctionAddress;
    }
    
    /**
     * @dev 获取所有拍卖合约数量
     */
    function getAuctionsCount() external view returns (uint256) {
        return allAuctions.length;
    }
    
    /**
     * @dev 获取用户创建的拍卖合约数量
     */
    function getUserAuctionsCount(address user) external view returns (uint256) {
        return userAuctions[user].length;
    }
    
    // 升级授权函数
    function _authorizeUpgrade(address) internal view override {
        require(msg.sender == admin, "Only admin can upgrade");
    }
}