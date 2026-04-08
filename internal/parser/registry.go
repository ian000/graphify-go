package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// LanguageType 定义了支持的编程语言
type LanguageType string

const (
	LangJavaScript LanguageType = "javascript"
	LangTypeScript LanguageType = "typescript"
	LangPython     LanguageType = "python"
	LangGo         LanguageType = "go"
	LangUnknown    LanguageType = "unknown"
)

// Registry 负责根据文件后缀管理和分配 Tree-sitter 解析器
type Registry struct {
	parsers map[LanguageType]*sitter.Language
}

// NewRegistry 初始化并返回一个包含所有受支持语言的注册表
func NewRegistry() *Registry {
	return &Registry{
		parsers: map[LanguageType]*sitter.Language{
			LangJavaScript: javascript.GetLanguage(),
			LangTypeScript: typescript.GetLanguage(),
			LangPython:     python.GetLanguage(),
			LangGo:         golang.GetLanguage(),
		},
	}
}

// DetectLanguage 根据文件名后缀探测语言类型
func (r *Registry) DetectLanguage(filename string) LanguageType {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".js", ".jsx":
		return LangJavaScript
	case ".ts", ".tsx":
		return LangTypeScript
	case ".py":
		return LangPython
	case ".go":
		return LangGo
	default:
		return LangUnknown
	}
}

// GetParser 返回一个配置好对应语言的 Tree-sitter Parser
func (r *Registry) GetParser(lang LanguageType) (*sitter.Parser, error) {
	langPtr, exists := r.parsers[lang]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(langPtr)
	return parser, nil
}
