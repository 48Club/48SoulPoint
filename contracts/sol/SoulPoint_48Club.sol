// SPDX-License-Identifier: MIT
pragma solidity >=0.8.28;

contract Token {
    function balanceOf(address account) external view returns (uint256) {}

    function getPooledBNBByShares(
        uint256 shares
    ) external view returns (uint256) {}

    function delegates(address account) external view returns (address) {}
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
    Token internal constant govBNB =
        Token(0x0000000000000000000000000000000000002005);

    DAO internal constant dao = DAO(0xa31F6B577704B4622d2ba63F6aa1b7e92fe8C8a9);

    CoinbaseHelper internal constant coinbaseHelper =
        CoinbaseHelper(0x32Bb57eAA91566488A76371043113bc38b144BDE);

    PuissantIndicator internal constant puissantIndicator =
        PuissantIndicator(0x5cC05FDe1D231A840061c1a2D7e913CeDc8EaBaF);

    address internal constant The48ClubGovAddr =
        0xaACc290a1A4c89F5D7bc29913122F5982916de48;
    uint256 internal constant decimals = 10 ** 18;

    function getPointDetail(
        address user
    )
        external
        view
        returns (address, uint256, uint256, uint256, uint256, uint256)
    {
        return (
            user,
            this.getKogePoint(user),
            this.getStakePoint(user),
            this.getNftPoint(user),
            this.getBscStakePoint(user),
            this.getGovBNBPoint(user)
        );
    }

    function getPoint(address user) external view returns (uint256) {
        return
            this.getKogePoint(user) +
            this.getStakePoint(user) +
            this.getNftPoint(user) +
            this.getBscStakePoint(user) +
            this.getGovBNBPoint(user);
    }

    function getGovBNBPoint(address user) external view returns (uint256) {
        if (govBNB.delegates(user) == The48ClubGovAddr) {
            return (govBNB.balanceOf(user) * 12) / decimals;
        }
        return 0;
    }

    function getKogePoint(address user) external view returns (uint256) {
        uint256 kogePoint = koge.balanceOf(user) / decimals;
        if (kogePoint > 48) {
            return 48;
        }
        return kogePoint;
    }

    function getStakePoint(address user) external view returns (uint256) {
        return dao.getStake(user) / decimals;
    }

    function getNftPoint(address user) external view returns (uint256) {
        uint256 nftPoint = nft.balanceOf(user);
        if (nftPoint > 0) {
            return 480;
        }
        return 0;
    }

    function getBscStakePoint(address user) external view returns (uint256) {
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
            bscStakePoint += Token(_credit).getPooledBNBByShares(balance);
        }

        return (bscStakePoint * 24) / decimals;
    }
}

contract OldContract {
    function getAllMembers() external view returns (address[] memory) {}
}

// abigen --abi contracts/sol/abi.json --pkg SoulPoint_48Club --out contracts/SoulPoint_48Club/SoulPoint_48Club.go

contract SoulPoint_48Club {
    mapping(address => bool) public isMember;
    address[] internal _members;
    Calculator internal calculator;
    address internal owner;

    constructor(address _calculator) {
        calculator = Calculator(_calculator);
        owner = msg.sender;
    }

    OldContract internal constant oldContract =
        OldContract(0x928dC5e31de14114f1486c756C30f39Ab9578A92);

    function upgrade(uint256 start, uint256 limit) external {
        require(msg.sender == owner, "Only owner can upgrade");
        address[] memory members = oldContract.getAllMembers();
        uint256 total = members.length;
        if (start >= total) {
            return;
        }
        uint256 end = start + limit;
        if (end > total) {
            end = total;
        }
        for (uint256 i = start; i < end; i++) {
            address _to = members[i];
            if (!isMember[_to]) {
                isMember[_to] = true;
                _members.push(_to);
                emit Minted(_to);
            }
        }
    }

    function getPoint(address account) external view returns (uint256) {
        return calculator.getPoint(account);
    }

    function setCalculator(address _calculator) external {
        require(msg.sender == owner, "Only owner can set calculator");
        calculator = Calculator(_calculator);
    }

    function transferOwnership(address newOwner) external {
        require(msg.sender == owner, "Only owner can transfer ownership");
        owner = newOwner;
    }

    function getMember(uint256 index) external view returns (address) {
        require(index < _members.length, "Invalid member index");
        return _members[index];
    }

    function getMembers(
        uint256 start,
        uint256 limit
    ) external view returns (address[] memory) {
        uint256 total = _members.length;
        if (start >= total) {
            return new address[](0);
        }

        uint256 end = start + limit;
        if (end > total) {
            end = total;
        }

        address[] memory result = new address[](end - start);
        for (uint256 i = start; i < end; i++) {
            result[i - start] = _members[i];
        }
        return result;
    }

    event Minted(address indexed member);

    function mint() external {
        address _to = msg.sender;
        require(!isMember[_to], "Already minted");
        isMember[_to] = true;
        _members.push(_to);
        emit Minted(_to);
    }

    function getMembersCount() external view returns (uint256) {
        return _members.length;
    }
}
