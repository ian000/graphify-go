---
tags: [技术架构, graphify-go]
alwaysApply: false
description: 当需要确定 graphify-go 项目的前后端框架、依赖管理、分发渠道或 CI/CD 选型时
---

# Graphify-Go 技术栈选型与部署基建 (Tech Stack)

## 1. 核心语言与编译选型
- **语言 (Language)**: Go 1.21+。原因：极速并发、静态编译、内存占用低。
- **构建工具 (Builder)**: `go build`，并启用 CGO。
- **CGO (C bindings)**: 由于 `go-tree-sitter` 底层是 C 编写的跨语言解析器，必须配置 GCC (如 Windows 的 MinGW-w64) 并在编译时开启 `CGO_ENABLED=1`。这是本项目最重要的工程依赖。

## 2. 核心架构选型
- **AST 解析层 (Parser)**: `github.com/smacker/go-tree-sitter`，并且静态链接各种语言的解析器（如 `javascript`, `python`, `golang` 等）。这是为了达到 1:1 还原 Python 版解析结果。
- **图论计算层 (Graph)**: 需要平替 Python 的 `networkx`。可考虑使用 `github.com/gonum/graph` 或者手写轻量级的邻接表（如果只需要做度中心性和简单遍历）。
- **社区发现层 (Clustering)**: Python 版使用了 `leidenalg`（Leiden 算法）。Go 生态中如果要寻找轻量平替，可实现 Louvain 算法（如使用社区开源库或自行实现），以计算图谱节点的 `community` 属性。
- **JSON 序列化层 (Serializer)**: Go 标准库 `encoding/json`，必须确保导出的格式符合 Python 版 `export.py` 的结构。

## 3. 部署与分发选型
- **分发渠道 (Distribution)**:
  - 核心目标是**单文件二进制分发** (Single Binary)。
  - 通过 GitHub Actions 配置交叉编译矩阵：
    - Windows (amd64) -> `.exe`
    - macOS (amd64, arm64) -> Mach-O
    - Linux (amd64, arm64) -> ELF
- **与 Node.js CLI 的集成 (Integration)**:
  - 编译后的二进制包将被发布到 NPM (如 `@aictx/graphify-win32-x64`)。
  - `aictx-cli` 的 Node.js 代码通过 `execa` 子进程调用这个隐藏的二进制文件，通过标准输出 (Stdout) 获取执行进度，通过文件系统 (File System) 读取最终的 `graph.json` 结果。

## 4. 测试与验证选型
- **单元测试 (Unit Testing)**: Go 标准库 `testing`。
- **端到端测试 (E2E Testing)**: 使用真实的 JS/TS 或 Python 源码库作为 fixture，运行 Python 原版和 Go 重写版，对比输出的 `graph.json` 结构是否一致（尤其是节点和边的数量、权重）。
