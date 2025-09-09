// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Import the NftAuction contract explicitly
import "./NftAuctionTwo.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/contracts/interfaces/IRouterClient.sol";
import {OwnerIsCreator} from "@chainlink/contracts/src/v0.8/shared/access/OwnerIsCreator.sol";
import {Client} from "@chainlink/contracts-ccip/contracts/libraries/Client.sol";
import {LinkTokenInterface} from "@chainlink/contracts/src/v0.8/shared/interfaces/LinkTokenInterface.sol";

contract NftAuctionFactory is Initializable, UUPSUpgradeable {

    // 结构体
    struct Auction {
        // 卖家
        address seller;
        // 拍卖持续时间
        uint256 duration;
        // 起始价格
        uint256 startPrice;
        // 开始时间
        uint256 startTime;
        // 是否结束
        bool ended;
        // 最高出价者
        address highestBidder;
        // 最高出价
        uint256 highestBid;
         // NFT合约地址
        address nftContract;
        // NFT ID
        uint256 tokenId;
        // 参与竞价的资产类型 0x 地址表示eth, 其他表示erc20
        address tokenAddress;
    }

    // 所有创建的拍卖合约地址数组
    address[] public auctions;
    // 管理员地址
    address public admin;
    // 手续费接收地址
    address public feeRecipient;
    // 卖家地址
    address public seller;
    // 当前拍卖
    Auction public currentAuction; 
    // 手续费率
    uint256 public feeRate;  

    // 映射：NFT 合约地址 + Token ID => 拍卖合约地址
    // mapping (uint256 tokenId => NftAuction) public auctionMap;
    mapping(uint256 => NftAuctionTwo) public auctionMap;

    // 事件：新拍卖创建
    // event AuctionCreated(address indexed auctionAddress,uint256 tokenId);
    event AuctionCreated(address indexed auction, address indexed seller, address nftContract, uint256 tokenId);
    
    function initialize(address _seller, uint256 _duration, uint256 _startPrice, address _nftContract, uint256 _tokenId, address _feeRecipient, uint256 _feeRate) public initializer {
        __UUPSUpgradeable_init();
        admin = msg.sender;
        currentAuction.seller = _seller; // 卖家
        currentAuction.duration = _duration; // 拍卖持续时间
        currentAuction.startPrice = _startPrice; // 起始价格
        currentAuction.startTime = block.timestamp; // 开始时间
        currentAuction.ended = false; // 是否结束
        currentAuction.highestBidder = address(0); // 最高出价者
        currentAuction.highestBid = 0; // 最高出价
        currentAuction.nftContract = _nftContract; // NFT合约地址
        currentAuction.tokenId = _tokenId; // NFT ID
        currentAuction.tokenAddress = address(0); // 参与竞价的资产类型 0x 地址表示eth, 其他表示erc20
        feeRecipient = _feeRecipient; // 手续费接收地址
        feeRate = _feeRate; // 手续费率
    }

    // Create a new auction
    function createAuctionFactory(
        // address _feeRecipient,
        // address _seller,
        // uint256 _duration,
        // uint256 _startPrice,
        // address _nftContract,
        uint256 tokenId
        // uint256 _feeRate
    ) external returns (address) {
        NftAuctionTwo auction = new NftAuctionTwo();
        auction.initialize(
            seller,
            currentAuction.duration,
            currentAuction.startPrice,
            currentAuction.nftContract,
            tokenId,
            feeRecipient,
            feeRate
        );
        auctions.push(address(auction));
        auctionMap[tokenId] = auction;
        // auctionMap[tokenId] = address(auction);

        // emit AuctionCreated(address(auction), tokenId);
        emit AuctionCreated(address(auction), seller, currentAuction.nftContract, tokenId);
        return address(auction);
    }

    function getAuctions() external view returns (address[] memory) {
        return auctions;
    }

    function getAuction(uint256 tokenId) external view returns (address) {
        // require(tokenId < auctions.length, "tokenId out of bounds");
        return auctions[tokenId];
    }

    function _authorizeUpgrade(address) internal view override {
        require(msg.sender == admin, "Not admin");

    }
}