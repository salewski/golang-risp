PACKAGES = {lexer,parser,runtime,std,util}

all:
	@go install github.com/raoulvdberge/risp

fmt:
	go fmt github.com/raoulvdberge/risp/${PACKAGES}
