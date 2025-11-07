# 🫧Gloob Programming Language
> *Gloob* is a small, playful, interpreted language — made just for fun and experimentation.  
> It's simple, readable, and unpretentious. Think of it as JavaScript's lazy cousin who just wants to print stuff and loop things.

---

## 🖨️ Output
```js
print('Hello', 'World')
println('Single line')
```
Basic output functions.  
Both `print()` and `println()` add a newline automatically and accept multiple arguments.

---

## 🧩 Variables
```js
var x = 5
x = 6
print(x) // 6

const PI = 3.14159
// PI = 3 // Error: cannot reassign constant
```
Variables are defined using `var` and can be reassigned freely.  
Constants are defined using `const` and cannot be reassigned.  
Everything is dynamically typed.

---

## 🔢 Primitive Data Types
```js
var num = 5        // int
var height = 5.4   // float
var name = 'John'  // string
var active = false // bool
var nothing = null // null
```

### Numbers
```js
var positive = 42
var negative = -10
var decimal = 3.14159
```
Supports both integers and floating-point numbers, including negative values.

### Boolean extensions
```js
var active = true | false | yes | no | on | off
```
Because sometimes you just *feel* like typing "yes" instead of "true".

---

## 📦 Collections

### Arrays
```js
arr = [1, 2, 3, 4]
second = arr[2] // 2 (arrays are 1-based!)
arr.push(5)
last = arr.pop()
arr.remove(2)   // Remove element at index 2
arr.insert(1, 99)  // Insert 99 at index 1
size = arr.len()

// New methods
index = arr.indexOf(3)  // 3 (1-based index!)
has = arr.contains(5)  // true
text = arr.join(", ")  // "1, 2, 3, 4"
arr.reverse()  // Reverses in-place

// Method chaining (methods that modify return the array)
arr.push(10).push(20).reverse()
```
Arrays are **1-based indexed** and come with built-in methods:
- `.push(value)` - Add element to end, returns array
- `.pop()` - Remove and return last element
- `.remove(index)` - Remove element at index (1-based), returns array
- `.insert(index, value)` - Insert element at index (1-based), returns array
- `.len()` - Get array length
- `.indexOf(element)` - Find 1-based index (0 if not found)
- `.contains(element)` - Check if array contains element
- `.join(separator)` - Join elements into string
- `.reverse()` - Reverse array in-place, returns array

### Array iteration
```js
loop element from arr {
    println(element)
}
```

### Objects
```js
user = { name: "Jane Doe", age: 25 }
age = user.age
user.active = false
```
Objects store key-value pairs.  
New properties can be added dynamically.

---

## 📝 Strings
```js
greeting = "Hello"
name = 'World'

// String indexing (1-based)
first = greeting[1]  // "H"
last = greeting[5]   // "o"
```
Strings can use single or double quotes.  
Strings are **1-based indexed** like arrays!

### String methods
```js
text = "  Hello World  "
length = text.len()        // 15
upper = text.upper()       // "  HELLO WORLD  "
lower = text.lower()       // "  hello world  "
clean = text.trim()        // "Hello World"
has = text.contains("World")  // true

// New methods
words = "hello world".split(" ")  // ["hello", "world"]
replaced = "hello".replace("ll", "y")  // "heyo"
index = "hello".indexOf("ll")  // 3 (1-based index!)
```
Built-in string methods:
- `.len()` - Get string length
- `.upper()` - Convert to uppercase
- `.lower()` - Convert to lowercase
- `.trim()` - Remove leading/trailing whitespace
- `.contains(substring)` - Check if string contains substring
- `.split(separator)` - Split string into array
- `.replace(old, new)` - Replace all occurrences
- `.indexOf(substring)` - Find 1-based index (0 if not found)

### Method chaining
```js
result = "  hello  ".trim().upper()  // "HELLO"
println("Gloob".len())  // 5 (call methods on literals!)
```

---

## ⚖️ Conditionals
```js
if grade > 75 {
    println('Cool')
} else if grade > 50 {
    println('Not bad')
} else {
    println('Not cool')
}
```
Straightforward branching — no parentheses, just clean blocks.  
Supports `if`, `else if` (or `elseif`), and `else`.

---

## 🔁 Loops
```js
// Range loop
loop i from 1 to 100 {
    println(i)
}

// Range with increment
loop i from 0 to 10: 2 {
    println(i)  // 0, 2, 4, 6, 8, 10
}

// Backward loop with negative increment
loop i from 10 to 0: -2 {
    println(i)  // 10, 8, 6, 4, 2, 0
}

// While-style loop
loop condition {
    // Do stuff
    break  // Break the loop
}

// For-each loop
loop element from arr {
    println(element)
}

// Infinite loop
loop {
    // Do stuff forever
    break  // until you break
}
```
Five loop types — simple, flexible, and readable.  
You can even mix styles if you dare.

**Loop control:**
- `break` - Exit the loop immediately

---

## 🧠 Functions
```js
fun doSomething(param1, param2) {
    // Explicit return
    return param1 + param2
}

fun implicit(x) {
    // Last expression is auto-returned
    x * 2
}

fun earlyExit(n) {
    if n < 0 {
        return  // Bare return stops execution
    }
    n * n
}

doSomething(1, 2)
```
Functions are first-class citizens, declared with `fun`.  
**Return behavior:**
- `return value` - Explicit return with value
- Last expression in function body is automatically returned
- `return` alone stops execution and returns `null`

---

## 🧍 Input
```js
var name = input('What's your name? ')
var age = input('Age: ')
```
Simple console input for interactive programs.  
Returns a string that you can convert with `number()`.

---

## 📚 Modules and Imports
```js
import "utils/helpers"
import "math.gloob"  // .gloob extension is optional
```
Import entire files to reuse code across your project.  
- Paths are relative to the importing file
- Supports recursive/nested imports
- Circular import detection built-in
- Extension (`.gloob`) is optional

---

## 🧮 Built-in Functions

### Math
```js
abs(-5)           // 5
round(3.7)        // 4
max(1, 5, 3)      // 5
min(1, 5, 3)      // 1
random()          // Random float 0-1
randInt(1, 10)    // Random int 1-10
```

### Type conversion
```js
number("42")      // 42
string(123)       // "123"
bool(1)           // true
type(42)          // "NUMERIC"
```

### Utility
```js
len("hello")      // 5 (works with strings and arrays)
len([1, 2, 3])    // 3
sleep(1)          // Sleep for 1 second
clear()           // Clear the terminal screen
```
**Note:** `len()` is available both as a standalone function and as a method (`.len()`). Use whichever feels more natural!

---

## 🎯 Operators

### Arithmetic
```js
+ - * / %
```

### Comparison
```js
== != > >= < <=
```

### Logical
```js
&& ||
```

### Assignment
```js
= 
```

---

## 💡 Examples

### Fibonacci
```js
fun fibonacci(n) {
    if n <= 1 {
        return n
    }
    fibonacci(n - 1) + fibonacci(n - 2)
}

println(fibonacci(10))  // 55
```

### Array operations
```js
numbers = [1, 2, 3, 4, 5]

loop num from numbers {
    println(num * 2)
}

numbers.push(6)
println(numbers.len())  // 6
```

### Method chaining
```js
result = [1, 2, 3].len()  // 3
text = "  hello  ".trim().upper()  // "HELLO"
```

---

## 🌟 Language Quirks

- **1-based arrays**: `arr[1]` gets the first element
- **Implicit returns**: Last expression in a function is automatically returned
- **Boolean alternatives**: Use `yes`/`no` or `on`/`off` instead of `true`/`false`
- **Method chaining**: Call methods directly on literals like `"hi".upper()`
- **No parentheses in conditionals**: `if x > 5 {` not `if (x > 5) {`
- **Dynamic typing**: Variables can hold any type of value

---

**Made with 🫧 for fun and learning**
