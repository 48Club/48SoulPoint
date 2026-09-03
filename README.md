# 48SoulPoint

48Club Soul Point 查询服务：定时从链上抓取成员积分写入 MySQL 形成每日快照，并以内存缓存对外提供实时积分与历史明细查询。

## 接口

所有请求的 `Host` 必须为 `soul-api.48.club`，否则返回 400。响应统一为 `{"code": <http status>, "message": "", "data": ...}`。

### `GET /`

| 参数 | 说明 |
|---|---|
| 无 | 全部 `points > 0` 的成员，按 `points` 降序 |
| `address` | 指定地址（校验和格式），无积分时 `message` 为 `no points`、`data` 为空数组 |
| `address` + `detail=true` | 在实时积分基础上附带每日快照明细 `detail` |

`points` / `level` 来自链上实时缓存（每 15 分钟刷新）；`detail` 来自数据库快照缓存，按快照总数取最近 48 / 7 / 2 次。

```jsonc
// GET /?address=0x...&detail=true
{
  "code": 200,
  "message": "",
  "data": {
    "address": "0x...",
    "count": 48,            // detail 条数
    "level": "Gold",
    "points": 520,
    "detail": [
      { "snap_time": "2026/09/03", "points": 520, "koge": 100, "stake": 300, "nft": 0, "bsc_stake": 120 }
    ]
  }
}
```

服务刚启动、实时缓存尚未构建完成时最多等待 30 s，超时返回 `503` 并携带 `Retry-After`。

### `GET /levels`

等级阈值表：Member 1–47、Entry 48–479、Gold 480–709、Platinum ≥ 710。

## 配置

默认读取工作目录下的 `config.json`，可用环境变量 `SP_CONFIG` 指定路径。仓库内的 `config.json` 为模板。

| 字段 | 说明 |
|---|---|
| `database` | MySQL 连接信息（`username` / `password` / `host` / `port` / `database`） |
| `rpc` | BSC JSON-RPC 端点，用于 multicall 查询积分与按时间定位区块 |
| `listen` | 监听地址，如 `:9091` |
| `test_addr` | 测试钱包，直接并入结果；字段 `addr` 及 `kogepoint` / `stakepoint` / `nftpoint` / `bscstakepoint` / `govbnbpoint`，缺省为 0 |

快照的"日"边界依赖系统时区，部署时请确认服务器时区符合预期。

## 部署

```sh
go install .            # 产出 $(go env GOPATH)/bin/sp
sudo ./install.sh       # 生成 systemd 单元 sp.service，配置复制到 ~/.config/sp/
```

## 开发

```sh
go build ./... && go vet ./...
go test ./...           # 单元测试，不依赖外部服务

# 集成测试：需要可用的 config.json；test/db 会向数据库真实写入一次快照
SP_CONFIG=/path/to/config.json go test -tags integration ./ethclient/ ./test/db/

# 合约 ABI 变更后重新生成绑定
abigen --abi contracts/sol/abi.json --pkg SoulPoint_48Club --out contracts/SoulPoint_48Club/SoulPoint_48Club.go
```
