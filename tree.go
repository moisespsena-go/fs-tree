package fstree

import (
	"os"
	"path/filepath"

	"github.com/moisespsena-go/error-wrap"

	"github.com/go-errors/errors"
)

type Options struct {
	PreVisitValid  ValidFunc
	PostVisitValid ValidFunc
	NotRecursive   bool
	Prefix         string
}

func New(dir string) (visitor *Visitor, err error) {
	root := &Node{Path: dir}
	root.FileInfo, err = os.Stat(dir)
	if err != nil {
		return nil, errwrap.Wrap(err, dir)
	}
	if !root.IsDir() {
		return nil, errors.New("not is dir")
	}
	if root.AbsPath, err = filepath.Abs(dir); err != nil {
		return nil, err
	}
	return &Visitor{root, map[string]*Node{}, nil}, nil
}
