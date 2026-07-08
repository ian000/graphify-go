---
tags:
  - aictx
  - moc
  - product
aliases:
  - [product Index, 目录]
entities:
  - [MOC]
roles:
  - [Maintainer]
---
# Product Map of Content (MOC)

> **路由表**: 这是 `product` 目录的核心索引文件。
> AI 助手在寻找特定业务逻辑时，必须**优先且仅**读取此文件，并通过这里的双向链接（如 `[[xxx]]`）去跳转到对应的原子文档。**严禁使用全局检索**。

## 📑 领域索引

<!-- aictx-index-start -->

| Doc Route (Bi-link) | Core Entities | Aliases | Description |
| --- | --- | --- | --- |
| [[graphify-go-PRD]] | `AST,Zero-LLM,Tree-sitter` | `Graphify-Go PRD,产品需求文档` | Graphify-Go 产品需求文档 (PRD) |
| [[graphify-go-language-full-support-todo]] | `AST,Tree-sitter,ExtractedData,Graph` | `Graphify-Go Language Full Support TODO,多语言完整支持计划` | Graphify-Go 五语言完整支持开发计划 (TODO) |

<!-- aictx-index-end -->

## 📌 业务模块说明
- `[[graphify-go-PRD]]`: 定义产品目标、核心红线与模块职责。
- `[[graphify-go-language-full-support-todo]]`: 承接 PRD 的多语言能力目标，拆解为可执行阶段任务与发布验收门槛。
