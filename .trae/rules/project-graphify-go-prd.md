---
tags: [PRD, graphify-go]
alwaysApply: false
description: 当设计或开发 graphify-go 项目的核心功能模块 (MOC) 和角色权限 (RBAC) 时
---

# Graphify-Go 产品需求文档 (PRD)

## 1. 核心功能模块列表 (MOC)

### 1.1 跨语言 AST 提取引擎 (Parser Layer)
- **目标**: 根据文件后缀自动探测语言并加载对应的 `go-tree-sitter` 解析器。
- **功能点**:
  - `Registry` 管理 `javascript`, `typescript`, `python`, `go`, `rust`, `java`, `c`, `cpp` 等多语言。
  - `Extractor` 负责执行预先定义的 Tree-sitter Queries，提取出 `Class`, `Function`, `Method`, `Import`, `Call` 节点。

### 1.2 依赖关系图谱构建 (Graph Layer)
- **目标**: 将各个文件中提取到的节点和边数据合并，构建有向加权连通图。
- **功能点**:
  - `AddNode`: 添加真实节点和合成占位符节点。
  - `AddEdge`: 根据提取的调用 (`Calls`) 和导入 (`Imports`) 关系，创建带权重和类型的边。
  - `WeightCalculation`: 根据调用频率计算边权重。

### 1.3 社区发现与图谱分析 (Analysis & Clustering)
- **目标**: 复刻原版 Python `analyze.py` 的算法，识别代码社区、上帝节点和异常依赖。
- **功能点**:
  - `CommunityDetection`: 实现 Louvain 算法划分代码模块。
  - `GodNodes`: 根据节点的度 (Degree) 计算出核心业务抽象（上帝节点），并剔除合成文件节点。
  - `SurprisingConnections`: 计算跨文件或跨社区的意外高权重依赖，发现技术债。

### 1.4 格式化导出器 (Exporter)
- **目标**: 将构建和分析后的图谱序列化并保存。
- **功能点**:
  - `ExportJSON`: 生成标准的 `graph.json`。
  - `GenerateReport`: 生成包含图谱摘要、上帝节点和意外连接的 `system-graph.md`，用于给 AI 喂养高浓度上下文。

## 2. 角色与权限划分 (RBAC)
由于 `graphify-go` 是一个本地 CLI 工具和底层的执行引擎，它不涉及复杂的用户权限系统。
- **调用者 (Caller)**: 通常是 `aictx-cli` 或者最终用户的终端环境，拥有当前目录的文件读写权限。
- **系统权限 (System Permissions)**: `graphify-go` 只能在指定的项目目录下执行只读（读取代码）和在特定的 `out_dir` 执行写操作（保存报告），绝不执行任何越权的系统调用。
