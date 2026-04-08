# Graphify-Go 🚀

*[Read in English](./README.md) | [中文阅读](./README_zh.md)*

**极速、零依赖的本地 AST 提取与代码架构图谱引擎。**

专为代码分析、重构评估以及 AI 辅助编程脚手架设计，提供极致的解析性能与无缝的跨平台分发体验。

### 💡 为什么开发 Graphify-Go？

在代码架构分析领域，传统工具往往面临几个致命缺陷：
- **依赖地狱**：许多基于 Python 的解析工具要求用户全局安装 Python 3.9+、pip，还要解决 `networkx` 和 `leidenalg` (底层依赖 C++) 等库带来的复杂编译报错。
- **性能瓶颈**：面对十万行级别的大型项目，传统的单线程脚本受限于内存开销，AST 提取往往需要数十秒甚至分钟级。
- **集成摩擦**：对于 Node.js 编写的前端 CLI 工具而言，跨平台无缝调用外部脚本（如 Python、Java）简直是一场噩梦。

**Graphify-Go 彻底解决了这些痛点：**
- **零环境依赖**：通过 CGO 静态链接编译为单一的二进制文件（`.exe`, ELF, Mach-O）。用户电脑上**不需要**安装 Python、Node、JVM 或 C 编译器。
- **极速并发**：底层采用 `tree-sitter`，结合 Go 语言标志性的 Goroutine 并发池，AST 提取与图谱构建速度达到毫秒级。
- **开箱即用**：直接输出标准的 `graph.json` 架构数据、可交互的 `graph.html` 可视化图谱和 Markdown 摘要，并内置了 Louvain 社区发现算法，前端工具可通过 npm 一键集成。

### ✨ 核心特性
- **高并发 AST 提取**：扫描工作区，利用 Tree-sitter S-expressions 精准捕获 `类`、`函数`、`方法`、`调用` 与 `导入`。
- **图论引擎构建**：自动合并相同实体的节点定义与调用，并根据代码调用频次自动累加边权重（Weight）。
- **社区发现算法**：内置 Louvain 模块度聚类算法，自动将高内聚的代码划分到同一个业务模块（Community）。
- **架构洞察分析**：智能过滤“合成占位符”，提纯计算真正的 `上帝节点 (God Nodes)`，并能敏锐地发现跨越不同模块的 `意外连接 (Surprising Connections)`。

### 🛠️ 支持的语言 (Tree-sitter)
- [x] JavaScript / TypeScript
- [x] Python
- [x] Go
- [ ] *即将支持原版全量生态 (Java, C, C++, Ruby, C#, Kotlin, Scala, PHP, Swift)*

### 🚀 使用方式

你可以直接从 [Releases](https://github.com/ian000/graphify-go/releases) 页面下载编译好的二进制文件，或者通过 npm 包装器一键安装：

```bash
npm install graphify-go
```

直接运行 CLI 命令：
```bash
# 扫描分析当前目录，在终端输出架构摘要
graphify-go

# 扫描指定目录，并将 JSON 和 Markdown 报告导出到目标文件夹
graphify-go -dir ./my-project -out ./reports
```

### 💻 本地开发
```bash
git clone https://github.com/ian000/graphify-go.git
cd graphify-go
# 注意：Windows 需要安装 MSYS2/GCC 环境以支持 cgo 编译
go build ./cmd/graphify
```

## 📄 License
MIT License