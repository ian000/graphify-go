---
tags:
  - aictx
  - architecture
  - graphify-go
aliases:
  - [Graphify-Go 技术架构设计, 架构纪实]
entities:
  - [Tree-sitter, CGO, Graph, Leiden/Louvain]
roles:
  - [Architect, Backend Developer]
---

# Graphify-Go 技术架构设计 (Architecture)

> **文档状态**: 纪实版 (As-Is & To-Be Mapping)。这份技术架构设计文档明确了在 `graphify-go` 的开发过程中，从语言选型到底层图论算法的技术债排雷，以及如何 1:1 复刻 Python 原版的魔改逻辑。

## 1. 核心语言与编译架构

- **语言**: Go 1.21+。
- **并发模型**: 利用 Goroutine 实现多文件 AST 并行提取，突破 Python 版 GIL (全局解释器锁) 带来的 CPU 瓶颈。
- **CGO 依赖**: 核心解析器 `github.com/smacker/go-tree-sitter` 强依赖 C 语言绑定。
  - **技术债预警**: 必须配置跨平台的 C 编译器（如 MinGW-w64）。最终分发时，所有的 C 代码必须被静态链接到 Go 二进制文件中，实现对终端用户的“零依赖”。

## 2. 架构分层设计 (Layered Design)

### 2.1 提取引擎层 (Parser & Extractor Layer)
- **实现路径**: `internal/parser/`
- **核心逻辑**:
  - `Registry`: 基于文件扩展名（`.js`, `.py`, `.go`）的动态路由分发，将每个文件映射到对应的 `tree-sitter.Language` 实例。
  - `Extractor`: 执行预先编写好的 Tree-sitter Queries (如 `queries.go`)。
  - **魔改特征映射**: 这里的 Queries 必须精准对应原版 Python 的节点标签：`@class.name`, `@function.name`, `@method.name`, `@call.function`, `@import.source`。

### 2.2 图论模型与计算层 (Graph Construction Layer)
- **实现路径**: `internal/graph/`
- **核心挑战**: 在 Go 生态中寻找或实现 Python `networkx` 的平替。
  - **选型建议**: 如果 `gonum/graph` 过于庞大且难以序列化为我们所需的 JSON 格式，建议手动实现一个轻量级的带权重的邻接表 (Adjacency List) 或邻接矩阵 (Adjacency Matrix)。
- **魔改特征映射**:
  - **Synthetic Nodes (合成节点)**: 在建图时，遇到文件导入（如 `./utils`），需要自动生成一个代表该文件的虚拟节点，以模拟代码库的物理层级。这些节点在后续分析时需要被标记。
  - **边权重 (Edge Weighting)**: 当文件 A 多次调用文件 B 的不同方法时，连接这两者的边的权重 (Weight) 必须进行累加。

### 2.3 社区聚类与洞察层 (Clustering & Analysis Layer)
- **实现路径**: `internal/cluster/` & `internal/graph/analyze.go`
- **核心挑战**: 原版 Python 依赖 `leidenalg` (Leiden 算法，基于 C++ 的 igraph) 进行模块划分。
  - **技术债预警**: Go 生态缺乏极其成熟且轻量级的 Leiden 算法库。
  - **降级方案**: 如果找不到合适的 Leiden 实现，可采用标准的 Louvain 算法计算模块度 (Modularity)，确保最终生成的每个节点都被打上 `community` 标签。
- **魔改特征映射**:
  - **God Nodes 过滤**: 必须复写 `_is_file_node` 和 `_is_concept_node` 逻辑。在对节点进行 Degree 排序前，剔除掉那些仅仅因为是被大量导入而变成中心的“文件占位符”，确保 God Nodes 都是真正的业务实体。
  - **Surprising Connections**: 通过遍历高权重的边，找出那些连接了两个不同 `community` 的调用关系。

### 2.4 序列化与持久化层 (Export Layer)
- **实现路径**: `internal/export/`
- **魔改特征映射**:
  - 必须导出一个严格的 `graph.json` 结构：
    ```json
    {
      "nodes": [
        {"id": "...", "label": "...", "type": "...", "community": 1, "source_file": "..."}
      ],
      "links": [
        {"source": "...", "target": "...", "weight": 1.5, "type": "calls"}
      ]
    }
    ```
  - **系统图谱摘要 (system-graph.md)**: 这个文件是给 AI 读的“上帝提示词”底座，必须高度浓缩提取出的上帝节点和社区划分情况。

## 3. 部署与分发机制 (Distribution Mechanism)

**核心红线**: `graphify-go` 不能是一个让用户手动编译的开源项目，它必须以即插即用的方式被 `aictx-cli` 唤起。

- **交叉编译矩阵 (Cross-Compilation Matrix)**:
  在 GitHub Actions 中配置 CI 流水线，针对 `windows-amd64`, `darwin-arm64`, `linux-amd64` 分别执行 `go build` 并生成对应的二进制包。
- **NPM 分发桥接**:
  把这些二进制文件发布为 NPM 包（或通过 `postinstall` 脚本下载），使得前端 CLI 工具可以通过 `execa` 启动一个隐藏子进程来执行 AST 提取，并在标准输出流中监听进度。
