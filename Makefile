PACKAGES = {builtin,lexer,list,math,parser,repl,runtime,strings,util}

all:
	@go install github.com/raoulvdberge/risp

fmt:
	go fmt github.com/raoulvdberge/risp/${PACKAGES}
