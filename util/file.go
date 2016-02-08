package util

import (
	"io/ioutil"
	"path"
	"strings"
)

type File struct {
	Path string
	Name string
	Data string
}

func NewFile(filePath string) (*File, error) {
	i, j := strings.LastIndex(filePath, "/")+1, strings.LastIndex(filePath, path.Ext(filePath))

	file := &File{Name: filePath[i:j], Path: filePath}

	data, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return nil, err
	}

	file.Data = string(data)

	return file, nil
}
