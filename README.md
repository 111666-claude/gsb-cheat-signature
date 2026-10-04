# cheat-signature

作弊特征扫描：用已知特征集合匹配玩家上报的特征，命中记账，只用 Go 标准库。

```
go run ./cmd/cheat-signature --sample=stale
go run ./cmd/cheat-signature --sample=dup
go run ./cmd/cheat-signature --sample=partial
go run ./cmd/cheat-signature --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **精确匹配**：只有上报特征与登记特征完全相同才算命中，部分包含不算。
- **特征表实时**：每次扫描都用当前登记的特征表，新登记的特征立刻生效，不缓存旧表。
- **幂等**：同一个 `sampleID` 只记一次，重复扫描返回命中但不重复记账。
- **命中台账**：`Hits` 返回该玩家命中过的不同特征 id（升序、去重）。
- **不变量**：命中只由当前特征表与样本特征决定；同一串扫描重复执行结果相同。
- **代价**：扫描不许全扫特征表，`Scanned` 不随特征数乘扫描次数增长；10 万条特征、每秒 10 万次扫描。

## 输出契约（不改格式）

```
hits=..
scanned=..
```

## 目录

```
signatures.go          特征匹配、实时特征表、幂等与命中台账
signatures_test.go     单元测试
cmd/cheat-signature/   命令行入口
```
