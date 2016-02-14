PACKAGES = {builtin,lexer,lists,math,parser,runtime,strings,util}

all:
	@go install github.com/raoulvdberge/risp

fmt:
	go fmt github.com/raoulvdberge/risp/${PACKAGES}
