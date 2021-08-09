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
var a = 10;
var b = a + 10; // b = 20
...
a = 12;
...
var c = a + 10; // c = 22
```

### Print
You can print the value of variables and/or the result of an expression
```
var name = "Mister Glox"
print "Hello, " + name;
```

### Conditionals
You can execute a piece of code if and only if a certain condition is true
```
if (10 > 2) {
    print "Yes 10 is larger than 2 as it should be";
} else {
    print "This is not what I expected";
}
```

### While loops
You can execute a piece of code multiple times using a while loop
```
var i = 0;
while (i < 10) {
    print "Value of i => " + i;
    i = i + 1;
}
```

### For loops
Similar to while loops, for loops will execute a piece of code multiple times but they 
have a shorter "syntactic sugar" shape
```
for (var i = 0; i < 10; i = i + 1) {
    print "Value of i => " + i;
}
```

### Functions
You can reuse a piece of code by creating a function
```
fun fibonacci(n) {
    if (n <= 1) {
        return 1;
    } else {
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
}
...
print fibonacci(5); // 8
```

### Closures
You can make functions that return other functions using closures
```
fun makeCounter() {
    var i = 0;

    fun count() {
        i = i + 1;
        print "Value of i => " + i;
    }
    return count
}
...
counter = makeCounter();
counter(); // 1
counter(); // 2
```

## Sample code
This is a sample of a valid program that can be currently executed with the 'glox' interpreter:
```
var a = "global a";
var b = "global b";
var c = "global c";

{
    var a = "outer a";
    var b = "outer b";
    {
        var a = "inner a";
        print a;
        print b;
        print c;
        print a + ", " + b + ", " + c;
    }
    print a;
    print b;
    print c;
    print a + ", " + b + ", " + c;
}
print a;
print b;
print c;
print a + ", " + b + ", " + c;

print (12 + 23) * 24 / 12 + 1;
print "hello my name is jeff";

if (10 > 0) {
    print clock();
} else {
    print "nope";
}

fun fib(n) {
    if (n <= 1) {
        return 1;
    } else {
        return fib(n-1) + fib(n-2);
    }
}

var start = clock();
print "START" + " " + start;
for (var i = 0; i < 30; i = i+1) {
    print "fib " + i + " = " + fib(i);
}
var end = clock();
print "END" + " " + end;
var total = (end - start);
print "total execution time = " + total + " seconds";

fun makeCounter() {
  var i = 0;
  fun count() {
    i = i + 1;
    print "VALUE OF I -> " + i;
  }

  return count;
}

var counter = makeCounter();
counter(); // "1".
counter(); // "2".
```
