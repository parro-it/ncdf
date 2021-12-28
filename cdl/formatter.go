package cdl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/parro-it/ncdf/ordmap"
	"github.com/parro-it/ncdf/types"
)

// CDL ...
func CDLFile(f *types.File) string {
	var res strings.Builder
	res.WriteString("netcdf filename {\n")
	res.WriteString(dimensionsCDL(f.Dimensions))
	res.WriteString("variables\n")
	for _, v := range f.Vars.Values() {
		res.WriteString("    ")
		res.WriteString(CDLVar(&v))
	}
	res.WriteRune('\n')
	res.WriteString("// global attributes:\n")
	res.WriteString(attributesCDL(f.Attrs, ""))
	res.WriteRune('}')

	return res.String()
}

// CDL ...
func CDLDimension(f *types.Dimension) string {
	return fmt.Sprintf("%s = %d;", f.Name, f.Len)
}

// CDL ...
func CDLAttr(f *types.Attr) string {
	switch v := f.Val.(type) {
	case []int:
		return fmt.Sprintf("%s = %d;", f.Name, v[0])
	case []int32:
		return fmt.Sprintf("%s = %d;", f.Name, v[0])
	case []int16:
		return fmt.Sprintf("%s = %d;", f.Name, v[0])
	case []byte:
		return fmt.Sprintf("%s = %d;", f.Name, v[0])
	case string:
		return fmt.Sprintf(`%s = "%s";`, f.Name, v)
	case []float32:
		return fmt.Sprintf("%s = %f;", f.Name, v[0])
	case []float64:
		return fmt.Sprintf("%s = %f;", f.Name, v[0])

	case int:
		return fmt.Sprintf("%s = %d;", f.Name, v)
	case int32:
		return fmt.Sprintf("%s = %d;", f.Name, v)
	case int16:
		return fmt.Sprintf("%s = %d;", f.Name, v)
	case byte:
		return fmt.Sprintf("%s = %d;", f.Name, v)
	case float32:
		return fmt.Sprintf("%s = %f;", f.Name, v)
	case float64:
		return fmt.Sprintf("%s = %f;", f.Name, v)
	}

	return fmt.Sprintf("~UNKNOWN TYPE %v~", reflect.TypeOf(f.Val))
}

// CDL ...
func CDLVar(v *types.Var) string {
	var dimS strings.Builder

	for i, d := range v.Dimensions {
		if i > 0 {
			dimS.WriteString(", ")
		}
		dimS.WriteString(d.Name)
	}

	var res strings.Builder
	res.WriteString(fmt.Sprintf("%s %s(%s);\n", CDLType(v.Type), v.Name, dimS.String()))

	res.WriteString(attributesCDL(v.Attrs, v.Name))
	return res.String()
}

func attributesCDL(attrs ordmap.OrderedMap[types.Attr, string], prefix string) string {
	var res strings.Builder
	for _, a := range attrs.Values() {
		res.WriteString("        ")
		res.WriteString(prefix + ":")
		res.WriteString(CDLAttr(&a))
		res.WriteRune('\n')
	}
	return res.String()
}

// CDL ...
func CDLType(t types.Type) string {
	switch t {
	case types.Byte:
		return "byte"
	case types.Char:
		return "char"
	case types.Short:
		return "short"
	case types.Int:
		return "int"
	case types.Float:
		return "float"
	case types.Double:
		return "double"
	}

	return fmt.Sprintf("[unknown type:%d]", t)
}

func dimensionsCDL(dd []types.Dimension) string {
	var res strings.Builder
	res.WriteString("dimensions:\n")
	for _, d := range dd {
		res.WriteString("    ")
		res.WriteString(CDLDimension(&d))
		res.WriteRune('\n')
	}
	return res.String()
}
