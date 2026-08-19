# Benzhi 评测镜像说明

本项目是一个纯标准库 + SQLite 的 Go 服务，实现加权一致哈希环管理（虚拟节点、副本查询、再平衡预览、SQLite 持久化与重启恢复）。

## 项目用途

服务支持创建命名环、增删节点、按 key 查询主节点与副本节点、查看负载分布与均衡度、预览/执行拓扑变更、按 key 持久化与重启恢复。HTTP 入口提供 `/healthz`、环与节点的增删查、查询与统计接口；`--smoke-test` 会执行不依赖外部服务的端到端自检并自行退出。

## 标准命令

```bash
go build ./...
go run . --smoke-test
go test ./...
go vet ./...
go run . server --addr :8080
```

## Benzhi Docker 构建

`build_benzhi_docker.sh` 固定使用 `benzhi.Dockerfile`，接受镜像名和目标平台两个参数，默认分别为 `my-project` 与 `linux/amd64`：

```bash
bash ./build_benzhi_docker.sh go-task116-chashring:amd64 linux/amd64
bash ./build_benzhi_docker.sh go-task116-chashring:arm64 linux/arm64
```

镜像使用 `docker.m.daocloud.io/library/golang:1.26.3-bookworm`，构建阶段会按 `go.mod`/`go.sum` 下载依赖并执行 `go build ./...`，容器启动后进入 bash：

```bash
docker run -it go-task116-chashring:amd64
```

项目自身的 `Dockerfile` 用于构建可运行的最小 HTTP 服务镜像，支持 `--smoke-test` 作为容器自检入口。
