package parser

import (
	"os"
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

// Scanner 负责遍历目录，并过滤掉被忽略的文件
type Scanner struct {
	ignoreMatcher *ignore.GitIgnore
	registry      *Registry
}

// NewScanner 初始化一个遍历器。它会尝试在工作区根目录读取 .aiignore 或 .gitignore
func NewScanner(workspace string, reg *Registry) *Scanner {
	var matcher *ignore.GitIgnore

	// 优先读取 .aiignore，如果没有则读取 .gitignore
	ignorePath := filepath.Join(workspace, ".aiignore")
	if _, err := os.Stat(ignorePath); os.IsNotExist(err) {
		ignorePath = filepath.Join(workspace, ".gitignore")
	}

	if _, err := os.Stat(ignorePath); err == nil {
		matcher, _ = ignore.CompileIgnoreFile(ignorePath)
	} else {
		// 如果没有忽略文件，给一个默认的忽略规则
		matcher = ignore.CompileIgnoreLines("node_modules", ".git", "dist", "build", "vendor")
	}

	return &Scanner{
		ignoreMatcher: matcher,
		registry:      reg,
	}
}

// Scan 递归遍历指定目录，返回所有需要被解析的源文件路径
func (s *Scanner) Scan(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对于根目录的相对路径，以便 ignore 引擎进行匹配
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return nil
		}
		// Windows 路径转为 Unix 风格的斜杠，因为 ignore 库通常用 /
		relPath = filepath.ToSlash(relPath)

		// 跳过被忽略的文件或目录
		if s.ignoreMatcher != nil && s.ignoreMatcher.MatchesPath(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 我们只收集非目录文件
		if !info.IsDir() {
			// 只收集我们在 Registry 中支持的语言文件
			lang := s.registry.DetectLanguage(path)
			if lang != LangUnknown {
				// 将被纳入提取池的合法文件
				files = append(files, path)
			}
		}
		return nil
	})

	return files, err
}
