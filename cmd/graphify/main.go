package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kings2017/graphify-go/internal/cluster"
	"github.com/kings2017/graphify-go/internal/export"
	"github.com/kings2017/graphify-go/internal/graph"
	"github.com/kings2017/graphify-go/internal/parser"
)

func main() {
	// 1. 定义命令行参数
	dirPtr := flag.String("dir", ".", "Directory to analyze (default: current directory)")
	outPtr := flag.String("out", "", "Output directory for JSON and Markdown reports")
	flag.Parse()

	workspace, err := filepath.Abs(*dirPtr)
	if err != nil {
		log.Fatalf("Invalid directory path: %v\n", err)
	}

	fmt.Println("🚀 Graphify-Go Scanner Starting 🚀")
	fmt.Printf("📂 Target Directory: %s\n", workspace)
	fmt.Println("-----------------------------------------")

	start := time.Now()

	// 2. 触发并发工作池进行文件提取
	results, err := parser.ProcessWorkspace(workspace)
	if err != nil {
		log.Fatalf("❌ Failed to process workspace: %v\n", err)
	}

	if len(results) == 0 {
		fmt.Println("⚠️ No supported source files found in the directory.")
		return
	}

	// 3. 构建依赖图谱
	fmt.Println("🏗️ Building Graph...")
	builder := graph.NewBuilder()
	g := builder.Build(results)
	fmt.Printf("✅ Graph built: %d Nodes, %d Edges\n", len(g.Nodes), len(g.Edges))

	// 4. 运行社区发现 (Louvain)
	fmt.Println("🧠 Running Community Detection (Louvain)...")
	cluster.DetectCommunities(g)

	// 统计社区分布
	commCount := make(map[int]int)
	for _, n := range g.Nodes {
		commCount[n.Community]++
	}
	fmt.Printf("✅ Discovered %d Communities.\n", len(commCount))

	// 5. 输出分析与摘要 (Markdown & JSON)
	mdSummary := export.ExportSystemGraphMD(g)
	jsonBytes, err := g.ToJSON()
	if err != nil {
		log.Fatalf("❌ Failed to serialize graph to JSON: %v\n", err)
	}

	// 如果指定了输出目录，则写入文件；否则仅在终端打印摘要
	if *outPtr != "" {
		outDir, _ := filepath.Abs(*outPtr)
		os.MkdirAll(outDir, 0755)

		jsonPath := filepath.Join(outDir, "graph.json")
		err = os.WriteFile(jsonPath, jsonBytes, 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write graph.json: %v\n", err)
		}

		// 导出 HTML
		htmlPath := filepath.Join(outDir, "graph.html")
		if err := export.ExportSystemGraphHTML(g, htmlPath); err != nil {
			fmt.Printf("⚠️ Skipped HTML export: %v\n", err)
		}

		mdPath := filepath.Join(outDir, "system-graph.md")
		err = os.WriteFile(mdPath, []byte(mdSummary), 0644)
		if err != nil {
			log.Fatalf("❌ Failed to write system-graph.md: %v\n", err)
		}

		fmt.Println("\n-----------------------------------------")
		fmt.Printf("💾 Reports saved to: %s\n", outDir)
		fmt.Println("  📄 graph.json")
		fmt.Println("  🌐 graph.html")
		fmt.Println("  📝 system-graph.md")
	} else {
		fmt.Println("\n-----------------------------------------")
		fmt.Println("📊 Analysis Summary")
		fmt.Println("-----------------------------------------")
		fmt.Println(mdSummary)
	}

	elapsed := time.Since(start)
	fmt.Println("-----------------------------------------")
	fmt.Printf("✨ Done in %v\n", elapsed)
}
