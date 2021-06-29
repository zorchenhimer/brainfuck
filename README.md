# Brainfuck Interpreter

A simple interpreter for the [brainfuck](https://en.wikipedia.org/wiki/Brainfuck)
language written in Go.

Dialects supported (& default extension):
- Brainfuck (`.bf`)
- [FuckFuck](https://github.com/MiffOttah/fuckfuck) (`.ff`)
- TenX (`.ten`)
- [Pikalang](https://www.dcode.fr/pikalang-language) (`.pika`)

## Running

```
$ cat examples/hello-world.bf
++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.
$ go build
$ brainfuck examples/hello-world.bf
Hello World!
```

### Arguments

`-d`, `--debug`
Prints out some debugging information.  Very verbose.

`--dialects`
Lists supported dialects.

`-l`, `--language`
Specifiy the dialect to use.

## Input caveat

Inputting data (with `,`) is not very user friendly yet.  Input is buffered so
the user has to hit the enter key to send the data to the program.  Fixing
this on Linux is easy enough but Windows is a different story. 

