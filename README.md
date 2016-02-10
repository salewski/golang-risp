# risp
risp or *Raoul's Lisp* is a Lisp made in Go.

It's feature set is currently very limited, more functions are being added.

## Usage
### Running a file
```
risp examples/hello-world.rp
```
### Running the REPL
```
risp -repl
```
### Running from stdin
```
echo "(println (+ 1 1))" | risp
```

## Building
Make sure you have Go installed and set up correctly.
```
make
```

## Packages
- `lexer` is the package that takes care of the  [lexical analysis](https://en.wikipedia.org/wiki/Lexical_analysis);
- `parser` transforms the tokens provided by `lexer` to an [abstract syntax tree](https://en.wikipedia.org/wiki/Abstract_syntax_tree) or AST;
- `runtime` is the execution engine of the language. It invokes the AST nodes provided by the `parser`;
- `std` is the standard library of the language. It provides basic values like `t`, `f`, `nil` and of course a bunch of useful functions;
- `util` contains utilities for ASCII colors and file reading.

## License
MIT license