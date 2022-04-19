# Tendermint

![banner](docs/tendermint-core-image.jpg)

[拜占庭容错](https://en.wikipedia.org/wiki/Byzantine_fault_tolerance)
[状态机](https://en.wikipedia.org/wiki/State_machine_replication).
或者 [区块链](<https://en.wikipedia.org/wiki/Blockchain_(database)>), 简单点.


[![version](https://img.shields.io/github/tag/tendermint/tendermint.svg)](https://github.com/tendermint/tendermint/releases/latest)
[![API Reference](https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667)](https://pkg.go.dev/github.com/tendermint/tendermint)
[![Go version](https://img.shields.io/badge/go-1.15-blue.svg)](https://github.com/moovweb/gvm)
[![Discord chat](https://img.shields.io/discord/669268347736686612.svg)](https://discord.gg/AzefAFd)
[![license](https://img.shields.io/github/license/tendermint/tendermint.svg)](https://github.com/tendermint/tendermint/blob/master/LICENSE)
[![tendermint/tendermint](https://tokei.rs/b1/github/tendermint/tendermint?category=lines)](https://github.com/tendermint/tendermint)
[![Sourcegraph](https://sourcegraph.com/github.com/tendermint/tendermint/-/badge.svg)](https://sourcegraph.com/github.com/tendermint/tendermint?badge)

| Branch | Tests                                                                                                                                                                                                                                                  | Coverage                                                                                                                             | Linting                                                                    |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------- |
| master | [![CircleCI](https://circleci.com/gh/tendermint/tendermint/tree/master.svg?style=shield)](https://circleci.com/gh/tendermint/tendermint/tree/master) </br> ![Tests](https://github.com/tendermint/tendermint/workflows/Tests/badge.svg?branch=master) | [![codecov](https://codecov.io/gh/tendermint/tendermint/branch/master/graph/badge.svg)](https://codecov.io/gh/tendermint/tendermint) | ![Lint](https://github.com/tendermint/tendermint/workflows/Lint/badge.svg) |

Tendermint Core 是 采用状态转换机的拜占庭容错 (BFT) 中间件 - 可以用任何编程语言编写 （现在是golang）-
还有 可以在多台机器上安全地复制.

查看更多协议细节, 点  [这里](https://github.com/tendermint/spec).

有关共识协议的详细分析，包括安全性和活性证明，
看我们最近的论文, "[The latest gossip on BFT consensus](https://arxiv.org/abs/1807.04938)".

## Releases

请不要使用 master 作为您的生产 。用这里 [releases](https://github.com/tendermint/tendermint/releases) instead.

Tendermint 正在私人和公共环境中用于生产,
最值得注意的是 [Cosmos Network](https://cosmos.network/).
但是，我们仍在对协议和 API 进行重大更改，尚未发布 v1.0.
有关更多详细信息，请参见下文 [versioning](#versioning).

无论如何，如果您打算在生产环境中运行 Tendermint，我们很乐意为您提供帮助。 你可以
邮件我们 [over email](mailto:hello@interchain.berlin) 或 [加入disscoed](https://discord.gg/AzefAFd).

## 安全性

报告安全漏洞, 看看 [漏洞奖励程序](https://hackerone.com/tendermint). 
有关我们正在寻找的bug错误的示例，请参阅 [我们的安全政策](SECURITY.md)

我们还维护一个专门的安全更新邮件列表。我们只会使用这个邮件列表去
通知您 Tendermint Core 中的漏洞和修复。你可以订阅[这里](http://eepurl.com/gZ5hQD).

## 最低要求

| Requirement | Notes            |
| ----------- | ---------------- |
| Go 版本  | Go1.15 or higher |

## Documentation

完整的文档可以在 [这里] (https://docs.tendermint.com/master/) 找到.

### 安装

请参阅[安装说明](/docs/introduction/install.md).

### 快速开始

- [单节点](/docs/introduction/quick-start.md)
- [使用 docker-compose 的本地集群](/docs/networks/docker-compose.md)
- [使用 Terraform 和 Ansible 的远程集群](/docs/networks/terraform-and-ansible.md)
- [加入 Cosmos 测试网](https://cosmos.network/testnet)

## 贡献

请在所有互动中遵守[行为准则](CODE_OF_CONDUCT.md)。

在为项目做出贡献之前，请查看 [贡献指南](CONTRIBUTING.md)
和 [风格指南](STYLE_GUIDE.md). 
还有[规格](https://github.com/tendermint/spec), 观看 [开发者会议](/docs/DEV_SESSIONS.md), 
并熟悉我们的[架构决策记录](https://github.com/tendermint/tendermint/tree/master/docs/architecture).

## 版本控制

### 语义版本控制

Tendermint 使用 [语义版本控制](http://semver.org/) 来确定版本何时以及如何更改。
根据 SemVer，公共 API 中的任何内容都可以在 1.0.0 版本之前随时更改

在这 0.X.X 天内为 Tendermint 用户提供一些稳定性, 使用次要版本
在整个公共 API 的一个子集中发出重大更改信号。这个子集包括所有
暴露给其他进程（cli、rpc、p2p 等）的接口，但不
包括 Go API。

也就是说，以下包中的重大更改将记录在
CHANGELOG 即使它们不会导致次要版本变动：

- crypto
- config
- libs
    - bech32
    - bits
    - bytes
    - json
    - log
    - math
    - net
    - os
    - protoio
    - rand
    - sync
    - strings
    - service
- node
- rpc/client
- types

### 升级
为了避免在 1.0.0 之前积累技术债务，
我们不保证重大更改（即次要版本中的颠簸）
将与现有的 Tendermint 区块链一起使用。在这些情况下，您将
必须启动一个新的区块链，或者编写一些自定义的东西来获取旧的
数据进入新链。但是，补丁版本中的任何不同都应该是
与现有的区块链历史兼容。
有关升级的更多信息，请参阅 [UPGRADING.md](./UPGRADING.md)。

### 支持的版本

因为我们是一个小型核心团队，所以我们只发布补丁更新，包括安全更新，
到最近的次要版本和第二次最近的次要版本。最后，
我们强烈建议让 Tendermint 保持最新。可以找到升级说明
在 [UPGRADING.md](./UPGRADING.md) 中。

## 资源

### Tendermint Core

有关区块链数据结构和 p2p 协议的详细信息, 查看
[Tendermint specification](https://docs.tendermint.com/master/spec/).

有关使用该软件的详细信息，请参阅 [文档](/docs/)，它也是
在: <https://docs.tendermint.com/master/>

### 工具

基准测试由 [`tm-load-test`](https://github.com/informalsystems/tm-load-test).
其他工具可以在[/docs/tools](/docs/tools).

### 应用程序

- [Cosmos SDK](http://github.com/cosmos/cosmos-sdk); a cryptocurrency application framework
- [Ethermint](http://github.com/cosmos/ethermint); Ethereum on Tendermint
- [Many more](https://tendermint.com/ecosystem)

### 继续研究
- [The latest gossip on BFT consensus](https://arxiv.org/abs/1807.04938)
- [Master's Thesis on Tendermint](https://atrium.lib.uoguelph.ca/xmlui/handle/10214/9769)
- [Original Whitepaper: "Tendermint: Consensus Without Mining"](https://tendermint.com/static/docs/tendermint.pdf)
- [Blog](https://blog.cosmos.network/tendermint/home)
