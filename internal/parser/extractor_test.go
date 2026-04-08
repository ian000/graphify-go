package parser_test

import (
	"testing"

	"github.com/kings2017/graphify-go/internal/parser"
)

func TestExtractor_ParseAndExtract_JavaScript(t *testing.T) {
	registry := parser.NewRegistry()
	extractor := parser.NewExtractor(registry)

	// Mock 包含 class, method, function, call, import 的 JS 代码
	jsCode := []byte(`
		import { helper } from './utils';

		class MyService {
			constructor() {
				this.init();
			}
			init() {
				helper();
			}
		}

		function standaloneFunc() {
			const s = new MyService();
			s.init();
		}
	`)

	data, err := extractor.ParseAndExtract("test.js", jsCode)
	if err != nil {
		t.Fatalf("ParseAndExtract failed: %v", err)
	}

	if len(data.Classes) != 1 || data.Classes[0] != "MyService" {
		t.Errorf("Expected 1 class 'MyService', got %v", data.Classes)
	}

	if len(data.Funcs) != 1 || data.Funcs[0] != "standaloneFunc" {
		t.Errorf("Expected 1 function 'standaloneFunc', got %v", data.Funcs)
	}

	// methods 应该包含 init (constructor 在 JS 里被识别为 method_definition)
	foundInit := false
	for _, m := range data.Methods {
		if m == "init" {
			foundInit = true
			break
		}
	}
	if !foundInit {
		t.Errorf("Expected method 'init' to be captured, got %v", data.Methods)
	}

	if len(data.Imports) != 1 || data.Imports[0] != "'./utils'" {
		t.Errorf("Expected import './utils', got %v", data.Imports)
	}
}

func TestExtractor_ParseAndExtract_Go(t *testing.T) {
	registry := parser.NewRegistry()
	extractor := parser.NewExtractor(registry)

	goCode := []byte(`
		package testpkg

		import "fmt"

		type Server struct {}

		func (s *Server) Start() {
			fmt.Println("started")
		}

		func Helper() {
			s := &Server{}
			s.Start()
		}
	`)

	data, err := extractor.ParseAndExtract("test.go", goCode)
	if err != nil {
		t.Fatalf("ParseAndExtract failed: %v", err)
	}

	if len(data.Classes) != 1 || data.Classes[0] != "Server" {
		t.Errorf("Expected 1 struct 'Server', got %v", data.Classes)
	}

	if len(data.Funcs) != 1 || data.Funcs[0] != "Helper" {
		t.Errorf("Expected 1 function 'Helper', got %v", data.Funcs)
	}

	if len(data.Methods) != 1 || data.Methods[0] != "Start" {
		t.Errorf("Expected 1 method 'Start', got %v", data.Methods)
	}

	if len(data.Imports) != 1 || data.Imports[0] != "\"fmt\"" {
		t.Errorf("Expected import '\"fmt\"', got %v", data.Imports)
	}
}
