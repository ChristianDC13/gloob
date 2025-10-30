# 🫧 Gloob Programming Language

> A small, playful, interpreted language — made just for fun and experimentation.

Gloob is a dynamically-typed programming language with a clean syntax, 1-based arrays, and a focus on readability. It's simple, expressive, and perfect for learning or quick scripting.

## ✨ Features

- 🎯 **Clean Syntax** - No semicolons, no parentheses in conditionals
- 📦 **1-Based Indexing** - Arrays and strings start at 1 (like humans count!)
- 🔄 **Multiple Loop Types** - Range loops, for-each, while-style, and infinite loops
- 🧩 **Dynamic Typing** - Variables can hold any type
- 🎨 **Boolean Alternatives** - Use `yes`/`no` or `on`/`off` instead of `true`/`false`
- 📚 **Module System** - Import files to organize your code
- 🔗 **Method Chaining** - Call methods on literals: `"hello".upper().len()`
- ⚡ **Implicit Returns** - Last expression in a function is auto-returned

## 🚀 Quick Start

### Installation

**macOS (Recommended):**

Install with a single command:
```bash
curl -fsSL https://raw.githubusercontent.com/ChristianDC13/gloob/main/install.sh | bash
```

This will download the latest pre-built binary and install it to `/usr/local/bin`.

**Build from Source (All Platforms):**

If you're on a different platform or prefer to build from source:

*Prerequisites:* Go 1.16 or higher

```bash
# Clone the repository
git clone https://github.com/ChristianDC13/gloob.git
cd gloob

# Build the interpreter
go build -o gloob cmd/gloob/main.go

# (Optional) Move to PATH
sudo mv gloob /usr/local/bin/
```

### Usage

**Run a Gloob file:**
```bash
gloob yourfile.gloob
```

**Start the interactive REPL:**
```bash
gloob
```

*Note: If you built from source without moving to PATH, use `./gloob` instead.*

### VS Code Extension

Get syntax highlighting for `.gloob` files in Visual Studio Code:

**Install from Marketplace:**
1. Open VS Code
2. Go to Extensions (`Cmd+Shift+X` / `Ctrl+Shift+X`)
3. Search for "Gloob Syntax Language Support"
4. Click Install

Or [install directly from the VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=ChristianDC13.gloob-language).

## 👋 Hello World

Create a file `hello.gloob`:

```js
println("Hello, World!")
```

Run it:
```bash
gloob hello.gloob
```

## 📝 Language Basics

### Variables and Functions
```js
var name = "Gloob"
const PI = 3.14159

fun greet(person) {
    "Hello, " + person + "!"  // Implicit return
}

println(greet(name))
```

### Arrays (1-based!)
```js
var fruits = ["apple", "banana", "cherry"]
println(fruits[1])  // "apple" (first element!)

fruits.push("date").reverse()  // Method chaining!
```

### Loops and Conditionals
```js
loop i from 1 to 5 {
    if i > 3 {
        println("Big number: " + string(i))
    }
}
```

## 📚 Learn More

- **[Complete Language Specification](SPECIFICATION.md)** - Full syntax reference, built-in functions, and detailed feature explanations
- **[Examples Directory](examples/)** - Sample programs including:
  - [Number guessing game](examples/guess-number.gloob)
  - [Module system demo](examples/modules/)
  - [Interactive quiz](examples/quiz.gloob)

## 🧮 Built-in Functions & Methods

Gloob comes with built-in functions for math, I/O, type conversion, and more:

```js
// Math & Utilities
abs(-5), round(3.7), max(1, 5, 3), random(), sleep(1)

// I/O
println("Hello"), var name = input("Name: ")

// String methods (chainable!)
"hello".upper().replace("H", "J")  // "JELLO"

// Array methods (chainable!)
arr.push(10).push(20).reverse()
```

## 🏗️ Architecture

```
Source Code → Lexer → Tokens → Parser → AST → Interpreter → Runtime Values
```

- **Lexer** (`internal/lexer/`) - Tokenizes source code
- **Parser** (`internal/parser/`) - Builds Abstract Syntax Tree
- **Interpreter** (`internal/interpreter/`) - Evaluates AST nodes
- **Scope** (`internal/scope/`) - Manages variables and functions
- **Built-ins** (`internal/builtins/`) - Native functions and methods

## 🤝 Contributing

Contributions are welcome! To add new features:

1. Update the lexer for new tokens (if needed)
2. Add AST nodes for new constructs
3. Implement parsing logic
4. Add evaluation/runtime logic
5. Update documentation

## 📄 License

MIT License - See LICENSE file for details

---

**Made with 🫧 Christian de la Cruz** | [GitHub](https://github.com/ChristianDC13/gloob) | [Report Issues](https://github.com/ChristianDC13/gloob/issues)
