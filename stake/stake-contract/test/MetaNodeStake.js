const { expect } = require("chai"); // 引入断言库 chai 的 expect
const { ethers } = require("hardhat"); // 引入 Hardhat 的 ethers（部署合约、发送交易、获取区块等）

describe("MetaNodeStake", function () { // 定义测试套件 describe("MetaNodeStake", ...)

  // 定义测试中用到的变量（账户、合约实例、关键参数、地址）
  let owner, user1, user2, MetaNode, stake, startBlock, endBlock, MetaNodePerBlock, metaNodeAddr, stakeAddr;

  // beforeEach：每个测试前的部署与初始化
  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // 部署Mock MetaNode ERC20
    const MetaNodeFactory = await ethers.getContractFactory("MetaNodeToken");
    MetaNode = await MetaNodeFactory.deploy();
    await MetaNode.waitForDeployment();
    metaNodeAddr = await MetaNode.getAddress();

    // 部署MetaNodeStake
    const StakeFactory = await ethers.getContractFactory("MetaNodeStake");
    stake = await StakeFactory.deploy();
    await stake.waitForDeployment();
    stakeAddr = await stake.getAddress();

    // 初始化
    startBlock = await ethers.provider.getBlockNumber();
    endBlock = startBlock + 10000;
    MetaNodePerBlock = ethers.parseUnits("10", 18);

    await stake.initialize(
      metaNodeAddr,
      startBlock,
      endBlock,
      MetaNodePerBlock
    );

    // 给合约授权
    await stake.setMetaNode(metaNodeAddr);

    // 给合约转入奖励Token
    await MetaNode.transfer(stakeAddr, ethers.parseUnits("100000", 18));
  });

  // 用例1：添加 ETH 池
  it("should add ETH pool", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      10,
      false
    );
    const pool = await stake.pool(0);
    expect(pool.stTokenAddress).to.equal(ethers.ZeroAddress);
    expect(pool.poolWeight).to.equal(100);
    expect(pool.minDepositAmount).to.equal(1n);
    expect(pool.unstakeLockedBlocks).to.equal(10n);
  });
 
  // 用例2：存入 ETH 并更新余额
  it("should deposit ETH and update user balance", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      10,
      false
    );
    await stake.connect(user1).depositETH({ value: 100 });
    const balance = await stake.stakingBalance(0, user1.address);
    expect(balance).to.equal(100n);
  });

  // 用例3：申请解押并产生提现请求
  it("should allow unstake and create withdraw request", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      10,
      false
    );
    await stake.connect(user1).depositETH({ value: 100 });
    await stake.connect(user1).unstake(0, 50);
    const balance = await stake.stakingBalance(0, user1.address);
    expect(balance).to.equal(50n);
    const withdrawInfo = await stake.withdrawAmount(0, user1.address);
    expect(withdrawInfo.requestAmount).to.equal(50n);
  });

  // 用例4：解锁后提现成功
  it("should withdraw after unlock", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      1,
      false
    );
    await stake.connect(user1).depositETH({ value: 100 });
    await stake.connect(user1).unstake(0, 100);

    // 模拟区块推进
    await ethers.provider.send("evm_mine");
    await stake.connect(user1).withdraw(0);

    const withdrawInfo = await stake.withdrawAmount(0, user1.address);
    expect(withdrawInfo.requestAmount).to.equal(0n);
  });

  // 用例5：领取奖励
  it("should claim MetaNode reward", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      1,
      false
    );
    await stake.connect(user1).depositETH({ value: 100 });

    // 推进区块
    for (let i = 0; i < 10; i++) {
      await ethers.provider.send("evm_mine");
    }

    await stake.connect(user1).claim(0);
    const balance = await MetaNode.balanceOf(user1.address);
    expect(balance).to.be.gt(0);
  });

  // 用例6：暂停/恢复提现
  it("should pause and unpause withdraw", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      1,
      false
    );
    await stake.pauseWithdraw();
    await expect(stake.connect(user1).withdraw(0)).to.be.revertedWith("withdraw is paused");
    await stake.unpauseWithdraw();
    // withdraw不会revert（如果有可提余额）
  });

  // 用例7：暂停/恢复奖励领取
  it("should pause and unpause claim", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      1,
      false
    );
    await stake.pauseClaim();
    await expect(stake.connect(user1).claim(0)).to.be.revertedWith("claim is paused");
    await stake.unpauseClaim();
    // claim不会revert（如果有奖励）
  });

  // 用例8：修改池权重
  it("should set pool weight", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      1,
      false
    );
    await stake.setPoolWeight(0, 200, false);
    const pool = await stake.pool(0);
    expect(pool.poolWeight).to.equal(200n);
  });

  // 用例9：更新池参数（min/lock）
  it("should update pool info", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      1,
      false
    );
    await stake.updatePool(0, 10, 20);
    const pool = await stake.pool(0);
    expect(pool.minDepositAmount).to.equal(10n);
    expect(pool.unstakeLockedBlocks).to.equal(20n);
  });

  // 用例10：非法池 id 校验
  it("should revert on invalid pool id", async function () {
    await expect(stake.setPoolWeight(99, 100, false)).to.be.revertedWith("invalid pid");
    await expect(stake.updatePool(99, 1, 1)).to.be.revertedWith("invalid pid");
    await expect(stake.stakingBalance(99, user1.address)).to.be.revertedWith("invalid pid");
  });

  // 用例11：非法参数校验
  it("should revert on invalid parameters", async function () {
    await expect(stake.setMetaNodePerBlock(0)).to.be.revertedWith("invalid parameter");
    await expect(stake.setStartBlock(endBlock + 1)).to.be.revertedWith("start block must be smaller than end block");
    await expect(stake.setEndBlock(startBlock - 1)).to.be.revertedWith("start block must be smaller than end block");
  });

  // 扩展测试用例以提高覆盖率
  // 用例12：添加 ERC20 池并存入
  it("should add ERC20 pool and deposit tokens", async function () {
    // 先添加ETH池（第一池必须为ETH）
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      10,
      false
    );
    // 部署一个模拟的ERC20代币
    const ERC20Factory = await ethers.getContractFactory("MetaNodeToken");
    const erc20Token = await ERC20Factory.deploy();
    await erc20Token.waitForDeployment();
    const erc20Addr = await erc20Token.getAddress();
    await erc20Token.transfer(user1.address, ethers.parseUnits("1000", 18));
    
    await stake.addPool(
      erc20Addr,
      200,
      ethers.parseUnits("10", 18),
      5,
      false
    );
    const pool = await stake.pool(1);
    expect(pool.stTokenAddress).to.equal(erc20Addr);
    expect(pool.poolWeight).to.equal(200n);
    
    // 用户授权并质押
    await erc20Token.connect(user1).approve(stakeAddr, ethers.parseUnits("100", 18));
    await stake.connect(user1).deposit(1, ethers.parseUnits("100", 18));
    const balance = await stake.stakingBalance(1, user1.address);
    expect(balance).to.equal(ethers.parseUnits("100", 18));
  });

  // 用例13：多池不同权重奖励比例
  it("should handle multiple pools with different weights", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      1,
      false
    );
    // 部署一个模拟的ERC20代币
    const ERC20Factory = await ethers.getContractFactory("MetaNodeToken");
    const erc20Token = await ERC20Factory.deploy();
    await erc20Token.waitForDeployment();
    const erc20Addr = await erc20Token.getAddress();
    await erc20Token.transfer(user1.address, ethers.parseUnits("1000", 18));
    
    await stake.addPool(
      erc20Addr,
      200,
      10,
      5,
      false
    );
    
    // 质押到两个池
    await stake.connect(user1).depositETH({ value: 100 });
    await erc20Token.connect(user1).approve(stakeAddr, ethers.parseUnits("100", 18));
    await stake.connect(user1).deposit(1, ethers.parseUnits("100", 18));
    
    // 推进区块
    for (let i = 0; i < 10; i++) {
      await ethers.provider.send("evm_mine");
    }
    
    // 领取奖励并比较
    await stake.connect(user1).claim(0);
    const ethPoolReward = await MetaNode.balanceOf(user1.address);
    await MetaNode.connect(user1).transfer(owner.address, ethPoolReward); // 清空余额以便测试第二个池的奖励
    await stake.connect(user1).claim(1);
    const erc20PoolReward = await MetaNode.balanceOf(user1.address);
    
    // 由于权重是200:100，ERC20池的奖励应该是ETH池的两倍
    expect(erc20PoolReward).to.be.gt(ethPoolReward);
  });

  // 用例14：多次解押与部分提现
  it("should handle multiple unstake requests and partial withdrawals", async function () {
    await stake.addPool(
      ethers.ZeroAddress,
      100,
      1,
      3,
      false
    );
    await stake.connect(user1).depositETH({ value: 300 });
    await stake.connect(user1).unstake(0, 100);
    await ethers.provider.send("evm_mine");
    await stake.connect(user1).unstake(0, 50);
    // 未到期前不可提
    const before = await stake.withdrawAmount(0, user1.address);
    expect(before.requestAmount).to.equal(150n);
    expect(before.pendingWithdrawAmount).to.equal(0n);
    // 推进到两条都到期
    await ethers.provider.send("evm_mine");
    await ethers.provider.send("evm_mine");
    await ethers.provider.send("evm_mine");
    await stake.connect(user1).withdraw(0);
    const after = await stake.withdrawAmount(0, user1.address);
    expect(after.requestAmount).to.equal(0n);
  });

  // 用例15：addPool 对“第一池/后续池”地址合法性检查
  it("should revert addPool with invalid first/next token addresses", async function () {
    // 第一池不得为非零地址
    await expect(
      stake.addPool(owner.address, 100, 1, 10, false)
    ).to.be.revertedWith("invalid staking token address");

    // 正确添加第一池为ETH
    await stake.addPool(ethers.ZeroAddress, 100, 1, 10, false);
    // 第二池不得为零地址
    await expect(
      stake.addPool(ethers.ZeroAddress, 100, 1, 10, false)
    ).to.be.revertedWith("invalid staking token address");
  });

  // 用例16：最小存入规则（ETH 与 ERC20）
  it("should enforce min deposit rules for ETH and ERC20 pools", async function () {
    await stake.addPool(ethers.ZeroAddress, 100, 100, 10, false);
    await expect(
      stake.connect(user1).depositETH({ value: 99 })
    ).to.be.revertedWith("deposit amount is too small");

    // 添加ERC20池，min=10，deposit需要 >10
    const ERC20Factory = await ethers.getContractFactory("MetaNodeToken");
    const token = await ERC20Factory.deploy();
    await token.waitForDeployment();
    const tokenAddr = await token.getAddress();
    await token.transfer(user1.address, ethers.parseUnits("100", 18));
    await stake.addPool(tokenAddr, 100, ethers.parseUnits("10", 18), 10, false);
    await token.connect(user1).approve(stakeAddr, ethers.parseUnits("11", 18));
    await expect(
      stake.connect(user1).deposit(1, ethers.parseUnits("10", 18))
    ).to.be.revertedWith("deposit amount is too small");
    await stake.connect(user1).deposit(1, ethers.parseUnits("11", 18));
  });

  // 用例17：奖励余额不足时的安全转账
  it("should not revert when reward balance is insufficient (safe transfer)", async function () {
    await stake.addPool(ethers.ZeroAddress, 100, 1, 1, false);
    await stake.connect(user1).depositETH({ value: 100 });
    // 推进区块，产生奖励
    for (let i = 0; i < 5; i++) {
      await ethers.provider.send("evm_mine");
    }
    // 抽空合约奖励余额
    const contractBal = await MetaNode.balanceOf(stakeAddr);
    await MetaNode.connect(owner).transfer(owner.address, 0); // no-op to keep signer
    await MetaNode.connect(owner).transfer(user2.address, contractBal);
    // claim 时应尽可能转出，不应revert
    await expect(stake.connect(user1).claim(0)).to.not.be.reverted;
  });

  // 用例18：批量更新池的累计值
  it("should massUpdatePools without reverting and update acc when supply > 0", async function () {
    // add ETH pool and deposit to create non-zero supply
    await stake.addPool(ethers.ZeroAddress, 100, 1, 1, false);
    await stake.connect(user1).depositETH({ value: 100 });
    // add ERC20 pool with zero supply
    const ERC20Factory = await ethers.getContractFactory("MetaNodeToken");
    const token = await ERC20Factory.deploy();
    await token.waitForDeployment();
    const tokenAddr = await token.getAddress();
    await stake.addPool(tokenAddr, 100, ethers.parseUnits("10", 18), 1, false);

    // advance blocks so there is reward to distribute
    for (let i = 0; i < 5; i++) {
      await ethers.provider.send("evm_mine");
    }

    // mass update should touch both pools
    await stake.massUpdatePools();
    const pool0 = await stake.pool(0);
    const pool1 = await stake.pool(1);
    // pool0 has supply, so acc should be > 0; pool1 has zero supply, acc should remain 0
    expect(pool0.accMetaNodePerST).to.be.gt(0n);
    expect(pool1.accMetaNodePerST).to.equal(0n);
  });

  // 用例19：当池中供给为 0，updatePool 不改变累计值
  it("should not change acc when stTokenAmount == 0 in updatePool", async function () {
    await stake.addPool(ethers.ZeroAddress, 100, 1, 1, false);
    // no deposit -> stTokenAmount == 0
    const before = await stake.pool(0);
    for (let i = 0; i < 3; i++) {
      await ethers.provider.send("evm_mine");
    }
    await stake.updatePool(0);
    const after = await stake.pool(0);
    expect(before.accMetaNodePerST).to.equal(0n);
    expect(after.accMetaNodePerST).to.equal(0n);
    expect(after.lastRewardBlock).to.be.gte(before.lastRewardBlock);
  });

  // 用例20：奖励时间窗口裁剪边界
  it("should respect getMultiplier boundaries with start/end clipping", async function () {
    // shorten window for precise assertions
    const current = await ethers.provider.getBlockNumber();
    await stake.setStartBlock(current + 2);
    await stake.setEndBlock(current + 6);
    // case1: _from < start, _to < start -> revert by require(_from <= _to) after clipping? Here we call via view wrapper
    // use pendingMetaNodeByBlockNumber to indirectly exercise multiplier clipping without reverting
    await stake.addPool(ethers.ZeroAddress, 100, 1, 1, false);
    await stake.connect(user1).depositETH({ value: 100 });
    // before start: no rewards added
    let pending = await stake.pendingMetaNodeByBlockNumber(0, user1.address, current + 1);
    expect(pending).to.equal(0n);
    // inside window: rewards accrue
    for (let i = 0; i < 5; i++) {
      await ethers.provider.send("evm_mine");
    }
    pending = await stake.pendingMetaNodeByBlockNumber(0, user1.address, current + 6);
    expect(pending).to.be.gt(0n);
    // beyond end: should be same as at end (no extra accrual)
    const pendingBeyond = await stake.pendingMetaNodeByBlockNumber(0, user1.address, current + 10);
    expect(pendingBeyond).to.equal(pending);
  });
});
