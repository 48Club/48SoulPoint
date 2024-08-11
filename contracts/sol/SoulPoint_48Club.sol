// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0 <=0.8.20;

contract Token {
    function balanceOf(address account) external view returns (uint256) {}

    function getPooledBNBByShares(
        uint256 shares
    ) external view returns (uint256) {}
}

contract DAO {
    function getStake(address _user) external view returns (uint256) {}
}

contract CoinbaseHelper {
    function coinbaseToCreditAddress(
        address coinbase
    ) external view returns (address) {}
}

contract PuissantIndicator {
    function getPuissants() external view returns (address[] memory) {}
}

contract Calculator {
    Token internal constant koge =
        Token(0xe6DF05CE8C8301223373CF5B969AFCb1498c5528);
    Token internal constant nft =
        Token(0x57b81C140BdfD35dbfbB395360a66D54a650666D);
    DAO internal constant dao = DAO(0xa31F6B577704B4622d2ba63F6aa1b7e92fe8C8a9);
    CoinbaseHelper internal constant coinbaseHelper =
        CoinbaseHelper(0x32Bb57eAA91566488A76371043113bc38b144BDE);
    PuissantIndicator internal constant puissantIndicator =
        PuissantIndicator(0x5cC05FDe1D231A840061c1a2D7e913CeDc8EaBaF);

    function getPointDetail(
        address user
    ) external view returns (address, uint256, uint256, uint256, uint256) {
        uint256 kogePoint = koge.balanceOf(user) / 1e18;
        if (kogePoint > 48) {
            kogePoint = 48;
        }

        uint256 stakePoint = dao.getStake(user) / 1e18;

        uint256 nftPoint = nft.balanceOf(user);
        if (nftPoint > 0) {
            nftPoint = 480;
        }

        address[] memory puissants = puissantIndicator.getPuissants();
        uint256 bscStakePoint = 0;
        for (uint256 i = 0; i < puissants.length; i++) {
            address _credit = coinbaseHelper.coinbaseToCreditAddress(
                puissants[i]
            );
            if (_credit == address(0)) {
                continue;
            }
            uint256 balance = Token(_credit).balanceOf(user);
            bscStakePoint +=
                (Token(_credit).getPooledBNBByShares(balance) * 48) /
                1e18;
        }

        return (user, kogePoint, stakePoint, nftPoint, bscStakePoint);
    }

    function getPoint(address user) external view returns (uint256) {
        (, uint256 a, uint256 b, uint256 c, uint256 d) = this.getPointDetail(
            user
        );
        return a + b + c + d;
    }
}

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// abigen --abi contracts/sol/abi.json --pkg SoulPoint_48Club --out contracts/SoulPoint_48Club/SoulPoint_48Club.go

contract SoulPoint_48Club is ERC20, Ownable(msg.sender) {
    constructor() ERC20("48SoulPoint", "48SP") {}

    uint256 internal maxHolding = 1;
    mapping(address => bool) public isMember;
    address[] internal _members;
    Calculator internal calculator;

    function getPoint(address account) external view returns (uint256) {
        return calculator.getPoint(account);
    }

    function setCaculator(address _calculator) external onlyOwner {
        calculator = Calculator(_calculator);
    }

    function decimals() public view virtual override returns (uint8) {
        return 0;
    }

    function mint() external {
        address _to = msg.sender;
        require(this.getPoint(_to) > 0, "the minimum soul point allowed is 1");
        _mint(_to, maxHolding);
    }

    function getAllMembers() external view returns (address[] memory) {
        address[] memory members = new address[](_members.length);
        uint256 index = 0;
        for (uint256 i = 0; i < _members.length; i++) {
            if (isMember[_members[i]]) {
                members[index] = _members[i];
                index++;
            }
        }

        return members;
    }

    function _update(
        address from,
        address to,
        uint256 value
    ) internal virtual override {
        require(
            address(0) == from || address(0) == to,
            "ERC20: transfer is not allowed"
        );

        super._update(from, to, value);
        require(balanceOf(to) <= maxHolding, "Insufficient balance"); // 1 soul per address

        _updateIsMember(to);
        _updateIsMember(from);
    }

    function _updateIsMember(address account) internal {
        if (address(0) != account) {
            if (balanceOf(account) == maxHolding) {
                isMember[account] = true;
                _members.push(account);
            } else {
                isMember[account] = false;
            }
        }
    }
}
