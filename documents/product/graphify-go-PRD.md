---
tags:
  - aictx
  - PRD
  - graphify-go
aliases:
  - [Graphify-Go PRD, 产品需求文档]
entities:
  - [AST, Zero-LLM, Tree-sitter]
roles:
  - [Product Manager]
---

# Graphify-Go 产品需求文档 (PRD)

> **文档状态**: 纪实版 (As-Is & To-Be Mapping)。本文档基于 `aictx-cli` 团队对原版 Python `graphify` 的“魔改”经验，提炼出 Go 重写版必须严格复刻的核心产品能力与防腐红线。

## 1. 产品定位 (Product Positioning)

**Graphify-Go** 是一款专为 AI 辅助编程（尤其是 Context as Code 理念）打造的**纯本地、零依赖、极速 AST 依赖图谱提取引擎**。

它不再是一个面向普通人类的 GUI 可视化工具，而是作为 `aictx-cli` (或其他 AI IDE 插件) 的底层“物理探针”。它的唯一使命是：**在几百毫秒内，将数十万行的遗留代码（屎山）压缩成一份结构化的物理拓扑图谱，从而彻底消除大语言模型（LLM）在阅读全局代码时的上下文爆炸与幻觉。**

## 2. 核心“魔改”特征继承 (Core Features Inherited)

在原版的 Python 实践中，我们发现传统的代码调用图对于 AI 来说噪音太大，因此我们加入了一系列“魔改”逻辑。Go 版本必须 100% 继承这些产品特性：

### 2.1 智能的上帝节点提纯 (God Nodes Purification)
- **痛点**: 原版的连通图会把诸如 `utils.py` 或 `models/` 这样的“文件级节点”，以及 `__init__()` 这样的“高频基础方法”误认为是系统的核心架构（God Nodes），这会严重误导 AI 的判断。
- **产品要求**: `analyze.go` 必须复刻过滤逻辑 (`_is_file_node` 和 `_is_concept_node`)。在计算出度/入度 (Degree) 时，**必须排除合成的占位符节点和文件中心节点**，确保交给 AI 的 God Nodes 都是真实的业务实体（如 `UserManager` 或 `OrderService`）。

### 2.2 跨维度的意外连接洞察 (Surprising Connections)
- **痛点**: AI 需要知道系统的技术债在哪里，而不是常识性的调用。
- **产品要求**: 
  - 对于**多文件项目**：必须能跨文件识别高权重的调用关系，挖掘出“不该耦合的模块耦合在了一起”。
  - 对于**单文件项目**：利用介数中心性 (Betweenness Centrality) 或跨社区 (Cross-Community) 权重，找出桥接图谱两端的隐藏结构。

### 2.3 零配置的多语言平滑降级 (Multi-lang & Fallback)
- **产品要求**: 无需用户指定语言，系统通过遍历目录，利用文件后缀自动分配 `go-tree-sitter` 的 Parser（首批支持 JS/TS/Py/Go）。遇到无法解析的文件静默跳过，绝不阻塞整个图谱的生成流程。

## 3. 功能模块规划 (Functional Modules)

| 模块名称 | 核心职责 | 输入 | 输出 |
| :--- | :--- | :--- | :--- |
| **Scanner (遍历器)** | 快速遍历工作区目录，严格遵循 `.aiignore` / `.gitignore` 的过滤规则。 | `cwd` 路径 | `[]string` 文件路径列表 |
| **Parser (解析器)** | 自动嗅探语言，调用对应的 Tree-sitter Parser，提取类、函数、导入、调用。 | 文件路径 + 字节码 | `ExtractedData` 结构体 |
| **Graph Builder (建图器)** | 将离散的 AST 数据转换为统一的有向加权图，处理合成节点，累加边权重。 | `[]ExtractedData` | 内存图谱模型 |
| **Analyzer (分析器)** | 运行 Louvain 社区发现算法，计算上帝节点与意外连接。 | 内存图谱模型 | 节点社区属性 & 洞察指标 |
| **Exporter (导出器)** | 严格按照原版协议导出 `graph.json`，并生成 `system-graph.md` 摘要。 | 图谱与洞察指标 | `JSON` 和 `Markdown` 文件 |

## 4. 成功指标 (Success Metrics)

1. **零依赖安装成功率**: 100%。通过发布单一的跨平台二进制文件，彻底终结 Python 环境带来的安装失败噩梦。
2. **提取速度跃升**: 相比原版 Python，解析 10 万行代码的时间应从 30 秒级缩减至 **3 秒以内**。
3. **输出无损兼容**: 生成的 `graph.json` 能够被现有的 `aictx-cli` 路由和渲染逻辑完美读取，不抛出任何 schema 错误。
