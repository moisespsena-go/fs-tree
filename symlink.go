package fstree

import (
	"os"
	"time"
)

type SymlinkInfo struct {
	Link os.FileInfo
	Real os.FileInfo
}

func (si SymlinkInfo) Name() string {
	return si.Link.Name()
}

func (si SymlinkInfo) Size() int64 {
	return si.Real.Size()
}

func (si SymlinkInfo) Mode() os.FileMode {
	return si.Real.Mode()
}

func (si SymlinkInfo) ModTime() time.Time {
	return si.Real.ModTime()
}

func (si SymlinkInfo) IsDir() bool {
	return si.Real.IsDir()
}

func (si SymlinkInfo) Sys() interface{} {
	return si.Real.Sys()
}
