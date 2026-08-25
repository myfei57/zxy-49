基于 Go 实现的 VenueOps 会展场馆智能化联动平台项目，一款后端服务，完成门禁模式落盘与放行、灯光场景与疏散覆盖、大屏播控、展位用电配额、云台巡看与巡更路线的场馆联动控制。

## 构建与运行

项目使用 Go 1.23 与离线 vendor 依赖构建，运行时只依赖本地文件目录持久化，不依赖 MySQL 或外部中间件。

```bash
go build -mod=vendor -o venueops.exe ./cmd/venueops
./venueops.exe -addr 127.0.0.1:8080 -data ./data -web ./web -seed
```

启动后访问：

- http://127.0.0.1:8080/halls 展馆与展区控制台
- http://127.0.0.1:8080/screens 大屏播控控制台
- http://127.0.0.1:8080/lights 灯光场景控制台
- http://127.0.0.1:8080/guards 巡更路线控制台

## 镜像构建

```bash
bash build_benzhi_docker.sh
```

镜像使用 golang:1.23 基础镜像，离线 vendor 构建，容器默认监听 8080 端口并自动写入演示数据。
