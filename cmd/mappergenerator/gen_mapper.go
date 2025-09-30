package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
)

type Field struct {
	Name string
	Type string
}

// StructInfo represents metadata about a struct used to generate a mapper.
//
// Name is the name of the struct being mapped.
//
// StructFields contains fields whose types are user-defined structs.
// These require a call to a corresponding Map<Type> function.
//
// StructPtrFields contains field whose types are pointer to user-defined structs
// These require a call to a corresponding Map<Type> function with pre-check for a nil value
//
// SliceFields contains fields that are slices of user-defined structs.
// These require element-wise mapping via Map<Type> in a loop.
type StructInfo struct {
	Name            string
	Fields          []Field
	StructFields    []Field
	StructPtrFields []Field
	SliceFields     []Field
}

var builtinTypes = map[string]bool{
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true,
	"float32": true, "float64": true,
	"complex64": true, "complex128": true,
	"bool": true, "byte": true, "rune": true, "string": true, "any": true, "error": true,
}

func main() {
	fromStructInfos := parse("pkg/mpdapi")
	toStructInfos := parse("internal/api/dto")
	extendedStructInfos := extractStructInfos(fromStructInfos, toStructInfos)
	generate(extendedStructInfos)
}

func extractStructInfos(fromStructInfos []StructInfo, toStructInfos []StructInfo) []StructInfo {
	var structInfos []StructInfo
	// итерируемся по структурам "from" и "to" чтобы найти одноименные
	for _, fromInto := range fromStructInfos {
		for _, toInfo := range toStructInfos {
			if fromInto.Name != toInfo.Name {
				continue
			}
			// имена структур from и to совпали
			var fields []Field
			// итерируемся по полям в найденных структурах, чтобы найти одноименные
			for _, fieldFrom := range fromInto.Fields {
				for _, fieldTo := range toInfo.Fields {
					if fieldFrom.Name != fieldTo.Name {
						continue
					}
					fields = append(fields, Field{Name: fieldFrom.Name})
					break
				}
			}
			info := StructInfo{
				Name:   fromInto.Name,
				Fields: fields,
			}
			// добавляем результат в structInfos
			structInfos = append(structInfos, info)
			break
		}
	}
	extendedStructInfos := make([]StructInfo, len(structInfos))
	// т.к. в structInfo находятся только данные по "обычным" полям, обогащаем ее данными по полям-структурам,
	// полям-слайсам, полям-указателям на структуры. Результат сохраним в extendedStructInfos
	for i, structInfo := range structInfos {
		// итерируемся по собранным на предыдущем цикле данным
		for _, toStructInfo := range toStructInfos {
			// итерируемся по полям-структурам из "to"
			if toStructInfo.Name != structInfo.Name {
				continue
			}
			for _, sliceField := range toStructInfo.SliceFields {
				for _, info := range structInfos {
					if info.Name == sliceField.Type {
						structInfo.SliceFields = append(structInfo.SliceFields, sliceField)
					}
				}
			}
			for _, structField := range toStructInfo.StructFields {
				for _, info := range structInfos {
					if info.Name == structField.Type {
						structInfo.StructFields = append(structInfo.StructFields, structField)
					}
				}
			}
			for _, structPtrField := range toStructInfo.StructPtrFields {
				for _, info := range structInfos {
					if info.Name == structPtrField.Type {
						structInfo.StructPtrFields = append(structInfo.StructPtrFields, structPtrField)
					}
				}
			}
		}
		extendedStructInfos[i] = structInfo
	}
	return extendedStructInfos
}

func generate(extendedStructInfos []StructInfo) {
	t := template.Must(template.ParseFiles("templates/dto_mapper.tmpl"))
	parseData := struct {
		StructInfos []StructInfo
		Notice      string
	}{
		StructInfos: extendedStructInfos,
		Notice:      "Generated code. DO NOT EDIT.",
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, parseData)
	if err != nil {
		log.Fatal(err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("internal/api/dto/generated_mapper.go", formatted, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func parse(dirPath string) []StructInfo {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	var result []StructInfo
	for _, v := range pkgs {
		for _, file := range v.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					// если это не объявление типа
					continue
				}
				for _, spec := range genDecl.Specs {
					spec := spec.(*ast.TypeSpec)
					st, ok := spec.Type.(*ast.StructType)
					if !ok {
						// если это не структура
						continue
					}
					structInfo := StructInfo{
						Name: spec.Name.Name,
					}
					for _, field := range st.Fields.List {
						switch x := field.Type.(type) {
						case *ast.Ident:
							// поле - Ident (или тип из текущего пакета или встроенный тип)
							if builtinTypes[x.Name] {
								// это встроенный тип
								for _, name := range field.Names {
									structInfo.Fields = append(structInfo.Fields, Field{Name: name.Name})
								}
							} else {
								// это структура
								for _, name := range field.Names {
									structInfo.StructFields = append(structInfo.StructFields, Field{Name: name.Name, Type: x.Name})
								}
							}
						case *ast.ArrayType:
							// это слайс
							if t, ok := x.Elt.(*ast.Ident); ok {
								typeName := t.Name
								for _, name := range field.Names {
									structInfo.SliceFields = append(structInfo.Fields, Field{Name: name.Name, Type: typeName})
								}
							}
						case *ast.SelectorExpr:
							if pkgIdent, ok := x.X.(*ast.Ident); ok {
								if pkgIdent.Name == "time" && x.Sel.Name == "Time" {
									for _, name := range field.Names {
										structInfo.Fields = append(structInfo.Fields, Field{Name: name.Name})
									}
								}
							}
						case *ast.StarExpr:
							// это указатель
							if _, ok := x.X.(*ast.Ident); ok {
								if builtinTypes[x.X.(*ast.Ident).Name] {
									// это указатель на встроенный тип
									for _, name := range field.Names {
										structInfo.Fields = append(structInfo.Fields, Field{Name: name.Name})
									}
								} else {
									// это указатель на структуру
									for _, name := range field.Names {
										structInfo.StructPtrFields = append(structInfo.StructPtrFields, Field{Name: name.Name, Type: x.X.(*ast.Ident).Name})
									}
								}
							}
						}
					}
					result = append(result, structInfo)
				}
			}
		}
	}
	return result
}
