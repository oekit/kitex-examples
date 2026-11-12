package main

// #112 护栏：本仓 AGENTS.md 日志节两条硬规则——
// ① 错误字段名统一用 logx.ErrKey 常量（值 "err"），禁止手写 "err"/"error" 字面量；
// ② 业务服务必须接线 logx.WithService/UseSlog（包级 logx.* 经 Default() 输出，
//    默认 logger 必须带 service 字段）。
//
// 本测试扫描两仓业务源码（server/ + client/，排除 kitex_gen 生成物）：
//   - 任何把字面量 "err"/"error" 作为调用参数（logx kv 对的字段名位）即 fail；
//   - main 函数必须包含 logx.UseSlog 接线（示例代码是下游照抄的模板）。

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLogxErrKey_NoHandwrittenLiteral(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file for repo root")
	}
	root := filepath.Dir(filepath.Dir(thisFile)) // server/ 的上一级 = 仓根

	var violations []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, _ := filepath.Rel(root, path)
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "kitex_gen" || name == "testdata" {
				return filepath.SkipDir
			}
			if rel != "." {
				if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			for _, arg := range call.Args {
				lit, ok := arg.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				v := strings.Trim(lit.Value, `"`)
				if v == "err" || v == "error" {
					pos := fset.Position(lit.Pos())
					violations = append(violations,
						pos.String()+"：手写错误字段名 "+lit.Value+"，改用 logx.ErrKey（AGENTS.md 日志节红线，#112）")
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	for _, v := range violations {
		t.Error(v)
	}
}

// TestLogxWithServiceWired 钉住第二条硬规则：两个业务 main 必须接线
// logx.UseSlog（或 WithService），使包级 logx.* 输出带 service 字段。
func TestLogxWithServiceWired(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file for repo root")
	}
	root := filepath.Dir(filepath.Dir(thisFile)) // server/ 的上一级 = 仓根

	for _, dir := range []string{"server", "client"} {
		mainGo := filepath.Join(root, dir, "main.go")
		src, err := os.ReadFile(mainGo)
		if err != nil {
			t.Fatalf("read %s: %v", mainGo, err)
		}
		if !strings.Contains(string(src), "logx.UseSlog(") &&
			!strings.Contains(string(src), "logx.WithService(") {
			t.Errorf("%s：业务 main 未接线 logx.UseSlog/WithService——包级 logx.* 输出将无 service 字段（AGENTS.md 日志节硬规则，#112）", mainGo)
		}
	}
}
