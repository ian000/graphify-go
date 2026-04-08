package graph_test

import (
	"testing"

	"github.com/kings2017/graphify-go/internal/graph"
	"github.com/kings2017/graphify-go/internal/parser"
)

func TestBuilder_BuildAndAnalyze(t *testing.T) {
	// 1. Mock 两份 ExtractedData
	file1 := &parser.ExtractedData{
		FilePath: "src/utils.js",
		Funcs:    []string{"helper"},
	}

	file2 := &parser.ExtractedData{
		FilePath: "src/main.js",
		Classes:  []string{"App"},
		Methods:  []string{"init"},
		Calls:    []string{"helper"},
		Imports:  []string{"./utils"},
	}

	dataset := []*parser.ExtractedData{file1, file2}

	// 2. Build Graph
	builder := graph.NewBuilder()
	g := builder.Build(dataset)

	// 验证节点和边的建立
	if len(g.Nodes) == 0 || len(g.Edges) == 0 {
		t.Fatalf("Graph is empty after build")
	}

	// 检查 "helper" 函数是否被合并 (ID 应该一致)
	helperID := graph.GenerateNodeID("entity", "helper")
	helperNode, exists := g.Nodes[helperID]
	if !exists {
		t.Fatalf("Expected helper node to exist")
	}
	// helper 的 sourceFile 应该是 src/utils.js
	if helperNode.SourceFile != "src/utils.js" {
		t.Errorf("Expected helper node to have source file 'src/utils.js', got %v", helperNode.SourceFile)
	}

	// 检查边: main.js -> helper() (calls)
	mainFileID := graph.GenerateNodeID("file", "main.js")
	edgeKey := mainFileID + "|" + helperID + "|calls"
	if _, ok := g.Edges[edgeKey]; !ok {
		t.Errorf("Expected call edge from main.js to helper()")
	}

	// 3. Test Analyzer 净化逻辑
	analyzer := graph.NewAnalyzer(g)

	// main.js 本身是一个 File Node，应当被剔除出 GodNodes
	godNodes := analyzer.GodNodes(10)
	for _, n := range godNodes {
		if n.ID == mainFileID {
			t.Errorf("File node 'main.js' should be filtered out from God Nodes")
		}
		if n.ID == graph.GenerateNodeID("entity", "init") {
			t.Errorf("Method stub 'init' (which is isolated or stub) might be filtered based on logic")
		}
	}
}
