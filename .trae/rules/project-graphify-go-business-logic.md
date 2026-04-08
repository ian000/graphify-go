---
tags: [业务逻辑, graphify-go]
alwaysApply: false
description: 当开发或讨论 graphify-go 项目的核心业务逻辑、状态机或红线拦截时
---

# Graphify-Go 业务逻辑与红线 (Business Logic)

## 1. 核心领域模型 (Domain Models)
- **AST Node (抽象语法树节点)**: 通过 `go-tree-sitter` 提取出的代码结构，如 Class, Function, Method, Import, Call。
- **Graph Node (图谱节点)**: 将 AST 节点转换成的图论节点，包含属性：`id`, `label`, `type`, `community`, `source_file`。
- **Synthetic Node (合成节点)**: 提取器自动生成的用于表示“文件”本身或“方法占位符”的节点（如 `.login()`）。这些节点在图谱分析时必须被特殊处理（过滤）。
- **Graph Edge (图谱边)**: 节点之间的关系，包含 `source`, `target`, `weight`, `type`（如 `contains`, `calls`, `imports`）。
- **God Node (上帝节点)**: 具有最高连接度（Degree）的**真实**业务节点，代表系统的核心抽象。
- **Surprising Connection (意外连接)**: 跨文件或跨社区的异常高权重边，通常代表潜在的技术债或高危耦合。

## 2. 数据流转状态机 (State Machine)
1. **解析阶段 (Parsing)**: `Registry` 探测文件类型 -> 分配 Parser -> 生成 AST。
2. **提取阶段 (Extraction)**: `Extractor` 执行 Query -> 捕获节点和依赖 -> 生成 `ExtractedData`。
3. **图谱构建阶段 (Graph Building)**: 将多个文件的 `ExtractedData` 合并 -> 生成有向图 -> 计算权重。
4. **聚类与分析阶段 (Clustering & Analysis)**: 执行社区发现 (Louvain) 划分模块 -> 计算上帝节点 -> 寻找意外连接。
5. **导出阶段 (Exporting)**: 生成完全兼容原版 Python 的 `graph.json` 和摘要 `system-graph.md`。

## 3. 业务红线 (Business Guardrails)
- **【硬阻断】禁止破坏 JSON 契约**: 输出的 `graph.json` 必须 100% 兼容原版 Python 的数据结构。任何新增字段可以添加，但绝不能修改或删除原有的 `id`, `source`, `target`, `community` 等核心键名。
- **【硬阻断】必须过滤合成节点**: 在计算 God Nodes 和 Surprising Connections 时，必须严格复刻原版 `analyze.py` 的逻辑，排除文件级中心节点和以 `.` 开头的方法占位符节点，防止图谱分析被垃圾数据污染。
- **【硬阻断】零运行时依赖**: 所有的 CGO (Tree-sitter) 必须在编译期静态链接，绝不能要求最终用户安装 GCC 或 Python。
