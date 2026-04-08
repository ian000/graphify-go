package parser

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

// ExtractedData 保存从单个文件中提取出的节点和边信息
type ExtractedData struct {
	FilePath string
	Classes  []string
	Funcs    []string
	Methods  []string
	Calls    []string
	Imports  []string
}

// Extractor 负责执行 AST 查询并封装数据
type Extractor struct {
	registry *Registry
}

func NewExtractor(reg *Registry) *Extractor {
	return &Extractor{registry: reg}
}

// ParseAndExtract 接收文件路径和代码字节，返回提取到的结构化数据
func (e *Extractor) ParseAndExtract(filePath string, sourceCode []byte) (*ExtractedData, error) {
	// 1. 探测语言类型
	lang := e.registry.DetectLanguage(filePath)
	if lang == LangUnknown {
		return nil, fmt.Errorf("unsupported file extension for: %s", filePath)
	}

	// 2. 获取解析器
	parser, err := e.registry.GetParser(lang)
	if err != nil {
		return nil, err
	}

	// 3. 生成 AST
	tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", filePath, err)
	}
	defer tree.Close()

	// 4. 准备 Query
	var queryStr string
	switch lang {
	case LangJavaScript, LangTypeScript:
		queryStr = JSQuery
	case LangPython:
		queryStr = PYQuery
	case LangGo:
		queryStr = GOQuery
	default:
		return nil, fmt.Errorf("queries for language %s are not implemented yet", lang)
	}

	q, err := sitter.NewQuery([]byte(queryStr), e.registry.parsers[lang])
	if err != nil {
		return nil, fmt.Errorf("invalid query for %s: %v", lang, err)
	}
	defer q.Close()

	// 5. 执行查询
	qc := sitter.NewQueryCursor()
	defer qc.Close()
	qc.Exec(q, tree.RootNode())

	// 6. 收集数据
	data := &ExtractedData{
		FilePath: filePath,
	}

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			captureName := q.CaptureNameForId(c.Index)
			nodeText := c.Node.Content(sourceCode)

			switch captureName {
			case "class.name":
				data.Classes = append(data.Classes, nodeText)
			case "function.name":
				data.Funcs = append(data.Funcs, nodeText)
			case "method.name":
				data.Methods = append(data.Methods, nodeText)
			case "call.function", "call.method":
				data.Calls = append(data.Calls, nodeText)
			case "import.source":
				data.Imports = append(data.Imports, nodeText)
			}
		}
	}

	return data, nil
}
