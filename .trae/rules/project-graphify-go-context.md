---
tags: [项目背景, graphify-go]
alwaysApply: false
description: 当需要了解 graphify-go 项目的真实痛点、商业变现目标或 MVP 边界时
---

# Graphify-Go 项目背景与目标 (Context)

## 1. 客户真实诉求与痛点
- **现状 (As-Is)**: 原版 `graphify` (基于 Python 和 `networkx`) 在解析十几万行的大型工程时，耗时极长，内存占用高。更致命的是，它要求使用者必须安装 Python 环境、编译 C 扩展，导致在 Windows 等环境下的安装失败率极高（经常报依赖错误）。
- **痛点 (Pain Points)**: 作为一个底层引擎，它阻碍了上层工具（如 `aictx-cli`）的顺滑分发。用户只想要一个“开箱即用”的代码架构提取器，而不是去折腾 Python 环境。
- **解决方案 (To-Be)**: 用 Go 语言和 `go-tree-sitter` 完全重写底层引擎，并编译为单一二进制文件（Single Binary）。

## 2. 商业变现目标 (Quote-to-Cash)
- 本项目是 `aictx-cli` (Context as Code 工具链) 的核心底座。
- 我们的商业目标是提供**“零 Token 消耗的秒级代码接盘能力”**，吸引独立黑客和全栈外包接单侠使用。
- 只有底层引擎做到极致的快和稳，上层的 `aictx init --onboard` 才能发挥出降维打击的威力。

## 3. 一期 MVP 边界声明 (Out of Scope)
- **坚决不做**：我们不造新的图谱查询语言，不造可视化大屏 (GUI)。可视化交给上层的 Obsidian 或 VSCode 插件处理。
- **坚决不做**：我们不修改原版 `graph.json` 的协议格式。重构的目的是“平替”，而不是“创造新协议”。
- **坚决不做**：我们不为了性能去牺牲准确率。Tree-sitter Query 必须严格对照原版 Python 实现，确保提取的 AST 节点一模一样。
