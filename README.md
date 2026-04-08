# Graphify-Go 🚀

*[Read in English](./README.md) | [中文阅读](./README_zh.md)*

**The blazing-fast, zero-dependency AST extraction and code architecture graph engine.**

Designed for code analysis, refactoring assessment, and AI-assisted engineering scaffolding, providing ultimate parsing performance and a seamless cross-platform distribution experience.

### 💡 Why Build Graphify-Go?

In the field of code architecture analysis, traditional tools often face several critical flaws:
- **Dependency Hell**: Many Python-based parsing tools require users to globally install Python 3.9+, pip, and deal with complex compilation errors from libraries like `networkx` and `leidenalg` (which rely on C++).
- **Performance Bottlenecks**: When facing large-scale projects with hundreds of thousands of lines of code, traditional single-threaded scripts are limited by memory overhead, and AST extraction often takes tens of seconds or even minutes.
- **Integration Friction**: For frontend CLI tools written in Node.js, calling external scripts (like Python or Java) cross-platform is a nightmare.

**Graphify-Go solves all of this:**
- **Zero Dependencies**: Compiled down to a single static binary (`.exe`, ELF, Mach-O). No Python, Node, JVM, or C-compilers required on the user's machine.
- **Blazing Fast**: Powered by `tree-sitter` (via `go-tree-sitter`) and Go's native concurrency (goroutines). Extracts and builds the graph in milliseconds.
- **Out of the Box**: Directly outputs standard `graph.json` architecture data, interactive `graph.html` visualization, and Markdown summaries, with built-in Louvain community detection. Frontend tools can integrate it with a single npm command.

### ✨ Features
- **Concurrent AST Extraction**: Scans workspaces and extracts `classes`, `functions`, `methods`, `calls`, and `imports` using Tree-sitter S-expressions.
- **Graph Construction**: Automatically merges identical entities and accumulates edge weights based on call frequencies.
- **Community Detection**: Built-in Louvain modularity algorithm to group highly cohesive code modules together.
- **Architecture Insights**: Purifies "God Nodes" and discovers "Surprising Connections" across different modules.

### 🛠️ Supported Languages (Tree-sitter)
- [x] JavaScript / TypeScript
- [x] Python
- [x] Go
- [ ] *Coming soon (Java, C, C++, Ruby, C#, Kotlin, Scala, PHP, Swift)*

### 🚀 Usage

You can download the pre-compiled binaries from the [Releases](https://github.com/ian000/graphify-go/releases) page, or install via npm wrapper:

```bash
npm install graphify-go
```

Run the CLI directly:
```bash
# Analyze the current directory
graphify-go

# Analyze a specific directory and save reports
graphify-go -dir ./my-project -out ./reports
```

### 💻 Development
```bash
git clone https://github.com/ian000/graphify-go.git
cd graphify-go
# Requires GCC for cgo compilation of tree-sitter
go build ./cmd/graphify
```

## 📄 License
MIT License