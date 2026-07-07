package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ian000/graphify-go/internal/cluster"
	"github.com/ian000/graphify-go/internal/export"
	"github.com/ian000/graphify-go/internal/graph"
	"github.com/ian000/graphify-go/internal/parser"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		switch args[0] {
		case "analyze":
			return runAnalyze(args[1:])
		case "print":
			return runPrint(args[1:])
		}
	}

	return runLegacy(args)
}

func runLegacy(args []string) error {
	fs := flag.NewFlagSet("graphify-go", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	dirPtr := fs.String("dir", ".", "Directory to analyze (default: current directory)")
	outPtr := fs.String("out", "", "Output directory for JSON and Markdown reports")
	if err := fs.Parse(args); err != nil {
		return err
	}

	return analyzeAndRender(*dirPtr, *outPtr, os.Stdout, os.Stderr)
}

func runAnalyze(args []string) error {
	fs := flag.NewFlagSet("graphify-go analyze", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	dirPtr := fs.String("dir", ".", "Directory to analyze (default: current directory)")
	outPtr := fs.String("out", "", "Output directory for JSON, HTML, and Markdown reports")
	if err := fs.Parse(args); err != nil {
		return err
	}

	return analyzeAndRender(*dirPtr, *outPtr, os.Stdout, os.Stderr)
}

func runPrint(args []string) error {
	fs := flag.NewFlagSet("graphify-go print", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	dirPtr := fs.String("dir", ".", "Directory to analyze (default: current directory)")
	formatPtr := fs.String("format", "markdown", "Output format: json or markdown")
	if err := fs.Parse(args); err != nil {
		return err
	}

	_, _, mdSummary, jsonBytes, err := analyzeWorkspace(*dirPtr, nil)
	if err != nil {
		return err
	}

	switch strings.ToLower(*formatPtr) {
	case "json":
		fmt.Fprint(os.Stdout, string(jsonBytes))
	case "markdown", "md":
		fmt.Fprint(os.Stdout, mdSummary)
	default:
		return fmt.Errorf("unsupported format %q, expected json or markdown", *formatPtr)
	}

	return nil
}

func analyzeAndRender(dir string, outDir string, stdout io.Writer, progress io.Writer) error {
	workspace, g, mdSummary, jsonBytes, err := analyzeWorkspace(dir, progress)
	if err != nil {
		return err
	}

	if outDir != "" {
		absOutDir, err := filepath.Abs(outDir)
		if err != nil {
			return fmt.Errorf("invalid output directory: %w", err)
		}
		if err := writeArtifacts(g, mdSummary, jsonBytes, absOutDir, progress); err != nil {
			return err
		}
		return nil
	}

	if progress != nil {
		fmt.Fprintln(progress, "\n-----------------------------------------")
		fmt.Fprintln(progress, "📊 Analysis Summary")
		fmt.Fprintln(progress, "-----------------------------------------")
	}
	fmt.Fprint(stdout, mdSummary)
	if progress != nil {
		fmt.Fprintln(progress, "-----------------------------------------")
		fmt.Fprintf(progress, "✨ Done for %s\n", workspace)
	}
	return nil
}

func analyzeWorkspace(dir string, progress io.Writer) (string, *graph.Graph, string, []byte, error) {
	workspace, err := filepath.Abs(dir)
	if err != nil {
		return "", nil, "", nil, fmt.Errorf("invalid directory path: %w", err)
	}

	if progress != nil {
		fmt.Fprintln(progress, "🚀 Graphify-Go Scanner Starting 🚀")
		fmt.Fprintf(progress, "📂 Target Directory: %s\n", workspace)
		fmt.Fprintln(progress, "-----------------------------------------")
	}

	start := time.Now()
	results, err := parser.ProcessWorkspace(workspace, progress)
	if err != nil {
		return "", nil, "", nil, fmt.Errorf("failed to process workspace: %w", err)
	}

	if len(results) == 0 {
		return "", nil, "", nil, fmt.Errorf("no supported source files found in the directory")
	}

	if progress != nil {
		fmt.Fprintln(progress, "🏗️ Building Graph...")
	}
	builder := graph.NewBuilder()
	g := builder.Build(results)
	if progress != nil {
		fmt.Fprintf(progress, "✅ Graph built: %d Nodes, %d Edges\n", len(g.Nodes), len(g.Edges))
		fmt.Fprintln(progress, "🧠 Running Community Detection (Louvain)...")
	}

	cluster.DetectCommunities(g)

	if progress != nil {
		commCount := make(map[int]int)
		for _, n := range g.Nodes {
			commCount[n.Community]++
		}
		fmt.Fprintf(progress, "✅ Discovered %d Communities.\n", len(commCount))
		fmt.Fprintf(progress, "⏱️ Analysis completed in %v\n", time.Since(start))
	}

	mdSummary := export.ExportSystemGraphMD(g)
	jsonBytes, err := g.ToJSON()
	if err != nil {
		return "", nil, "", nil, fmt.Errorf("failed to serialize graph to JSON: %w", err)
	}

	return workspace, g, mdSummary, jsonBytes, nil
}

func writeArtifacts(g *graph.Graph, mdSummary string, jsonBytes []byte, outDir string, progress io.Writer) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	jsonPath := filepath.Join(outDir, "graph.json")
	if err := os.WriteFile(jsonPath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write graph.json: %w", err)
	}

	htmlPath := filepath.Join(outDir, "graph.html")
	if err := export.ExportSystemGraphHTML(g, htmlPath); err != nil && progress != nil {
		fmt.Fprintf(progress, "⚠️ Skipped HTML export: %v\n", err)
	}

	mdPath := filepath.Join(outDir, "system-graph.md")
	if err := os.WriteFile(mdPath, []byte(mdSummary), 0644); err != nil {
		return fmt.Errorf("failed to write system-graph.md: %w", err)
	}

	if progress != nil {
		fmt.Fprintln(progress, "\n-----------------------------------------")
		fmt.Fprintf(progress, "💾 Reports saved to: %s\n", outDir)
		fmt.Fprintln(progress, "  📄 graph.json")
		fmt.Fprintln(progress, "  🌐 graph.html")
		fmt.Fprintln(progress, "  📝 system-graph.md")
	}

	return nil
}
