package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
)

var (
	importsTemplate = template.Must(template.New("importsTemplate").Parse(`
import ({{ range .Imports }}
	{{ if .Name }}{{ .Name }} {{ end }}{{ .Path }}{{ end }}
	"{{ .Config.MockPackage }}"
)`))
	mockStructTemplate = template.Must(template.New("mockStructTemplate").Parse(`
{{ range $name, $val := .Interfaces }}{{ if $val.IsMockOfMocks }}
type {{ $name }}Mock struct { {{ range $methodName, $m := $val.Methods }}
	{{ $methodName }}Mock {{ range $m.Results }}{{ .Type }}Mock{{ end }}{{ end }}
}
{{ else }}
type {{ $name }}Mock struct { {{ $.Config.EmbeddedMockStruct }} }
{{ end }}{{ end }}
`))
	mockMethodsTemplate = template.Must(template.New("mockMethodsTemplate").Parse(`
{{ range $name, $intr := .Interfaces }}{{ range $methodName, $method := .Methods }}
func (m *{{ $name }}Mock) {{ $methodName }}({{
	range $method.Parameters
}}{{
		.Name
		}} {{
		.Type
		}},{{
	end
}}) {{
	if $method.Results
}}({{
		range $method.Results
		}}{{
			if and $intr.IsMockOfMocks $.Config.LocalPackage
				}}{{ 
				$.PackageName
				}}.{{
			end
			}}{{
			.Type
			}}, {{
		end
}}){{
	end
}}{
{{ 
	if $intr.IsMockOfMocks
}}	return &m.{{ $methodName }}Mock{{
	else
}}	args := m.Called({{ 
		range $method.Parameters 
		}}{{
			.Name
		}},{{
		end
	}})
	return {{
		range $i, $r := $method.Results
		}}{{ 
			if $i 
				}}, {{ 
			end 
			}}{{ 
			if eq $r.Type "error"
			}}args.Error({{ $i }}){{
			else if eq $r.Type "int"
			}}args.Int({{ $i }}){{
			else if eq $r.Type "bool"
			}}args.Bool({{ $i }}){{
			else if eq $r.Type "string"
			}}args.String({{ $i }}){{
			else
			}}args.Get({{ $i }}).({{ $r.Type }}){{
			end
		}}{{
		end 
	}}{{
	end
}}
}
{{ end }}{{ end }}
`))
)

type Config struct {
	Package            string  `yaml:"-"`
	MockPackage        string  `yaml:"mock_package"`
	LocalPackage       *string `yaml:"-"`
	EmbeddedMockStruct string  `yaml:"embedded_mock"`
	Debug              bool    `yaml:"debug"`
}

func (c *Config) ReadConfig(file string) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

type Import struct {
	Name *string
	Path string
}

func (f *Import) DebugString() string {
	if f.Name == nil {
		return fmt.Sprintf("{(nil) %s}", f.Path)
	}

	return fmt.Sprintf("{%s %s}", *f.Name, f.Path)
}

type Field struct {
	Name *string
	Type string
}

func (f *Field) DebugString() string {
	if f.Name == nil {
		return fmt.Sprintf("{(nil) %s}", f.Type)
	}

	return fmt.Sprintf("{%s %s}", *f.Name, f.Type)
}

type Method struct {
	Parameters []Field
	Results    []Field
}

type Interface struct {
	Methods       map[string]Method
	IsMockOfMocks bool
}

type GenStore struct {
	Config       *Config
	Imports      []Import
	Interfaces   map[string]Interface
	LocalStructs map[string]struct{}
	PackageName  string
	Fset         *token.FileSet
}

func (store *GenStore) Debug(format string, a ...interface{}) {
	if store.Config.Debug {
		fmt.Printf(format, a...)
	}
}

func (store *GenStore) CollectType(t ast.Expr) (string, error) {
	switch expr := t.(type) {
	case *ast.Ident:
		return expr.Name, nil
	case *ast.SelectorExpr:
		tt, err := store.CollectType(expr.X)
		if err != nil {
			return "", fmt.Errorf("collect type of %#v: %w", expr.X, err)
		}

		return fmt.Sprintf("%s.%s", tt, expr.Sel.Name), nil
	case *ast.StarExpr:
		tt, err := store.CollectType(expr.X)
		if err != nil {
			return "", fmt.Errorf("collect type of %#v: %w", expr.X, err)
		}

		return fmt.Sprintf("*%s", tt), nil
	case *ast.ArrayType:
		if expr.Len != nil {
			return "", fmt.Errorf("accept only slice types by received %#v", expr)
		}
		tt, err := store.CollectType(expr.Elt)
		if err != nil {
			return "", fmt.Errorf("collect type of %#v: %w", expr.Elt, err)
		}

		return fmt.Sprintf("[]%s", tt), nil
	case *ast.MapType:
		tkey, err := store.CollectType(expr.Key)
		if err != nil {
			return "", fmt.Errorf("collect type of %#v: %w", expr.Key, err)
		}
		tval, err := store.CollectType(expr.Value)
		if err != nil {
			return "", fmt.Errorf("collect type of %#v: %w", expr.Value, err)
		}
		return fmt.Sprintf("map[%s]%s", tkey, tval), nil
	case *ast.InterfaceType:
		if expr.Methods != nil && len(expr.Methods.List) > 0 {
			return "", fmt.Errorf("move interface methods to named interface")
		}
		return "interface{}", nil
	default:
		return "", fmt.Errorf("unknown type %#v", t)
	}
}

func (store *GenStore) ProcessField(f *ast.Field) (Field, error) {
	if f == nil {
		return Field{}, fmt.Errorf("field is nil")
	}

	var fieldName *string

	if len(f.Names) > 0 {
		fieldName = &f.Names[0].Name
	}

	fieldType, err := store.CollectType(f.Type)
	if err != nil {
		return Field{}, fmt.Errorf("error on collect type: %w", err)
	}

	return Field{
		Name: fieldName,
		Type: fieldType,
	}, nil
}

func (store *GenStore) ProcessMethod(fn *ast.FuncType) (Method, error) {
	parameters := make([]Field, 0)

	store.Debug("  <-\n")

	for i, param := range fn.Params.List {
		field, err := store.ProcessField(param)
		if err != nil {
			return Method{}, fmt.Errorf("process %dth param: %w", i+1, err)
		}

		if field.Name == nil {
			arg := fmt.Sprintf("arg%d", i)
			field.Name = &arg
		}

		store.Debug("    %s\n", field.DebugString())

		parameters = append(parameters, field)
	}

	results := make([]Field, 0)

	if fn.Results != nil && len(fn.Results.List) > 0 {
		store.Debug("  ->\n")

		for i, result := range fn.Results.List {
			field, err := store.ProcessField(result)
			if err != nil {
				return Method{}, fmt.Errorf("process %dth result: %w", i+1, err)
			}
			store.Debug("    %s\n", field.DebugString())

			results = append(results, field)
		}
	}

	return Method{
		Parameters: parameters,
		Results:    results,
	}, nil
}

func (store *GenStore) ProcessInterfaceType(expr *ast.InterfaceType) (Interface, error) {
	if expr.Methods == nil {
		return Interface{}, fmt.Errorf("methods is nil")
	}

	methods := make(map[string]Method)

	for _, field := range expr.Methods.List {
		if field == nil {
			return Interface{}, fmt.Errorf("field is nil")
		}

		methodName := field.Names[0].Name
		store.Debug("* %s\n", methodName)

		fn, ok := field.Type.(*ast.FuncType)
		if !ok {
			log.Fatalf("field is not FuncType")
		}

		m, err := store.ProcessMethod(fn)
		if err != nil {
			return Interface{}, fmt.Errorf("process method %s: %w", methodName, err)
		}

		methods[methodName] = m
	}

	return Interface{
		Methods: methods,
	}, nil
}

func (store *GenStore) ProcessTokenType(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch expr := spec.Type.(type) {
		case *ast.InterfaceType:
			store.Debug("interface %s\n", spec.Name.Name)

			if expr.Incomplete {
				log.Fatalf("%s is incomplete interface", spec.Name.Name)
			}

			intr, err := store.ProcessInterfaceType(expr)
			if err != nil {
				log.Fatalf("%s ProcessInterfaceType: %v", spec.Name.Name, err)
			}

			store.Interfaces[spec.Name.Name] = intr

		case *ast.StructType:
			store.Debug("struct %s\n", spec.Name.Name)

			if expr.Incomplete {
				log.Fatalf("%s is incomplete interface", spec.Name.Name)
			}

			store.LocalStructs[spec.Name.Name] = struct{}{}
		default:
			continue
		}
	}
}

func (store *GenStore) ProcessImportType(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}
		store.Debug("ProcessImportType: spec = %#v\n", spec)

		var importName *string
		if spec.Name != nil {
			importName = &spec.Name.Name
		}
		imp := Import{
			Name: importName,
			Path: spec.Path.Value,
		}
		store.Debug("ProcessImportType: imp = %s\n", imp.DebugString())

		store.Imports = append(store.Imports, imp)
	}
}

func (store *GenStore) ProcessGenDecl(decl *ast.GenDecl) {
	switch decl.Tok {
	case token.TYPE:
		store.ProcessTokenType(decl)
	case token.IMPORT:
		store.ProcessImportType(decl)
	}
}

func (store *GenStore) FixMockOfMocks() {
	for name, i := range store.Interfaces {
		isMockOfMocks := true

		for _, m := range i.Methods {
			if len(m.Parameters) > 0 {
				isMockOfMocks = false
				break
			}

			if len(m.Results) != 1 {
				isMockOfMocks = false
				break
			}

			if _, exists := store.Interfaces[m.Results[0].Type]; !exists {
				isMockOfMocks = false
				break
			}
		}

		i.IsMockOfMocks = isMockOfMocks
		if isMockOfMocks {
			store.Debug("%s isMockOfMocks\n", name)
		}
		store.Interfaces[name] = i
	}
}

func (store *GenStore) fixLocalFields(fields []Field) {
	for idx, f := range fields {
		var local bool
		if _, exists := store.LocalStructs[f.Type]; exists {
			local = true
		}

		if _, exists := store.Interfaces[f.Type]; exists {
			local = true
		}

		if !local {
			continue
		}

		newType := fmt.Sprintf("%s.%s", store.PackageName, f.Type)
		store.Debug("Fixed %s => %s\n", f.Type, newType)

		f.Type = newType
		fields[idx] = f
	}
}

func (store *GenStore) FixLocalFields() {
	for _, i := range store.Interfaces {
		if i.IsMockOfMocks {
			continue
		}

		for _, m := range i.Methods {
			store.fixLocalFields(m.Parameters)
			store.fixLocalFields(m.Results)
		}
	}
}

func (store *GenStore) Generate(out io.Writer) {
	fmt.Fprintf(out, "// This file was generated by portmocksgen. DO NOT EDIT.\npackage %s\n", store.Config.Package)

	if err := importsTemplate.Execute(out, store); err != nil {
		log.Fatalf("importsTemplate.Execute: %v", err)
	}

	if err := mockStructTemplate.Execute(out, store); err != nil {
		log.Fatalf("mockStructTemplate.Execute: %v", err)
	}

	if err := mockMethodsTemplate.Execute(out, store); err != nil {
		log.Fatalf("mockMethodsTemplate.Execute: %v", err)
	}
}

func process(config Config, filein string, fileout string) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filein, nil, 0)
	if err != nil {
		log.Fatalf("parse file: %v", err)
	}

	out, err := os.Create(fileout)
	if err != nil {
		log.Fatalf("create file: %v", err)
	}

	store := GenStore{
		Config:       &config,
		Imports:      make([]Import, 0),
		Interfaces:   make(map[string]Interface),
		LocalStructs: make(map[string]struct{}),
		Fset:         fset,
		PackageName:  node.Name.Name,
	}

	if config.LocalPackage != nil {
		store.Imports = append(store.Imports, Import{
			Path: fmt.Sprintf(`"%s"`, *config.LocalPackage),
		})
	}

	for _, decl := range node.Decls {
		switch a := decl.(type) {
		case *ast.GenDecl:
			store.ProcessGenDecl(a)
		default:
			continue
		}
	}

	store.FixMockOfMocks()

	if config.LocalPackage != nil {
		store.FixLocalFields()
	}

	var buf bytes.Buffer
	store.Generate(&buf)

	p, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("format error: %v", err)
	}
	out.Write(p)
	out.Close()
}

func main() {
	configPath := flag.String("config", "portmocksgen.yml", "Shared config for all files")
	fileIn := flag.String("in", "", "input file with interfaces")
	fileOut := flag.String("out", "", "output file where mocks would be stored")
	fileInPackage := flag.String("inpkg", "", "input file full package. If not set equals outpkg.")
	fileOutPackage := flag.String("outpkg", "", "output file package name")
	flag.Parse()

	config := Config{}
	config.ReadConfig(*configPath)

	var errors error

	if *fileIn == "" {
		errors = multierror.Append(errors, fmt.Errorf("--in is not set"))
	}

	if *fileOut == "" {
		errors = multierror.Append(errors, fmt.Errorf("--out is not set"))
	}

	if *fileOutPackage == "" {
		errors = multierror.Append(errors, fmt.Errorf("--outpkg is not set"))
	} else {
		config.Package = *fileOutPackage
	}

	if *fileInPackage != "" {
		config.LocalPackage = fileInPackage
	}

	if errors != nil {
		log.Fatalf("%v", errors)
	}

	process(config, *fileIn, *fileOut)
	fmt.Printf("portsmockgen: %s -> %s\n", *fileIn, *fileOut)
}
