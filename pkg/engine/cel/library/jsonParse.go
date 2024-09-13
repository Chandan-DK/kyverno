package library

import (
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func JsonParseLib() cel.EnvOption {
	return cel.Lib(jsonParseLib)
}

var jsonParseLib = &jsonParse{}

type jsonParse struct{}

func (*jsonParse) LibraryName() string {
	return "kyverno"
}

var jsonParseLibraryDecls = map[string][]cel.FunctionOpt{
	"json_parse": {
		cel.Overload("kyverno_json_parse_string", []*cel.Type{cel.StringType}, cel.DynType,
			cel.UnaryBinding(jsonParseString),
		),
	},
}

func (*jsonParse) CompileOptions() []cel.EnvOption {
	options := []cel.EnvOption{}
	for name, overloads := range jsonParseLibraryDecls {
		options = append(options, cel.Function(name, overloads...))
	}
	return options
}

func (*jsonParse) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func jsonParseString(val ref.Val) ref.Val {
	jsonString, ok := val.Value().(string)
	if !ok {
		return types.NewErr("expected a string, got %T", val.Type())
	}

	var parsedJSON any
	if err := json.Unmarshal([]byte(jsonString), &parsedJSON); err != nil {
		return types.NewErr("error while parsing JSON: %v", err)
	}

	return convertToCelValue(parsedJSON)
}

func convertToCelValue(value any) ref.Val {
	fmt.Printf("Value Type: %T", value)
	switch v := value.(type) {
	case map[string]any:
		return types.NewStringInterfaceMap(types.DefaultTypeAdapter, v)
	case []any:
		return types.NewDynamicList(types.DefaultTypeAdapter, v)
	case string:
		return types.String(v)
	case float64:
		return types.Double(v)
	case bool:
		return types.Bool(v)
	case nil:
		return types.NullValue
	default:
		return types.NewErr("unsupported JSON type %T", v)
	}
}
