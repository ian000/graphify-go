package export

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kings2017/graphify-go/internal/graph"
)

// ExportSystemGraphMD 1:1 复刻 Python 版的 system-graph.md 摘要
func ExportSystemGraphMD(g *graph.Graph) string {
	analyzer := graph.NewAnalyzer(g)

	godNodes := analyzer.GodNodes(10)
	surprises := analyzer.SurprisingConnections(5)

	var buf bytes.Buffer

	// 标题和基本统计
	buf.WriteString("# System Graph Summary\n\n")
	buf.WriteString(fmt.Sprintf("**Total Nodes:** %d\n", len(g.Nodes)))
	buf.WriteString(fmt.Sprintf("**Total Edges:** %d\n\n", len(g.Edges)))

	// God Nodes 模块
	buf.WriteString("## 🏛️ God Nodes (Top Entities)\n")
	buf.WriteString("> The most connected real abstractions (excluding files and synthetic stubs).\n\n")
	if len(godNodes) == 0 {
		buf.WriteString("*No significant God Nodes found.*\n\n")
	} else {
		for _, n := range godNodes {
			buf.WriteString(fmt.Sprintf("- **%s** (`%s`): %d connections\n", n.Label, n.ID, n.Edges))
		}
		buf.WriteString("\n")
	}

	// Surprising Connections 模块
	buf.WriteString("## ⚡ Surprising Connections\n")
	buf.WriteString("> Non-obvious cross-community or cross-file connections that bridge the architecture.\n\n")
	if len(surprises) == 0 {
		buf.WriteString("*No surprising connections found.*\n\n")
	} else {
		for _, s := range surprises {
			src := s.SourceFiles[0]
			if src == "" {
				src = "unknown"
			}
			tgt := s.SourceFiles[1]
			if tgt == "" {
				tgt = "unknown"
			}
			buf.WriteString(fmt.Sprintf("- **%s** (`%s`) ↔ **%s** (`%s`)\n", s.Source, src, s.Target, tgt))
			buf.WriteString(fmt.Sprintf("  - *Relation:* `%s`\n", s.Relation))
			buf.WriteString(fmt.Sprintf("  - *Why:* %s\n", s.Note))
		}
		buf.WriteString("\n")
	}

	// Communities (模块化划分)
	buf.WriteString("## 🧩 Architecture Communities (Modules)\n")
	buf.WriteString("> Code organized by graph clustering (Louvain algorithm).\n\n")

	commMap := make(map[int][]string)
	for _, n := range g.Nodes {
		// 只统计真实节点
		if n.Label != "" {
			commMap[n.Community] = append(commMap[n.Community], n.Label)
		}
	}

	if len(commMap) == 0 {
		buf.WriteString("*No communities detected.*\n\n")
	} else {
		for cid, nodes := range commMap {
			// 如果一个社区太大，只展示前 5 个代表性节点
			displayCount := len(nodes)
			if displayCount > 5 {
				displayCount = 5
			}
			sample := nodes[:displayCount]
			buf.WriteString(fmt.Sprintf("- **Community %d** (%d nodes): %s", cid, len(nodes), strings.Join(sample, ", ")))
			if len(nodes) > 5 {
				buf.WriteString(", ...\n")
			} else {
				buf.WriteString("\n")
			}
		}
	}

	return buf.String()
}
