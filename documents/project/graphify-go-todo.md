---
tags: [project, todo, graphify-go]
alwaysApply: false
---

# Graphify-Go 全局开发任务清单 (TODO)

> **目标**: 1:1 像素级复刻 Python 版 `graphify`，实现极速、零依赖的单文件 AST 提取器与代码依赖图谱引擎。

## 阶段一：基础解析器与核心结构 (Parser & Extraction) - [PoC 验证]
- [x] **基建**: 搭建 Go + CGO 编译环境 (Windows MSYS2 GCC / Linux / macOS)
- [x] **基建**: 定义项目核心包结构 (`internal/parser`, `internal/graph`, `internal/cluster`, `internal/export`)
- [x] **基建**: 建立跨语言解析器注册表 (`Registry`)，支持语言探测。
- [x] **提取**: 编写针对 `JavaScript`/`TypeScript` 的 Tree-sitter 查询语句 (`JSQuery`)。
- [x] **提取**: 实现单个文件的 `Extractor` 逻辑，成功提取出类、函数、方法、导入和调用。

## 阶段二：产品文档与规范注入 (Context & Rules) - [已完成]
- [x] **红线**: 梳理并写入 Graphify-Go 的架构设计红线与技术栈选型 (`.trae/rules/`)。
- [x] **PRD**: 提取原版 Python 引擎的魔改特征，重写并固化业务逻辑与产品需求文档。
- [x] **MOC**: 构建并刷新 `documents/` 下的双链路由索引表 (`00-Index.md`)。

## 阶段三：全量文件扫描与多语言支持 (Scanner & Multi-lang) - [已完成]
- [x] **扫描**: 在 `internal/parser` 中实现目录遍历器 (`Scanner`)，支持根据 `.aiignore` 或 `.gitignore` 过滤文件。
- [x] **并发**: 实现基于 Goroutine 的多文件并发提取池 (Worker Pool)，榨干 CPU 性能。
- [x] **语言**: 补充 `Python` 的 Tree-sitter Query 语句 (`PYQuery`)。
- [x] **语言**: 补充 `Go` 的 Tree-sitter Query 语句 (`GOQuery`)。

## 阶段四：图论引擎与社区发现 (Graph Engine & Clustering) - [已完成]
- [x] **图论**: 在 `internal/graph` 中实现加权有向图的数据结构（使用 `gonum/graph` 或手写邻接表）。
- [x] **构建**: 将所有文件的 `ExtractedData` 转换为 `Node` 和 `Edge`。
    - 处理合成节点（如以 `.` 开头的方法占位符和文件节点）。
    - 处理导入边 (`imports`) 和调用边 (`calls`)。
    - 边权重计算逻辑（复刻 Python 版的频次累加）。
- [x] **聚类**: 在 `internal/cluster` 中实现社区发现算法（寻找轻量级的 Louvain 算法 Go 实现，计算节点的 `community` 属性）。

## 阶段五：核心业务分析与导出 (Analysis & Export) - [已完成]
- [x] **分析**: 实现 `God Nodes` (上帝节点) 计算，严格过滤合成节点。
- [x] **分析**: 实现 `Surprising Connections` (意外连接) 算法，分为跨文件策略和跨社区策略。
- [x] **导出**: 在 `internal/export` 中实现 JSON 序列化，严格保证与 Python 版 `graph.json` 的字段 100% 对齐。
- [x] **导出**: 实现 `system-graph.md` 的 Markdown 摘要生成器。

## 阶段六：CLI 入口与跨平台分发 (CLI & Distribution) - [已完成]
- [x] **CLI**: 在 `cmd/graphify` 中编写标准的命令行参数解析 (`flag` 或 `cobra`)，如 `--dir`, `--out`, `--ignore`。
- [x] **测试**: 使用真实的大型开源库作为测试夹具，跑通 E2E 测试并对比 Python 版结果。
- [x] **CI/CD**: 配置 GitHub Actions，实现 Windows, macOS, Linux 的交叉编译矩阵。
- [x] **分发**: 封装 NPM 包结构并发布，供 `aictx-cli` 一键拉取和静默调用。

---
*注：每完成一项任务，请将 `[ ]` 修改为 `[x]`。*
