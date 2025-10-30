package builtins

import (
	"bufio"
	"fmt"
	"gloob-interpreter/internal/parser"
	"gloob-interpreter/internal/scope"
	"gloob-interpreter/internal/values"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// SetupNativeFunctions adds all built-in native functions to the scope
func SetupNativeFunctions(s *scope.Scope) {
	// Math functions
	DeclareNativeFunction(s, "abs", AbsFunction)
	DeclareNativeFunction(s, "round", RoundFunction)
	DeclareNativeFunction(s, "max", MaxFunction)
	DeclareNativeFunction(s, "min", MinFunction)
	DeclareNativeFunction(s, "random", RandomFunction)
	DeclareNativeFunction(s, "randInt", RandIntFunction)

	// I/O functions
	DeclareNativeFunction(s, "input", InputFunction)
	DeclareNativeFunction(s, "print", PrintFunction)
	DeclareNativeFunction(s, "println", PrintlnFunction)

	// Note: len() works with both strings and arrays but is kept as a standalone
	// function for convenience. For consistency, .len() method is also available.
	DeclareNativeFunction(s, "len", LenFunction)

	// Type conversion functions
	DeclareNativeFunction(s, "number", NumberFunction)
	DeclareNativeFunction(s, "string", StringFunction)
	DeclareNativeFunction(s, "bool", BoolFunction)
	DeclareNativeFunction(s, "type", TypeFunction)

	// System functions
	DeclareNativeFunction(s, "sleep", SleepFunction)
	DeclareNativeFunction(s, "clear", ClearFunction)

	// Note: String methods (upper, lower, trim, contains, split, replace, indexOf)
	// are now available as string methods: "hello".upper(), "text".split(" "), etc.
	// Array methods (contains, indexOf, join, reverse) are available as: arr.contains(x), arr.join(", "), etc.
}

func DeclareNativeFunction(s *scope.Scope, name string, expression func(args []values.RuntimeValue, scope interface{}) values.RuntimeValue) {
	s.Declare(name, &values.NativeFunctionValue{
		Type:       parser.NodeTypeNativeFunction,
		Expression: expression,
	}, true)
}

// PrintFunction prints arguments to stdout without newline
func PrintFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Print("\n")
	return &values.NullValue{Type: parser.NodeTypeNull}
}

// PrintlnFunction prints arguments to stdout with newline
func PrintlnFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Print("\n")
	return &values.NullValue{Type: parser.NodeTypeNull}
}

func InputFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	prompt := ""
	if len(args) > 0 {
		prompt = fmt.Sprint(args[0])
	}
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
		return nil
	}
	// Trim the newline character but keep the string as is
	value = strings.TrimSpace(value)
	return &values.StringValue{Type: parser.NodeTypeString,
		Value: value,
	}
}

func RandomFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	return &values.NumericValue{Type: parser.NodeTypeNumeric,
		Value: rand.Float64(),
	}
}

func RandIntFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {

	if len(args) > 2 {
		fmt.Printf("RandInt function expects 1 or 2 arguments\n")
		os.Exit(1)
		return nil
	}

	var min *values.NumericValue = &values.NumericValue{Type: parser.NodeTypeNumeric, Value: 0}
	var limit *values.NumericValue = &values.NumericValue{Type: parser.NodeTypeNumeric, Value: 100}
	var ok bool
	if len(args) > 1 {
		min, ok = args[0].(*values.NumericValue)
		if !ok {
			fmt.Printf("RandInt function expects a numeric argument\n")
			os.Exit(1)
			return nil
		}
		limit, ok = args[1].(*values.NumericValue)
		if !ok {
			fmt.Printf("RandInt function expects a numeric argument\n")
			os.Exit(1)
			return nil
		}
	}

	randomNumber := float64(rand.Intn(int(limit.Value)-int(min.Value)+1) + int(min.Value))

	return &values.NumericValue{Type: parser.NodeTypeNumeric,
		Value: randomNumber,
	}
}

func AbsFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	number, ok := args[0].(*values.NumericValue)
	if !ok {
		fmt.Printf("Abs function expects a numeric argument\n")
		os.Exit(1)
		return nil
	}
	return &values.NumericValue{Type: parser.NodeTypeNumeric,
		Value: math.Abs(number.Value),
	}
}

func RoundFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	number, ok := args[0].(*values.NumericValue)
	if !ok {
		fmt.Printf("Round function expects a numeric argument\n")
		os.Exit(1)
		return nil
	}
	return &values.NumericValue{Type: parser.NodeTypeNumeric,
		Value: math.Round(number.Value),
	}
}

func MaxFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	number1, ok := args[0].(*values.NumericValue)
	if !ok {
		fmt.Printf("Max function expects a numeric argument\n")
		os.Exit(1)
		return nil
	}
	number2, ok := args[1].(*values.NumericValue)
	if !ok {
		fmt.Printf("Max function expects a numeric argument\n")
		os.Exit(1)
		return nil
	}
	return &values.NumericValue{Type: parser.NodeTypeNumeric,
		Value: math.Max(number1.Value, number2.Value),
	}
}

func MinFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	number1, ok := args[0].(*values.NumericValue)
	if !ok {
		fmt.Printf("Min function expects a numeric argument\n")
		os.Exit(1)
		return nil
	}
	number2, ok := args[1].(*values.NumericValue)
	if !ok {
		fmt.Printf("Min function expects a numeric argument\n")
		os.Exit(1)
		return nil
	}
	return &values.NumericValue{Type: parser.NodeTypeNumeric,
		Value: math.Min(number1.Value, number2.Value),
	}
}

func LenFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	if len(args) != 1 {
		fmt.Printf("len() expects 1 argument, got %d\n", len(args))
		os.Exit(1)
		return nil
	}

	// Handle strings
	if strVal, ok := args[0].(*values.StringValue); ok {
		return &values.NumericValue{
			Type:  parser.NodeTypeNumeric,
			Value: float64(len(strVal.Value)),
		}
	}

	// Handle arrays
	if arrVal, ok := args[0].(*values.ArrayValue); ok {
		return &values.NumericValue{
			Type:  parser.NodeTypeNumeric,
			Value: float64(len(arrVal.Elements)),
		}
	}

	fmt.Printf("len() expects a string or array argument\n")
	os.Exit(1)
	return nil
}

func NumberFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	stringValue, ok := args[0].(*values.StringValue)
	if !ok {
		fmt.Printf("Number function expects a string argument\n")
		os.Exit(1)
		return nil
	}
	value, err := strconv.ParseFloat(stringValue.Value, 64)
	if err != nil {
		fmt.Printf("Error parsing number: %v\n", err)
		os.Exit(1)
		return nil
	}
	return &values.NumericValue{Type: parser.NodeTypeNumeric,
		Value: value,
	}
}

func StringFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	numberValue, ok := args[0].(*values.NumericValue)
	if !ok {
		fmt.Printf("String function expects a numeric argument\n")
		os.Exit(1)
		return nil
	}
	return &values.StringValue{Type: parser.NodeTypeString,
		Value: strconv.FormatFloat(numberValue.Value, 'f', -1, 64),
	}
}

func BoolFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	boolValue, ok := args[0].(*values.StringValue)
	if !ok {
		fmt.Printf("Bool function expects a boolean argument\n")
		os.Exit(1)
		return nil
	}
	return &values.BooleanValue{Type: parser.NodeTypeBoolean, Value: boolValue.Value == "true"}
}

func TypeFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	typeValue := args[0]
	return &values.StringValue{
		Type:  parser.NodeTypeString,
		Value: strings.ToLower(fmt.Sprint(typeValue.NodeType())),
	}
}

func SleepFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	// Convert seconds to milliseconds to handle decimal values
	duration := time.Duration(args[0].(*values.NumericValue).Value*1000) * time.Millisecond
	time.Sleep(duration)
	return &values.NullValue{Type: parser.NodeTypeNull}
}

func ClearFunction(args []values.RuntimeValue, scope interface{}) values.RuntimeValue {
	fmt.Printf("\x1b[2J")
	return &values.NullValue{Type: parser.NodeTypeNull}
}
