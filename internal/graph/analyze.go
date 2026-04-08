package graph

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// GodNodeResult 包含上帝节点的统计信息
type GodNodeResult struct {
	ID    string
	Label string
	Edges int
}

// SurprisingConnectionResult 包含意外连接的信息
type SurprisingConnectionResult struct {
	Source      string
	Target      string
	SourceFiles []string
	Relation    string
	Note        string
}

// Analyze 封装了图的各种分析算法
type Analyzer struct {
	graph *Graph
}

func NewAnalyzer(g *Graph) *Analyzer {
	return &Analyzer{graph: g}
}

// isFileNode 判定是否为合成的文件节点或方法占位符
func (a *Analyzer) isFileNode(n *Node) bool {
	if n.Label == "" {
		return false
	}

	// 1. 如果是一个文件级中枢节点（名称正好等于 source_file 的 basename）
	if n.SourceFile != "" {
		if n.Label == filepath.Base(n.SourceFile) {
			return true
		}
	}

	// 2. 如果是一个方法占位符（AST 提取器产生的以 '.' 开头并以 '()' 结尾）
	if strings.HasPrefix(n.Label, ".") && strings.HasSuffix(n.Label, "()") {
		return true
	}

	// 3. 如果是一个孤立的函数节点（度数 <= 1，通常只包含一个 contains 边）
	if strings.HasSuffix(n.Label, "()") {
		if a.Degree(n.ID) <= 1 {
			return true
		}
	}

	return false
}

// isConceptNode 判定是否为手动注入的语义概念节点
func (a *Analyzer) isConceptNode(n *Node) bool {
	// 如果没有 SourceFile，通常认为是外部概念或没有真实落地
	if n.SourceFile == "" {
		return true
	}
	// 如果 SourceFile 没有扩展名，大概率也是个语义概念而非真实文件
	if !strings.Contains(filepath.Base(n.SourceFile), ".") {
		return true
	}
	return false
}

// Degree 返回某个节点相连的边的数量（入度+出度）
func (a *Analyzer) Degree(nodeID string) int {
	deg := 0
	for _, edge := range a.graph.Edges {
		if edge.Source == nodeID || edge.Target == nodeID {
			deg++
		}
	}
	return deg
}

// GodNodes 获取最核心的业务抽象实体，剔除所有的文件/合成占位符
func (a *Analyzer) GodNodes(topN int) []GodNodeResult {
	// 计算所有节点的度数
	degreeMap := make(map[string]int)
	for _, edge := range a.graph.Edges {
		degreeMap[edge.Source]++
		degreeMap[edge.Target]++
	}

	// 过滤合成节点，并转为 Slice
	var candidates []GodNodeResult
	for id, deg := range degreeMap {
		node := a.graph.Nodes[id]
		if node == nil {
			continue
		}
		if a.isFileNode(node) || a.isConceptNode(node) {
			continue
		}
		candidates = append(candidates, GodNodeResult{
			ID:    id,
			Label: node.Label,
			Edges: deg,
		})
	}

	// 降序排列
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Edges > candidates[j].Edges
	})

	if len(candidates) > topN {
		return candidates[:topN]
	}
	return candidates
}

// SurprisingConnections 查找跨文件、跨社区且非结构性的非预期连接
func (a *Analyzer) SurprisingConnections(topN int) []SurprisingConnectionResult {
	var candidates []SurprisingConnectionResult

	// 为了去重跨越社区的边界连接
	seenPairs := make(map[string]bool)

	for _, edge := range a.graph.Edges {
		u := a.graph.Nodes[edge.Source]
		v := a.graph.Nodes[edge.Target]

		// 保护逻辑
		if u == nil || v == nil {
			continue
		}

		// 如果同属于一个社区，则认为关联是显而易见的，不够 surprise
		if u.Community == v.Community {
			continue
		}

		// 跳过显式的结构性关系边
		if edge.Relation == "imports" || edge.Relation == "imports_from" || edge.Relation == "contains" || edge.Relation == "method" {
			continue
		}

		// 剔除上帝节点过滤掉的占位符文件节点
		if a.isFileNode(u) || a.isFileNode(v) {
			continue
		}
		if a.isConceptNode(u) || a.isConceptNode(v) {
			continue
		}

		// 构建唯一跨界对
		var pair string
		if u.Community < v.Community {
			pair = fmt.Sprintf("%d-%d", u.Community, v.Community)
		} else {
			pair = fmt.Sprintf("%d-%d", v.Community, u.Community)
		}

		if seenPairs[pair] {
			continue
		}
		seenPairs[pair] = true

		candidates = append(candidates, SurprisingConnectionResult{
			Source:      u.Label,
			Target:      v.Label,
			SourceFiles: []string{u.SourceFile, v.SourceFile},
			Relation:    edge.Relation,
			Note:        fmt.Sprintf("Bridges Community %d → Community %d", u.Community, v.Community),
		})

		if len(candidates) >= topN {
			break
		}
	}

	return candidates
}
