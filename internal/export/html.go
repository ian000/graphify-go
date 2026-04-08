package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kings2017/graphify-go/internal/graph"
)

var communityColors = []string{
	"#4E79A7", "#F28E2B", "#E15759", "#76B7B2", "#59A14F",
	"#EDC948", "#B07AA1", "#FF9DA7", "#9C755F", "#BAB0AC",
}

const maxNodesForViz = 5000

func htmlStyles() string {
	return `<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { background: #0f0f1a; color: #e0e0e0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; display: flex; height: 100vh; overflow: hidden; }
  #graph { flex: 1; }
  #sidebar { width: 280px; background: #1a1a2e; border-left: 1px solid #2a2a4e; display: flex; flex-direction: column; overflow: hidden; }
  #search-wrap { padding: 12px; border-bottom: 1px solid #2a2a4e; }
  #search { width: 100%; background: #0f0f1a; border: 1px solid #3a3a5e; color: #e0e0e0; padding: 7px 10px; border-radius: 6px; font-size: 13px; outline: none; }
  #search:focus { border-color: #4E79A7; }
  #search-results { max-height: 140px; overflow-y: auto; padding: 4px 12px; border-bottom: 1px solid #2a2a4e; display: none; }
  .search-item { padding: 4px 6px; cursor: pointer; border-radius: 4px; font-size: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .search-item:hover { background: #2a2a4e; }
  #info-panel { padding: 14px; border-bottom: 1px solid #2a2a4e; min-height: 140px; }
  #info-panel h3 { font-size: 13px; color: #aaa; margin-bottom: 8px; text-transform: uppercase; letter-spacing: 0.05em; }
  #info-content { font-size: 13px; color: #ccc; line-height: 1.6; }
  #info-content .field { margin-bottom: 5px; }
  #info-content .field b { color: #e0e0e0; }
  #info-content .empty { color: #555; font-style: italic; }
  .neighbor-link { display: block; padding: 2px 6px; margin: 2px 0; border-radius: 3px; cursor: pointer; font-size: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; border-left: 3px solid #333; }
  .neighbor-link:hover { background: #2a2a4e; }
  #neighbors-list { max-height: 160px; overflow-y: auto; margin-top: 4px; }
  #legend-wrap { flex: 1; overflow-y: auto; padding: 12px; }
  #legend-wrap h3 { font-size: 13px; color: #aaa; margin-bottom: 10px; text-transform: uppercase; letter-spacing: 0.05em; }
  .legend-item { display: flex; align-items: center; gap: 8px; padding: 4px 0; cursor: pointer; border-radius: 4px; font-size: 12px; }
  .legend-item:hover { background: #2a2a4e; padding-left: 4px; }
  .legend-item.dimmed { opacity: 0.35; }
  .legend-dot { width: 12px; height: 12px; border-radius: 50%; flex-shrink: 0; }
  .legend-label { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .legend-count { color: #666; font-size: 11px; }
  #stats { padding: 10px 14px; border-top: 1px solid #2a2a4e; font-size: 11px; color: #555; }
</style>`
}

func htmlScript(nodesJSON, edgesJSON, legendJSON string) string {
	return fmt.Sprintf(`<script>
const RAW_NODES = %s;
const RAW_EDGES = %s;
const LEGEND = %s;

// Build vis datasets
const nodesDS = new vis.DataSet(RAW_NODES.map(n => ({
  id: n.id, label: n.label, color: n.color, size: n.size,
  font: n.font, title: n.title,
  _community: n.community, _community_name: n.community_name,
  _source_file: n.source_file, _file_type: n.file_type, _degree: n.degree,
})));

const edgesDS = new vis.DataSet(RAW_EDGES.map((e, i) => ({
  id: i, from: e.from, to: e.to,
  label: '',
  title: e.title,
  dashes: e.dashes,
  width: e.width,
  color: e.color,
  arrows: { to: { enabled: true, scaleFactor: 0.5 } },
})));

const container = document.getElementById('graph');
const network = new vis.Network(container, { nodes: nodesDS, edges: edgesDS }, {
  physics: {
    enabled: true,
    solver: 'forceAtlas2Based',
    forceAtlas2Based: {
      gravitationalConstant: -60,
      centralGravity: 0.005,
      springLength: 120,
      springConstant: 0.08,
      damping: 0.4,
      avoidOverlap: 0.8,
    },
    stabilization: { iterations: 200, fit: true },
  },
  interaction: {
    hover: true,
    tooltipDelay: 100,
    hideEdgesOnDrag: true,
    navigationButtons: false,
    keyboard: false,
  },
  nodes: { shape: 'dot', borderWidth: 1.5 },
  edges: { smooth: { type: 'continuous', roundness: 0.2 }, selectionWidth: 3 },
});

network.once('stabilizationIterationsDone', () => {
  network.setOptions({ physics: { enabled: false } });
});

function showInfo(nodeId) {
  const n = nodesDS.get(nodeId);
  if (!n) return;
  const neighborIds = network.getConnectedNodes(nodeId);
  const neighborItems = neighborIds.map(nid => {
    const nb = nodesDS.get(nid);
    const color = nb ? nb.color.background : '#555';
    return '<span class="neighbor-link" style="border-left-color:' + color + '" onclick="focusNode(\'' + nid + '\')">' + (nb ? nb.label : nid) + '</span>';
  }).join('');
  document.getElementById('info-content').innerHTML = 
    '<div class="field"><b>' + n.label + '</b></div>' +
    '<div class="field">Type: ' + (n._file_type || 'unknown') + '</div>' +
    '<div class="field">Community: ' + n._community_name + '</div>' +
    '<div class="field">Source: ' + (n._source_file || '-') + '</div>' +
    '<div class="field">Degree: ' + n._degree + '</div>' +
    (neighborIds.length ? '<div class="field" style="margin-top:8px;color:#aaa;font-size:11px">Neighbors (' + neighborIds.length + ')</div><div id="neighbors-list">' + neighborItems + '</div>' : '');
}

function focusNode(nodeId) {
  network.focus(nodeId, { scale: 1.4, animation: true });
  network.selectNodes([nodeId]);
  showInfo(nodeId);
}

network.on('click', params => {
  if (params.nodes.length > 0) showInfo(params.nodes[0]);
  else document.getElementById('info-content').innerHTML = '<span class="empty">Click a node to inspect it</span>';
});

const searchInput = document.getElementById('search');
const searchResults = document.getElementById('search-results');
searchInput.addEventListener('input', () => {
  const q = searchInput.value.toLowerCase().trim();
  searchResults.innerHTML = '';
  if (!q) { searchResults.style.display = 'none'; return; }
  const matches = RAW_NODES.filter(n => n.label.toLowerCase().includes(q)).slice(0, 20);
  if (!matches.length) { searchResults.style.display = 'none'; return; }
  searchResults.style.display = 'block';
  matches.forEach(n => {
    const el = document.createElement('div');
    el.className = 'search-item';
    el.textContent = n.label;
    el.style.borderLeft = '3px solid ' + n.color.background;
    el.style.paddingLeft = '8px';
    el.onclick = () => {
      network.focus(n.id, { scale: 1.5, animation: true });
      network.selectNodes([n.id]);
      showInfo(n.id);
      searchResults.style.display = 'none';
      searchInput.value = '';
    };
    searchResults.appendChild(el);
  });
});
document.addEventListener('click', e => {
  if (!searchResults.contains(e.target) && e.target !== searchInput)
    searchResults.style.display = 'none';
});

const hiddenCommunities = new Set();
const legendEl = document.getElementById('legend');
LEGEND.forEach(c => {
  const item = document.createElement('div');
  item.className = 'legend-item';
  item.innerHTML = '<div class="legend-dot" style="background:' + c.color + '"></div>' +
    '<span class="legend-label">' + c.label + '</span>' +
    '<span class="legend-count">' + c.count + '</span>';
  item.onclick = () => {
    if (hiddenCommunities.has(c.cid)) {
      hiddenCommunities.delete(c.cid);
      item.classList.remove('dimmed');
    } else {
      hiddenCommunities.add(c.cid);
      item.classList.add('dimmed');
    }
    const updates = RAW_NODES
      .filter(n => n.community === c.cid)
      .map(n => ({ id: n.id, hidden: hiddenCommunities.has(c.cid) }));
    nodesDS.update(updates);
  };
  legendEl.appendChild(item);
});
</script>`, nodesJSON, edgesJSON, legendJSON)
}

// ExportSystemGraphHTML 1:1 复刻 Python 版的 HTML 可视化图谱
func ExportSystemGraphHTML(g *graph.Graph, outputPath string) error {
	if len(g.Nodes) > maxNodesForViz {
		return fmt.Errorf("graph has %d nodes - too large for HTML viz. Max allowed is %d", len(g.Nodes), maxNodesForViz)
	}

	// 1. 统计 Degree
	degrees := make(map[string]int)
	maxDeg := 1
	for _, e := range g.Edges {
		degrees[e.Source]++
		degrees[e.Target]++
	}
	for _, deg := range degrees {
		if deg > maxDeg {
			maxDeg = deg
		}
	}

	// 2. 统计 Community 数量和大小
	communityCounts := make(map[int]int)
	for _, n := range g.Nodes {
		communityCounts[n.Community]++
	}

	// 3. 构建 Nodes 数据
	type visNodeColor struct {
		Background string `json:"background"`
		Border     string `json:"border"`
		Highlight  struct {
			Background string `json:"background"`
			Border     string `json:"border"`
		} `json:"highlight"`
	}
	type visNodeFont struct {
		Size  int    `json:"size"`
		Color string `json:"color"`
	}
	type visNode struct {
		ID            string       `json:"id"`
		Label         string       `json:"label"`
		Color         visNodeColor `json:"color"`
		Size          float64      `json:"size"`
		Font          visNodeFont  `json:"font"`
		Title         string       `json:"title"`
		Community     int          `json:"community"`
		CommunityName string       `json:"community_name"`
		SourceFile    string       `json:"source_file"`
		FileType      string       `json:"file_type"`
		Degree        int          `json:"degree"`
	}

	var visNodes []visNode
	for _, n := range g.Nodes {
		cid := n.Community
		colorHex := communityColors[cid%len(communityColors)]
		deg := degrees[n.ID]
		if deg == 0 {
			deg = 1
		}
		size := 10.0 + 30.0*(float64(deg)/float64(maxDeg))

		fontSize := 0
		if float64(deg) >= float64(maxDeg)*0.15 {
			fontSize = 12
		}

		commName := fmt.Sprintf("Community %d", cid)

		visNodes = append(visNodes, visNode{
			ID:    n.ID,
			Label: n.Label,
			Color: visNodeColor{
				Background: colorHex,
				Border:     colorHex,
				Highlight: struct {
					Background string `json:"background"`
					Border     string `json:"border"`
				}{Background: "#ffffff", Border: colorHex},
			},
			Size:          size,
			Font:          visNodeFont{Size: fontSize, Color: "#ffffff"},
			Title:         n.Label,
			Community:     cid,
			CommunityName: commName,
			SourceFile:    n.SourceFile,
			FileType:      n.Type,
			Degree:        deg,
		})
	}

	// 4. 构建 Edges 数据
	type visEdgeColor struct {
		Opacity float64 `json:"opacity"`
	}
	type visEdge struct {
		From       string       `json:"from"`
		To         string       `json:"to"`
		Label      string       `json:"label"`
		Title      string       `json:"title"`
		Dashes     bool         `json:"dashes"`
		Width      int          `json:"width"`
		Color      visEdgeColor `json:"color"`
		Confidence string       `json:"confidence"`
	}

	var visEdges []visEdge
	for _, e := range g.Edges {
		confidence := "EXTRACTED" // 目前 Go 版本全是静态提取的，所以都是 EXTRACTED
		dashes := false
		width := 2
		opacity := 0.7

		visEdges = append(visEdges, visEdge{
			From:       e.Source,
			To:         e.Target,
			Label:      e.Relation,
			Title:      fmt.Sprintf("%s [%s]", e.Relation, confidence),
			Dashes:     dashes,
			Width:      width,
			Color:      visEdgeColor{Opacity: opacity},
			Confidence: confidence,
		})
	}

	// 5. 构建 Legend 数据
	type visLegend struct {
		CID   int    `json:"cid"`
		Color string `json:"color"`
		Label string `json:"label"`
		Count int    `json:"count"`
	}
	var visLegends []visLegend
	for cid, count := range communityCounts {
		visLegends = append(visLegends, visLegend{
			CID:   cid,
			Color: communityColors[cid%len(communityColors)],
			Label: fmt.Sprintf("Community %d", cid),
			Count: count,
		})
	}

	nodesJSON, _ := json.Marshal(visNodes)
	edgesJSON, _ := json.Marshal(visEdges)
	legendJSON, _ := json.Marshal(visLegends)

	title := filepath.Base(outputPath)
	stats := fmt.Sprintf("%d nodes &middot; %d edges &middot; %d communities", len(g.Nodes), len(g.Edges), len(communityCounts))

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>graphify - %s</title>
<script src="https://unpkg.com/vis-network/standalone/umd/vis-network.min.js"></script>
%s
</head>
<body>
<div id="graph"></div>
<div id="sidebar">
  <div id="search-wrap">
    <input id="search" type="text" placeholder="Search nodes..." autocomplete="off">
    <div id="search-results"></div>
  </div>
  <div id="info-panel">
    <h3>Node Info</h3>
    <div id="info-content"><span class="empty">Click a node to inspect it</span></div>
  </div>
  <div id="legend-wrap">
    <h3>Communities</h3>
    <div id="legend"></div>
  </div>
  <div id="stats">%s</div>
</div>
%s
</body>
</html>`, title, htmlStyles(), stats, htmlScript(string(nodesJSON), string(edgesJSON), string(legendJSON)))

	return os.WriteFile(outputPath, []byte(html), 0644)
}
