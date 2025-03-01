package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/yiigo/yiigo/internal"
	"golang.org/x/tools/go/packages"
)

type GenBody struct {
	PkgName string
	Imports map[string]struct{}
	Structs []Struct
}

// Struct represents a Go struct
type Struct struct {
	Receiver string
	Name     string
	Fields   []Field
}

// Field represents a struct field
type Field struct {
	Name    string
	Type    string
	Default Default
}

type Default struct {
	GenType *GenType // 泛型
	Value   string
}

type GenType struct {
	Ident string
	Type  string
}

// Template for generating Get methods
const tmpl = `package {{ .PkgName }}

// Code generated by gg; DO NOT EDIT.

{{- with .Imports }}
{{- $len := len . }}
{{- if gt $len 0 }}

{{- if eq $len 1 }}
{{- range $k, $v := . }}

import "{{ $k }}"

{{- end }}
{{- else }}

import (
{{- range $k, $v := . }}
	"{{ $k }}"
{{- end }}
)

{{- end }}
{{- end }}

{{- end }}

{{- range $s := .Structs }}

// Get methods for {{ $s.Name }}

{{- range $s.Fields }}

func ({{ $s.Receiver }} *{{ $s.Name }}) Get{{ .Name }}() {{ .Type }} {
	if {{ $s.Receiver }} != nil {
		return {{ $s.Receiver }}.{{ .Name }}
	}
	{{- if .Default.GenType }}
	var v {{ .Default.GenType.Ident }}
	return v
	{{- else }}
	return {{ .Default.Value }}
	{{- end }}
}

{{- end }}

{{- end }}
`

func main() {
	cmd := &cobra.Command{
		Use:     "gg",
		Short:   "生成Get方法^_^",
		Long:    "为结构体生成`Get`方法，避免空指针导致Panic",
		Version: "v1.0.0",
		Example: internal.CmdExamples(
			"gg xxx.go",
			"//go:generate gg xxx.go",
		),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("usage: gg <source file>")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, filename := range args {
				genGetter(filename)
			}
		},
	}
	// 执行
	_ = cmd.Execute()
}

// 自定义导入逻辑
type customImporter struct{}

func (ci *customImporter) Import(path string) (*types.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedDeps | packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("packages.Load: %w", err)
	}
	if len(pkgs) > 0 {
		return pkgs[0].Types, nil
	}
	return nil, fmt.Errorf("package not found: %s", path)
}

func genGetter(filename string) {
	sourceFile := filepath.Clean(filename)
	fset := token.NewFileSet()

	// Parse the source file
	node, err := parser.ParseFile(fset, sourceFile, nil, parser.AllErrors)
	if err != nil {
		log.Fatalln("parser.ParseFile:", internal.FmtErr(err))
	}

	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	conf := &types.Config{
		Importer: &customImporter{},
	}
	_, err = conf.Check("", fset, []*ast.File{node}, info)
	if err != nil {
		log.Fatalln("conf.Check:", internal.FmtErr(err))
	}

	gen := GenBody{
		PkgName: node.Name.Name,
		Imports: make(map[string]struct{}),
		Structs: make([]Struct, 0),
	}

	// Walk through the AST
	ast.Inspect(node, func(n ast.Node) bool {
		// Find type declarations
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if it's a struct
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Generics identifier
		var gt []GenType
		if ts.TypeParams != nil {
			for _, field := range ts.TypeParams.List {
				for _, ident := range field.Names {
					gt = append(gt, GenType{
						Ident: ident.Name,
						Type:  info.TypeOf(field.Type).String(),
					})
				}
			}
		}

		imports, data := buildStruct(info, ts, st, gt)
		for _, v := range imports {
			gen.Imports[v] = struct{}{}
		}
		gen.Structs = append(gen.Structs, data)

		return true
	})

	// Generate code
	var buf bytes.Buffer
	t := template.Must(template.New("getters").Parse(tmpl))
	if err := t.Execute(&buf, gen); err != nil {
		log.Fatalln("t.Execute:", internal.FmtErr(err))
	}

	// Write to a file
	outputFile := strings.ReplaceAll(sourceFile, ".go", "_getter.go")
	if err := os.WriteFile(outputFile, buf.Bytes(), 0o755); err != nil {
		log.Fatalln("os.WriteFile:", internal.FmtErr(err))
	}
	fmt.Println("Generated code saved to", outputFile)
}

func buildStruct(info *types.Info, ts *ast.TypeSpec, st *ast.StructType, gt []GenType) ([]string, Struct) {
	name := ts.Name.String()
	receiver := "s"
	if name != "<nil>" {
		receiver = strings.ToLower(name[:1])
	}
	// Generics identifier
	if len(gt) != 0 {
		name += "["
		name += gt[0].Ident
		for _, v := range gt[1:] {
			name += ", "
			name += v.Ident
		}
		name += "]"
	}

	// Collect imports and fields
	var (
		imports []string
		fields  []Field
	)
	for _, field := range st.Fields.List {
		// Skip embedded fields
		if len(field.Names) == 0 {
			continue
		}

		// 字段类型
		fieldType := info.TypeOf(field.Type).String()
		underlyingType := info.TypeOf(field.Type).Underlying().String()
		// fmt.Println(field.Names, "[fieldType]", fieldType, "[underlyingType]", underlyingType)
		if dotIndex := strings.LastIndex(fieldType, "."); dotIndex >= 0 {
			if pkg := fieldType[:dotIndex]; !filepath.IsAbs(pkg) {
				imports = append(imports, pkg)
				if slashIndex := strings.LastIndex(fieldType, "/"); slashIndex >= 0 {
					fieldType = fieldType[slashIndex+1:]
				}
			} else {
				fieldType = fieldType[dotIndex+1:]
			}
		}

		for _, name := range field.Names {
			fields = append(fields, Field{
				Name:    name.Name,
				Type:    fieldType,
				Default: getTypeValue(field.Type, fieldType, underlyingType, gt),
			})
		}
	}

	return imports, Struct{
		Receiver: receiver,
		Name:     name,
		Fields:   fields,
	}
}

// 获取类型的默认值
func getTypeValue(expr ast.Expr, fieldType, underlyingType string, gt []GenType) Default {
	switch expr.(type) {
	case *ast.Ident, *ast.SelectorExpr: // 基本类型 || 自定义类型 || 包路径类型
		return getDefaultValue(fieldType, underlyingType, gt)
	case *ast.ArrayType, *ast.MapType, *ast.InterfaceType, *ast.StarExpr:
		return Default{Value: "nil"}
	default:
		return Default{Value: "unknown"}
	}
}

func getDefaultValue(fieldType, underlyingType string, gt []GenType) Default {
	// 泛型字段
	for _, v := range gt {
		if fieldType == v.Ident {
			return Default{
				GenType: &GenType{
					Ident: v.Ident,
					Type:  v.Type,
				},
			}
		}
	}
	// 普通字段
	switch underlyingType {
	case "string":
		return Default{Value: `""`}
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return Default{Value: "0"}
	case "float32", "float64":
		return Default{Value: "0"}
	case "bool":
		return Default{Value: "false"}
	case "any":
		return Default{Value: "nil"}
	default:
		if strings.HasPrefix(underlyingType, "*") ||
			strings.HasPrefix(underlyingType, "interface") ||
			strings.HasPrefix(underlyingType, "[]") ||
			strings.HasPrefix(underlyingType, "map") {
			return Default{Value: "nil"}
		}
		if strings.HasPrefix(underlyingType, "struct") {
			return Default{Value: fmt.Sprintf("%s{}", fieldType)}
		}
		return Default{Value: "unknown"}
	}
}
