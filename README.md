# glox

## What is this?
A Golang implementation of the 'jlox' interpreter from the book [Crafting Interpreters](https://craftinginterpreters.com/) by Bob Nystrom.

## How can I use it?
If you already installed the Go compiler in your computer, you can start using 'glox' following these steps:

Clone this repository:
```
git clone https://github.com/jcbages/glox
```

Navigate into the 'glox' folder:
```
cd glox
```

Run 'glox' using the 'go' command. You can either pass a filename you want to run or just open the REPL: 
```
go run *.go [file]
```

## What can I do with this?
So far we support the following:

### Arithmetic expressions
Use operators like `*`, `/`, `+`, and `-` to perform calculations
```
1 + (12 / 3) * 43 - 10
```

### Variables
You can store a piece of data to recall later using variables
```
var a = 10
var b = a + 10 // b = 20
...
a = 12
...
var c = a + 10 // c = 22
```

### Print
You can print the value of variables and/or the result of an expression
```
var name = "Mister Glox"
print "Hello, " + name
```
