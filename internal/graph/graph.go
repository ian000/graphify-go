package graph

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Node 代表图中的一个节点
type Node struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Type       string `json:"type"` // e.g. "file", "class", "function", "method"
	Community  int    `json:"community"`
	SourceFile string `json:"source_file"`
}

// Edge 代表图中的一条有向边
type Edge struct {
	Source   string  `json:"source"`
	Target   string  `json:"target"`
	Weight   float64 `json:"weight"`
	Relation string  `json:"type"` // e.g. "contains", "calls", "imports"
}

// Graph 代表完整的依赖图谱
type Graph struct {
	Nodes map[string]*Node
	Edges map[string]*Edge
}

// NewGraph 初始化一个空的图
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string]*Edge),
	}
}

// AddNode 幂等地添加节点，如果节点已存在，尝试用更丰富的信息覆盖
func (g *Graph) AddNode(id, label, nodeType, sourceFile string) {
	if existing, exists := g.Nodes[id]; !exists {
		g.Nodes[id] = &Node{
			ID:         id,
			Label:      label,
			Type:       nodeType,
			SourceFile: sourceFile,
			Community:  0, // 默认 0，后续由 Leiden/Louvain 算法更新
		}
	} else {
		// Python 版逻辑：多次添加同一个 ID 时，后添加的更详细属性会覆盖先前的占位符
		if existing.SourceFile == "" && sourceFile != "" {
			existing.SourceFile = sourceFile
		}
		if existing.Type == "" && nodeType != "" {
			existing.Type = nodeType
		}
		// 尽量不覆盖已有 label，除非它本来为空
		if existing.Label == "" && label != "" {
			existing.Label = label
		}
	}
}

// AddEdge 幂等地添加边，如果存在同类型边，则累加权重
func (g *Graph) AddEdge(sourceID, targetID, relation string, weight float64) {
	edgeKey := fmt.Sprintf("%s|%s|%s", sourceID, targetID, relation)
	if edge, exists := g.Edges[edgeKey]; exists {
		edge.Weight += weight
	} else {
		g.Edges[edgeKey] = &Edge{
			Source:   sourceID,
			Target:   targetID,
			Weight:   weight,
			Relation: relation,
		}
	}
}

// ToJSON 将图序列化为符合原版 Python graphify 输出格式的 JSON
func (g *Graph) ToJSON() ([]byte, error) {
	// 组装最终格式
	type outputFormat struct {
		Nodes []*Node `json:"nodes"`
		Links []*Edge `json:"links"`
	}

	out := &outputFormat{
		Nodes: make([]*Node, 0, len(g.Nodes)),
		Links: make([]*Edge, 0, len(g.Edges)),
	}

	for _, n := range g.Nodes {
		out.Nodes = append(out.Nodes, n)
	}
	for _, e := range g.Edges {
		out.Links = append(out.Links, e)
	}

	return json.MarshalIndent(out, "", "  ")
}

var nonAlnumRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// GenerateNodeID 生成稳定唯一的节点 ID，1:1 复刻 Python 的 _make_id
func GenerateNodeID(parts ...string) string {
	var validParts []string
	for _, p := range parts {
		// 去除两端 "_" 和 "."
		cleanP := strings.Trim(p, "_.")
		if cleanP != "" {
			validParts = append(validParts, cleanP)
		}
	}
	combined := strings.Join(validParts, "_")
	cleaned := nonAlnumRegex.ReplaceAllString(combined, "_")
	return strings.ToLower(strings.Trim(cleaned, "_"))
}
