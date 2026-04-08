package graph

import (
	"path/filepath"
	"strings"

	"github.com/kings2017/graphify-go/internal/parser"
)

// Builder 负责将提取的数据转化为有向图
type Builder struct {
	graph *Graph
}

func NewBuilder() *Builder {
	return &Builder{
		graph: NewGraph(),
	}
}

// Build 从解析器的批量数据构建完整的图模型
func (b *Builder) Build(dataset []*parser.ExtractedData) *Graph {
	for _, data := range dataset {
		b.processFile(data)
	}
	return b.graph
}

// processFile 处理单个文件数据，创建节点和边
func (b *Builder) processFile(data *parser.ExtractedData) {
	// 1. 创建文件级节点（作为容器枢纽）
	fileName := filepath.Base(data.FilePath)
	fileID := GenerateNodeID("file", fileName)
	b.graph.AddNode(fileID, fileName, "file", data.FilePath)

	// 2. 处理 Classes
	for _, cls := range data.Classes {
		classID := GenerateNodeID("entity", cls)
		b.graph.AddNode(classID, cls, "class", data.FilePath)
		b.graph.AddEdge(fileID, classID, "contains", 1.0)
	}

	// 3. 处理 Funcs
	for _, fn := range data.Funcs {
		funcID := GenerateNodeID("entity", fn)
		funcLabel := fn + "()"
		b.graph.AddNode(funcID, funcLabel, "function", data.FilePath)
		b.graph.AddEdge(fileID, funcID, "contains", 1.0)
	}

	// 4. 处理 Methods (合成节点逻辑：以 '.' 开头的占位符)
	for _, method := range data.Methods {
		methodID := GenerateNodeID("entity", method)
		methodLabel := "." + method + "()"
		b.graph.AddNode(methodID, methodLabel, "method", data.FilePath)
		b.graph.AddEdge(fileID, methodID, "contains", 1.0)
	}

	// 5. 处理 Calls (从该文件向外发起的调用)
	for _, call := range data.Calls {
		// 使用相同的前缀 "entity"，这样定义和调用就能通过相同的 ID 自动合并！
		callID := GenerateNodeID("entity", call)

		callLabel := call
		if !strings.HasSuffix(callLabel, "()") {
			callLabel = callLabel + "()"
		}

		b.graph.AddNode(callID, callLabel, "function", "")
		b.graph.AddEdge(fileID, callID, "calls", 1.0)
	}

	// 6. 处理 Imports (跨文件模块引用)
	for _, imp := range data.Imports {
		// 去除引号和分号等冗余字符
		cleanImp := strings.Trim(imp, "\"'`; \t")
		impName := filepath.Base(cleanImp) // 比如从 'path/to/utils' 拿到 'utils'

		if impName == "" || impName == "." {
			continue
		}

		impID := GenerateNodeID("file", impName)
		// 创建一个被引入的文件节点，由于是外部的或还没有被扫描到的，SourceFile 暂时留空。如果后续被扫描到，会有更新（但目前 AddNode 只插入第一次，我们需要确保合并逻辑，或者留给 Graph 自己处理）
		b.graph.AddNode(impID, impName, "file", "")
		b.graph.AddEdge(fileID, impID, "imports", 1.0)
	}
}
