// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract SingleAuction is ReentrancyGuard {
    // 拍卖结构体
    struct Auction {
        address seller;
        uint256 duration;
        uint256 startPrice; // 以USD为单位，使用8位小数
        uint256 startTime;
        bool ended;
        address highestBidder;
        uint256 highestBid; // 实际支付代币数量
        address nftContract;
        uint256 tokenId;
        address tokenAddress; // 支付代币地址
    }
    
    // 拍卖信息
    Auction public auction;
    
    // 工厂合约地址
    address public factory;
    
    // 喂价合约映射
    mapping(address => address) public priceFeeds;
    
    // 事件
    event BidPlaced(address indexed bidder, uint256 amount, address tokenAddress);
    event AuctionEnded(address indexed winner, uint256 amount);
    
    // 修饰器：只有工厂合约可以调用
    modifier onlyFactory() {
        require(msg.sender == factory, "Only factory can call");
        _;
    }
    
    // 构造函数
    constructor(
        address _seller,
        uint256 _duration,
        uint256 _startPrice,
        address _nftAddress,
        uint256 _tokenId,
        address _tokenAddress,
        address _factory
    ) {
        factory = _factory;
        auction = Auction({
            seller: _seller,
            duration: _duration,
            startPrice: _startPrice,
            startTime: block.timestamp,
            ended: false,
            highestBidder: address(0),
            highestBid: 0,
            nftContract: _nftAddress,
            tokenId: _tokenId,
            tokenAddress: _tokenAddress
        });
    }
    
    // 设置价格喂价合约（只能由工厂调用）
    function setPriceFeed(address tokenAddress, address priceFeed) external onlyFactory {
        priceFeeds[tokenAddress] = priceFeed;
    }
    
    // 获取Chainlink价格数据
    function getChainlinkDataFeedLatestAnswer(address tokenAddress) public view returns (int256) {
        address feedAddress = priceFeeds[tokenAddress];
        require(feedAddress != address(0), "Price feed not set for this token");
        
        AggregatorV3Interface priceFeed = AggregatorV3Interface(feedAddress);
        (
            /* uint80 roundId */,
            int256 answer,
            /*uint256 startedAt*/,
            /*uint256 updatedAt*/,
            /*uint80 answeredInRound*/
        ) = priceFeed.latestRoundData();
        return answer;
    }
    
    // 买家参与拍卖
    function placeBid(uint256 amount) external payable nonReentrant {
        require(!auction.ended, "Auction has ended");
        require(block.timestamp < auction.startTime + auction.duration, "Auction has ended");
        
        // 获取统一的价值尺度（USD）
        uint256 bidValueUSD;
        if (auction.tokenAddress != address(0)) {
            // ERC20代币出价
            bidValueUSD = amount * uint256(getChainlinkDataFeedLatestAnswer(auction.tokenAddress));
        } else {
            // ETH出价
            amount = msg.value;
            bidValueUSD = amount * uint256(getChainlinkDataFeedLatestAnswer(address(0)));
        }
        
        // 计算当前最高出价的USD价值
        uint256 currentBidValueUSD;
        if (auction.highestBid > 0) {
            currentBidValueUSD = auction.highestBid * uint256(getChainlinkDataFeedLatestAnswer(auction.tokenAddress));
        }
        
        // 计算起拍价的USD价值
        uint256 startPriceValueUSD = auction.startPrice * 10**8; // 假设startPrice已经是8位小数的USD价格
        
        // 验证出价
        require(bidValueUSD >= startPriceValueUSD, "Bid must be higher than start price");
        require(bidValueUSD > currentBidValueUSD, "Bid must be higher than current highest bid");
        
        // 处理代币转移
        if (auction.tokenAddress != address(0)) {
            IERC20(auction.tokenAddress).transferFrom(msg.sender, address(this), amount);
        }
        
        // 退还前一个最高出价者的资金
        if (auction.highestBidder != address(0)) {
            if (auction.tokenAddress != address(0)) {
                IERC20(auction.tokenAddress).transfer(auction.highestBidder, auction.highestBid);
            } else {
                payable(auction.highestBidder).transfer(auction.highestBid);
            }
        }
        
        // 更新拍卖状态
        auction.highestBidder = msg.sender;
        auction.highestBid = amount;
        
        emit BidPlaced(msg.sender, amount, auction.tokenAddress);
    }
    
    // 结束拍卖
    function endAuction() external nonReentrant {
        require(!auction.ended, "Auction already ended");
        require(
            msg.sender == auction.seller || 
            block.timestamp >= auction.startTime + auction.duration,
            "Only seller can end auction before duration"
        );
        
        auction.ended = true;
        
        if (auction.highestBidder != address(0)) {
            // 有获胜者：转移NFT给获胜者，转移资金给卖家
            IERC721(auction.nftContract).safeTransferFrom(address(this), auction.highestBidder, auction.tokenId);
            
            if (auction.tokenAddress != address(0)) {
                IERC20(auction.tokenAddress).transfer(auction.seller, auction.highestBid);
            } else {
                payable(auction.seller).transfer(auction.highestBid);
            }
            
            emit AuctionEnded(auction.highestBidder, auction.highestBid);
        } else {
            // 无获胜者：退回NFT给卖家
            IERC721(auction.nftContract).safeTransferFrom(address(this), auction.seller, auction.tokenId);
            emit AuctionEnded(address(0), 0);
        }
    }
    
    // 获取拍卖信息
    function getAuctionDetails() external view returns (
        address seller,
        uint256 duration,
        uint256 startPrice,
        uint256 startTime,
        bool ended,
        address highestBidder,
        uint256 highestBid,
        address nftContract,
        uint256 tokenId,
        address tokenAddress
    ) {
        seller = auction.seller;
        duration = auction.duration;
        startPrice = auction.startPrice;
        startTime = auction.startTime;
        ended = auction.ended;
        highestBidder = auction.highestBidder;
        highestBid = auction.highestBid;
        nftContract = auction.nftContract;
        tokenId = auction.tokenId;
        tokenAddress = auction.tokenAddress;
    }
    
    // 接收ERC721令牌
    function onERC721Received(
        address,
        address,
        uint256,
        bytes calldata
    ) external pure returns (bytes4) {
        return this.onERC721Received.selector;
    }
}