package cluster

import (
	randv2 "math/rand/v2"
	"time"

	"github.com/kings2017/graphify-go/internal/graph"
	"gonum.org/v1/gonum/graph/community"
	"gonum.org/v1/gonum/graph/simple"
)

// DetectCommunities 使用 Louvain 算法为图谱节点打上 Community 标签
func DetectCommunities(g *graph.Graph) {
	if len(g.Nodes) == 0 {
		return
	}

	// 1. 构建 Gonum 所需的无向权重图
	ug := simple.NewWeightedUndirectedGraph(0, 0)

	// 为了将我们的字符串 ID 映射到 int64，我们需要一个双向映射表
	idToInt := make(map[string]int64)
	intToID := make(map[int64]string)

	var i int64 = 1
	for id := range g.Nodes {
		idToInt[id] = i
		intToID[i] = id
		ug.AddNode(simple.Node(i))
		i++
	}

	// 添加边（如果有重复的无向边，我们需要累加）
	for _, e := range g.Edges {
		u := idToInt[e.Source]
		v := idToInt[e.Target]

		// 忽略自环
		if u == v {
			continue
		}

		if ug.HasEdgeBetween(u, v) {
			existing := ug.WeightedEdge(u, v)
			// 累加权重
			newWeight := existing.Weight() + e.Weight
			ug.SetWeightedEdge(simple.WeightedEdge{
				F: simple.Node(u),
				T: simple.Node(v),
				W: newWeight,
			})
		} else {
			ug.SetWeightedEdge(simple.WeightedEdge{
				F: simple.Node(u),
				T: simple.Node(v),
				W: e.Weight,
			})
		}
	}

	// 2. 执行 Louvain 社区发现
	// resolution 通常设为 1.0 (标准 Modularity)
	src := randv2.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()))
	reduced := community.Modularize(ug, 1.0, src)

	// 3. 将社区结果写回原图
	communities := reduced.Communities()
	for cid, commNodes := range communities {
		for _, n := range commNodes {
			strID := intToID[n.ID()]
			if node, exists := g.Nodes[strID]; exists {
				node.Community = cid
			}
		}
	}
}
