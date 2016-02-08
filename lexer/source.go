package lexer

import "github.com/raoulvdberge/risp/util"

type Source interface {
	Name() string
	Data() string
}

type FileSource struct {
	file *util.File
}

func NewSourceFromFile(file *util.File) *FileSource {
	return &FileSource{file: file}
}

func (f *FileSource) Name() string {
	return f.file.Name
}

func (f *FileSource) Data() string {
	return f.file.Data
}

type StringSource struct {
	name string
	data string
}

func NewSourceFromString(name string, data string) *StringSource {
	return &StringSource{name: name, data: data}
}

func (f *StringSource) Name() string {
	return f.name
}

func (f *StringSource) Data() string {
	return f.data
}
