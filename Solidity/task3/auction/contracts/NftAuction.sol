// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/interfaces/IERC721Receiver.sol";
// import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract NftAuction is Initializable, UUPSUpgradeable, IERC721Receiver, ERC721URIStorage {
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

    // 状态变量
    mapping(uint256 => Auction) public auctions;
    // 下一个拍卖ID
    uint256 public nextAuctionId;
    // 管理员地址
    address public admin;
    // 手续费接收地址
    // address public feeRecipient;
    // 卖家地址
    // address public seller;
    // 当前拍卖
    // Auction public currentAuction; 
    // 手续费率
    // uint256 public feeRate;  

    // 喂价合约
    // AggregatorV3Interface internal priceETHFeed;
    mapping(address => AggregatorV3Interface) public priceFeeds;

    // 构造函数
    // 传入参数：卖家地址，NFT合约地址，NFT ID，起始价格，拍卖持续时间，手续费接收地址，手续费率
    // 传入参数：卖家地址，NFT合约地址，NFT ID，起始价格，拍卖持续时间，手续费接收地址，手续费率
    // constructor(
    //     address _seller,
    //     uint256 _duration,
    //     uint256 _startPrice,
    //     address _nftContract,
    //     uint256 _nftToken,
    //     address _feeRecipient,
    //     uint256 _feeRate
    // ) {
    //     currentAuction.seller = _seller; // 卖家
    //     currentAuction.duration = _duration; // 拍卖持续时间
    //     currentAuction.startPrice = _startPrice; // 起始价格
    //     currentAuction.startTime = block.timestamp; // 开始时间
    //     currentAuction.ended = false; // 是否结束
    //     currentAuction.highestBidder = address(0); // 最高出价者
    //     currentAuction.highestBid = 0; // 最高出价
    //     currentAuction.nftContract = _nftContract; // NFT合约地址
    //     currentAuction.tokenId = _nftToken; // NFT ID
    //     currentAuction.tokenAddress = address(0); // 参与竞价的资产类型 0x 地址表示eth, 其他表示erc20
    //     feeRecipient = _feeRecipient; // 手续费接收地址
    //     feeRate = _feeRate; // 手续费率
    // }

    function initialize() public initializer {
        __UUPSUpgradeable_init();
        admin = msg.sender;
    }
    // function initialize(address _seller, uint256 _duration, uint256 _startPrice, address _nftContract, uint256 _tokenId, address _feeRecipient, uint256 _feeRate) public initializer {
    //     __UUPSUpgradeable_init();
    //     admin = msg.sender;
    //     currentAuction.seller = _seller; // 卖家
    //     currentAuction.duration = _duration; // 拍卖持续时间
    //     currentAuction.startPrice = _startPrice; // 起始价格
    //     currentAuction.startTime = block.timestamp; // 开始时间
    //     currentAuction.ended = false; // 是否结束
    //     currentAuction.highestBidder = address(0); // 最高出价者
    //     currentAuction.highestBid = 0; // 最高出价
    //     currentAuction.nftContract = _nftContract; // NFT合约地址
    //     currentAuction.tokenId = _tokenId; // NFT ID
    //     currentAuction.tokenAddress = address(0); // 参与竞价的资产类型 0x 地址表示eth, 其他表示erc20
    //     feeRecipient = _feeRecipient; // 手续费接收地址
    //     feeRate = _feeRate; // 手续费率
    // }

    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    // 继承 ERC721URIStorage 来存储 tokenURI
    function mintNFT(address recipient, string memory tokenURI)
    public
    returns (uint256)
    {
        require(msg.sender == admin, "Only admin can mint");
        _tokenIds.increment();
        uint256 newItemId = _tokenIds.current();
        _mint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);
        return newItemId;
    }

    // 设置喂价合约
    // 0x 地址表示eth, 其他表示erc20
    function setPriceETHFeed(address tokenAddress, address _priceETHFeed) public {
        // require(msg.sender == admin, "Only admin can set price feed");
        // priceETHFeed = AggregatorV3Interface(_priceETHFeed);
        priceFeeds[tokenAddress] = AggregatorV3Interface(_priceETHFeed);
    } 

    // ETH -> USD => 4322 0063 0000 => 4322.00630000
    // USDC -> USD => 9998 5058 => 0.99985058
    function getChainlinkDataFeedLatestAnswer(address tokenAddress) public view returns (int) {
        // prettier-ignore
        AggregatorV3Interface priceFeed = priceFeeds[tokenAddress];
        (
            /* uint80 roundId */,
            int256 answer,
            /*uint256 startedAt*/,
            /*uint256 updatedAt*/,
            /*uint80 answeredInRound*/
        ) = priceFeed.latestRoundData();
        return answer;
    }

    // 创建拍卖
    function createAuction(uint256 _duration, uint256 _startPrice, address _nftAddress, uint256 _tokenId) external {
        // 只有管理员可以创建拍卖
        // require(msg.sender == admin, "Only admin can create auction");
        // 检查参数有效性
        require(_duration >= 10, "Duration must be at least 10 minute");
        require(_startPrice > 0, "Start price must be greater than 0");

        // 转移NFT到合约
        // IERC721(_nftAddress).approve(address(this), _tokenId);
        IERC721(_nftAddress).safeTransferFrom(msg.sender, address(this), _tokenId);

        auctions[nextAuctionId] = Auction({
            seller: msg.sender,
            duration: _duration,
            startPrice: _startPrice,
            ended: false,
            highestBidder: address(0),
            highestBid: 0,
            startTime: block.timestamp,
            nftContract: _nftAddress,
            tokenId: _tokenId,
            tokenAddress: address(0)
        });

        nextAuctionId++;
    }

    // 买家参与拍卖
    // TODO ERC20 也能参与
    function placeBid(uint256 _auctionId, uint256 amount, address _tokenAddress) external payable {
        // ETH 是 ？ 美金
        // 1个 USDC 是 ？ 美金
        Auction storage auction = auctions[_auctionId];
        // 判断当前拍卖是否结束
        require(!auction.ended && auction.startTime + auction.duration > block.timestamp, "Auction has ended");

        // 1. 获取统一的价值尺度美金
        uint payValue;
        int price = getChainlinkDataFeedLatestAnswer(_tokenAddress);
        uint8 decimals = priceFeeds[_tokenAddress].decimals();
        // 判断是否是ERC20代币
        if (_tokenAddress != address(0)) { // 如果是ERC20
            // 处理ERC20代币的出价逻辑
            // 获取ERC20 代币的价格为多少美金？
            // payValue = amount * uint(getChainlinkDataFeedLatestAnswer(_tokenAddress));
            payValue = amount * uint(price) / (10 ** decimals);
        } else { // 如果是其他，例如ETH？
            amount = msg.value;
            payValue = amount * uint(getChainlinkDataFeedLatestAnswer(address(0)));
        }
        // 获得起拍价的美金价值
        uint startPriceValue = auction.startPrice * uint(getChainlinkDataFeedLatestAnswer(auction.tokenAddress));

        // 获得最高价的美金价值
        uint highestBidValue = auction.highestBid * uint(getChainlinkDataFeedLatestAnswer(auction.tokenAddress));

        // 判断出价是否高于当前最高出价和起始价 || 判断是否符合拍卖价格范围
        require(payValue >= startPriceValue && payValue > highestBidValue, "Bid must be higher than the current highest bid");

        // 转移ERC20代币 到合约 -- 要先扣款，后退款。防止出现？闪电贷攻击
        if (_tokenAddress != address(0)){
            // IERC20(_tokenAddress).transferFrom(msg.sender, address(this), amount);
            require(IERC20(_tokenAddress).transferFrom(msg.sender, address(this), amount), "ERC20 transfer failed");
        }

        // 退还前最高价 -- 后退款
        if (auction.highestBid > 0) {
            // 判断是否是ERC20代币
            if (auction.tokenAddress != address(0)) { // 如果是ERC20
                IERC20(auction.tokenAddress).transfer(
                        auction.highestBidder,
                        auction.highestBid
                );
            } else { // 如果是其他，例如ETH？
                payable(auction.highestBidder).transfer(auction.highestBid);
            }
        }

        auction.tokenAddress = _tokenAddress; // 变更当前获得者token
        auction.highestBid = amount; // 变更当前最高出价
        auction.highestBidder = msg.sender; // 变更当前最高出价者身份
    }

    // 结束拍卖
    function endAuction(uint256 _auctionId) external {
        Auction storage auction = auctions[_auctionId];
        // 判断当前拍卖是否结束
        require(!auction.ended && auction.startTime + auction.duration < block.timestamp, "Auction is still ongoing");
        // 转移NFT到最高出价者
        // IERC721(auction.nftContract).safeTransferFrom(admin, auction.highestBidder, auction.tokenId);
        // 修正 NFT 转移 -- NFT 应该从合约转移给最高出价者，而不是从管理员转移？
        IERC721(auction.nftContract).safeTransferFrom(address(this), auction.highestBidder, auction.tokenId);
        // 转移剩余的资金到卖家
        payable(auction.seller).transfer(address(this).balance);
        auction.ended = true;

    }

    function _authorizeUpgrade(address) internal view override {
        // 只有管理员可以升级合约
        require(msg.sender == admin, "Only admin can upgrade");
    }

    // 实现 IERC721Receiver 接口的 onERC721Received 函数
    function onERC721Received(address, address, uint256, bytes calldata) external pure override returns (bytes4) {
        return IERC721Receiver.onERC721Received.selector;
    }

}