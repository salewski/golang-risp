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

## License
MIT license