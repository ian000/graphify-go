package parser

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
)

// ProcessWorkspace 是我们的并发调度引擎
// 它扫描目录，把所有文件丢进 Worker Pool 并发提取，然后收集所有结果
func ProcessWorkspace(workspace string) ([]*ExtractedData, error) {
	registry := NewRegistry()
	scanner := NewScanner(workspace, registry)
	extractor := NewExtractor(registry)

	// 1. 扫描获取所有待处理的文件路径
	fmt.Println("🔍 Scanning workspace for source files...")
	files, err := scanner.Scan(workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to scan workspace: %v", err)
	}
	fmt.Printf("✅ Found %d source files to process.\n", len(files))

	if len(files) == 0 {
		return []*ExtractedData{}, nil
	}

	// 2. 初始化并发池结构 (Goroutines)
	// 根据系统的 CPU 核心数来决定启动多少个 Worker 榨干性能
	numWorkers := runtime.NumCPU()
	if len(files) < numWorkers {
		numWorkers = len(files) // 文件太少没必要开那么多 worker
	}

	jobs := make(chan string, len(files))
	results := make(chan *ExtractedData, len(files))
	var wg sync.WaitGroup

	// 3. 启动 Worker 们
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for filePath := range jobs {
				// 读取文件内容
				sourceCode, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Printf("⚠️ Worker %d failed to read file %s: %v\n", workerID, filePath, err)
					continue
				}

				// 丢给我们的 AST 提取器
				data, err := extractor.ParseAndExtract(filePath, sourceCode)
				if err != nil {
					fmt.Printf("⚠️ Worker %d failed to extract AST for %s: %v\n", workerID, filePath, err)
					continue
				}

				// 将成功提取的数据丢入结果管道
				results <- data
			}
		}(w)
	}

	// 4. 将所有任务推送到 jobs 管道，然后关闭管道
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// 5. 启动一个后台协程，等待所有 Worker 完成后，关闭结果管道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 6. 收集所有结果
	var allData []*ExtractedData
	for data := range results {
		allData = append(allData, data)
	}

	fmt.Printf("🚀 Successfully extracted ASTs from %d files.\n", len(allData))
	return allData, nil
}
