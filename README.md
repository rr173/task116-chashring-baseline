# task116-chashring

加权一致哈希环管理服务（Go）。包含虚拟节点、副本查询、再平衡预览、SQLite 持久化与重启恢复。

运行自检：

```bash
go run . --smoke-test
```

启动服务：

```bash
go run . server --addr :8080
```

构建并运行（多架构）：

```bash
bash ./build_benzhi_docker.sh go-task116-chashring:amd64 linux/amd64
docker run --rm go-task116-chashring:amd64 --smoke-test
```
