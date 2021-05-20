# unizone

一个用 Go 写的工具，用于从多个云厂商获取服务信息（包括云主机，内网负载均衡，云数据等），最终生成 DNS Zone 文件，用于搭建内网 DNS

## 关键设计

**`unizone` 不会构建任何形式的命名空间**

现阶段各个厂商均提供类似云联网的功能，跨多个地域，甚至跨厂商能够实现内网打通，在这种情况下，极简设计反而是最好的。

设想这样一个场景：

在 云厂商 `cloud-a` 的 `vpc-a` 部署有一个云服务器 `server-a` 和一个云数据库 `mysql-a`

在 云厂商 `cloud-b` 的 `vpc-b` 部署有一个云服务器 `server-b` 和一个云数据库 `mysql-b`

如果使用带命名空间的设计方式，最终会生成类似如下形式的 DNS 记录

```text
server-a.cvm.vpc-a.cloud-a.infra
mysql-a.mysql.vpc-a.cloud-a.infra
server-b.cvm.vpc-b.cloud-b.infra
mysql-b.mysql.vpc-b.cloud-b.infra
```

那么任何 `VPC` 之间，云厂商之间，甚至于 自建数据库 与 云数据库之间的切换，都会导致上述 DNS 条目变化，最终会带来配置变更成本，那么自建 DNS 就失去了意义。

而如果使用全局扁平设计，不同厂商，不同 `VPC` 以及不同类型的服务全部位于同一个扁平的命名空间下

```text
server-a.infra
server-b.infra
mysql-a.infra
mysql-b.infra
```

那么，把一个服务从厂商 A 迁移到厂商 B，从一个 VPC 移动到另一个 VPC，从 自建服务 切换到 托管服务，都是零成本的。

**你唯一需要做的事情就是，打通网络，并制定一个命名规范，避免冲突即可**

## 支持厂商

## 许可证

Guo Y.K., MIT License
