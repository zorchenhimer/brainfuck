# Brainfuck Interpreter

A simple interpreter for the [brainfuck](https://en.wikipedia.org/wiki/Brainfuck)
language written in Go. 

## Running

```
$ cat hello-world.bf
++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.
$ go build
$ brainfuck hello-world.bf
Hello World!
```

## Input caveat

Inputting data (with `,`) is not very user friendly yet.  Input is buffered so
the user has to hit the enter key to send the data to the program.  Fixing
this on Linux is easy enough but Windows is a different story. 

