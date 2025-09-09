// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 任务1.1 创建一个名为Voting的合约
// 一个mapping来存储候选人的得票数
// 一个vote函数，允许用户投票给某个候选人
// 一个getVotes函数，返回某个候选人的得票数
// 一个resetVotes函数，重置所有候选人的得票数
contract Voting {
    mapping(string => uint256) private votesReceived;

    string[] private candidates;

    mapping(string => bool) private candidateExists;

    function vote(string memory candidate) public {
        if (!candidateExists[candidate]) {
            candidates.push(candidate);
            candidateExists[candidate] = true;
        }

        votesReceived[candidate] += 1;
    }

    function getVotes(string memory candidate) public view returns (uint256) {
        return votesReceived[candidate];
    }

    function resetVotes() public {
        for (uint i = 0; i < candidates.length; i++) {
            votesReceived[candidates[i]] = 0;
        }
    }
}


// 任务1.2 反转一个字符串。输入 "abcde"，输出 "edcba"
contract RollbackString {
    function rollbackString(string memory newStr) public pure returns (string memory) {
        bytes memory input = bytes(newStr);
        bytes memory result = new bytes(input.length);

        // 字符串反转逻辑
        for(uint i = 0; i < input.length; i++) {
            result[i] = input[input.length - i - 1];
        }

        // 更新状态变量
        string memory str = string(result);
        return str;
    }
}

// 任务1.3 用 solidity 实现整数转罗马数字
contract RomanToInteger {
    function romanToInt(string memory s) public pure returns (uint256) {
        // 将字符串转换为字节数组便于处理
        bytes memory roman = bytes(s);
        uint256 length = roman.length;
        uint256 total = 0;

        // 遍历每个字符
        for (uint256 i = 0; i < length; i++) {
            // 获取当前字符的值
            uint256 current = charToValue(roman[i]);

            // 如果还有下一个字符，则获取下一个字符的值
            if (i < length - 1) {
                uint256 next = charToValue(roman[i + 1]);

                // 检查是否需要应用减法规则（小值在左）
                if (current < next) {
                    total += (next - current);
                    i++; // 跳过下一个字符，因为已处理
                    continue;
                }
            }

            // 正常情况：大值在左或单独字符
            total += current;
        }

        return total;
    }

    // 辅助函数：罗马字符转数值
    function charToValue(bytes1 char) private pure returns (uint256) {
        if (char == "I") return 1;
        if (char == "V") return 5;
        if (char == "X") return 10;
        if (char == "L") return 50;
        if (char == "C") return 100;
        if (char == "D") return 500;
        if (char == "M") return 1000;
        revert("Invalid Roman character");
    }
}

// 任务1.4 用 solidity 实现整数转罗马数字
contract IntToRoman {
    function intToRoman(uint num) public pure returns (string memory)  {
        bytes memory result;
        uint[13] memory nums = [uint(1000), 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
        string[13] memory strNums = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];
        uint i = 0;
        while (num > 0) {
            if (num >= nums[i]) {
                result = abi.encodePacked(result, strNums[i]);
                num -= nums[i];
            } else {
                i++;
            }
        }

        return string(result);
    }
}

// 任务1.5 合并两个有序数组 (Merge Sorted Array)
contract Merge {
    function merge(uint[] memory nums1, uint m, uint[] memory nums2, uint n) public pure returns (uint[] memory)  {
        uint[] memory result = new uint[](m + n);
        uint index1 = 0; uint index2 = 0; uint k = 0;
        while (index1 < m && index2 < n){
            if (nums1[index1] <= nums2[index2]) {
                result[k] = nums1[index1];
                index1++;
            } else {
                result[k] = nums2[index2];
                index2++;
            }
            k++;
        }
        while (index1 < m) {
            result[k] = nums1[index1];
            index1++;
            k++;
        }
        while (index2 < n) {
            result[k] = nums2[index2];
            index2++;
            k++;
        }
        return result;
    }
}

// 任务1.6 二分查找 (Binary Search)
contract Search {
    function search(int[] memory nums, int target) public pure returns (int)  {
        int l = 0;
        int r = int(nums.length) - 1;

        while (l <= r) {
            int mid = l + (r - l) / 2;
            int midVal = nums[uint(mid)];

            if (midVal < target) {
                l = mid + 1;
            } else if (midVal > target) {
                r = mid - 1;
            } else {
                return int(uint(mid));
            }
        }
        return -1;
    }
}

