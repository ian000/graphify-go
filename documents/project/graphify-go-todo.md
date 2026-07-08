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

## 阶段七：五语言完整支持执行计划 (JS/TS/Python/Go/Verilog) - [进行中]
> 对齐产品计划：[[graphify-go-language-full-support-todo]]
> 执行红线：不改 `graph.json` 协议、不破坏现有功能、不扩展五门目标语言之外的新语言。

### 7.1 回归护栏与基线冻结
- [ ] **基线**: 固定 JS/TS/Python/Go 夹具，生成并提交当前 `graph.json` 快照基线。
- [ ] **对比**: 增加节点数、边数、关系类型分布、关键节点 ID 的自动 diff 脚本。
- [ ] **测试**: 补齐 TypeScript 与 Python 的 `extractor` 测试，达到与 JS/Go 同等级覆盖。
- [ ] **门禁**: CI 增加 `go test ./...` + 快照对比双门禁。

### 7.2 JavaScript 完整支持
- [ ] **声明**: 覆盖函数表达式、箭头函数、对象方法、导出声明。
- [ ] **调用**: 覆盖 `call/new`、链式调用、可选链调用。
- [ ] **依赖**: 覆盖 `import`、`require`、动态导入。
- [ ] **验收**: Node 工程夹具回归通过，快照无异常放大。

### 7.3 TypeScript 完整支持
- [ ] **语法**: 覆盖 `interface/type alias/enum/namespace` 的提取与映射。
- [ ] **类方法**: 覆盖访问修饰符、泛型、抽象类相关形态。
- [ ] **依赖**: 覆盖 `import type` 与路径别名常见场景。
- [ ] **验收**: TS 项目夹具无明显漏提取、误提取。

### 7.4 Python 完整支持
- [ ] **声明**: 覆盖嵌套函数、类方法/静态方法、装饰器函数。
- [ ] **调用**: 覆盖属性调用、链式调用、上下文对象调用。
- [ ] **依赖**: 覆盖 `import` 与 `from ... import ...` 并统一模块名提取。
- [ ] **验收**: Python 夹具快照稳定，统计口径不漂移。

### 7.5 Go 完整支持
- [ ] **声明**: 覆盖函数、方法、接口、结构体及关键类型声明映射。
- [ ] **调用**: 覆盖选择器调用、包函数调用及常见链路模式。
- [ ] **依赖**: 覆盖标准导入、别名导入、多段导入块。
- [ ] **验收**: 中大型 Go 仓库性能与准确率不倒退。

### 7.6 Verilog 完整支持
- [ ] **语法**: 覆盖 `module/function/task/always/instance` 主干结构提取。
- [ ] **映射**: 定义 Verilog 到统一图模型的节点与关系映射规则并固化测试。
- [ ] **依赖**: 覆盖 `` `include`` 与常见跨文件引用路径清洗。
- [ ] **验收**: 在 RTL 夹具上通过准确率与回归稳定性校验。

### 7.7 收敛与发版
- [ ] **文档**: 更新 README 与 PRD 的语言覆盖说明与边界声明。
- [ ] **记录**: 每阶段输出“新增覆盖点 + 已知限制 + 回归结果”。
- [ ] **准入**: 测试全绿、快照稳定、性能阈值通过后发版。

---
*注：每完成一项任务，请将 `[ ]` 修改为 `[x]`。*
