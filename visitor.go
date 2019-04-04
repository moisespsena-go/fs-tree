package fstree

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/moisespsena-go/os-common"

	"github.com/moisespsena-go/error-wrap"
)

type Visitor struct {
	Root    *Node
	Visited map[string]*Node
	Options *Options
}

func (v *Visitor) on(n *Node) bool {
	if _, ok := v.Visited[n.Path]; ok {
		return false
	}
	v.Visited[n.AbsPath] = n
	return true
}

func (v *Visitor) Visit() (err error) {
	if v.Options == nil {
		v.Options = &Options{}
	}

	if v.Options.PreVisitValid == nil {
		v.Options.PreVisitValid = AllValid
	}
	if v.Options.PostVisitValid == nil {
		v.Options.PostVisitValid = AllValid
	}
	if v.Visited == nil {
		v.Visited = map[string]*Node{}
	}

	root := v.Root
	if v.Options.Prefix != "" {
		prefix := path.Clean(strings.Replace(v.Options.Prefix, "\\", "/", -1))
		if prefix[0] == '/' {
			prefix = prefix[1:]
		}
		parts := strings.Split(prefix, "/")
		for _, name := range parts {
			if name != "" {
				n := &Node{FileInfo: oscommon.NewVirtualDirFileInfo(name)}
				root.children[name] = n
				root = n
			}
		}
	}

	return v.visit(root)
}

func (v *Visitor) visit(n *Node) (err error) {
	if !v.on(n) {
		return
	}

	n.children = map[string]*Node{}

	names, err := readDirNames(n.Path)
	if err != nil {
		return errwrap.Wrap(err, n.Path)
	}

	for _, name := range names {
		c := &Node{Path: filepath.Join(n.Path, name)}
		c.FileInfo, err = os.Lstat(c.Path)
		if err != nil {
			return errwrap.Wrap(err, c.Path)
		}
		if c.AbsPath, err = filepath.Abs(c.Path); err != nil {
			return err
		}
		if !v.Options.PreVisitValid(c) {
			continue
		}
		if c.IsDir() {
			if !v.Options.NotRecursive {
				if err = v.visit(c); err != nil {
					if err == filepath.SkipDir {
						continue
					} else {
						return err
					}
				}
			}
		} else if isSymlink(c.Mode()) {
			if err = v.followSymlink(c); err != nil {
				return errwrap.Wrap(err, n.Path)
			}
		}

		if v.Options.PostVisitValid(c) {
			n.children[name] = c
		}
	}
	return
}

func (v *Visitor) followSymlink(n *Node) (err error) {
	var (
		ok   bool
		link Link
		info os.FileInfo = n
		n2   *Node
	)

	for isSymlink(info.Mode()) {
		if link.Path, err = os.Readlink(n.AbsPath); err != nil {
			return
		}

		if !filepath.IsAbs(link.Path) {
			if link.AbsPath, err = filepath.Abs(filepath.Join(filepath.Dir(n.Path), link.Path)); err != nil {
				return
			}
		}

		if n2, ok = v.Visited[link.AbsPath]; ok {
			n.Link = n2
			return
		}

		if info, err = os.Lstat(link.AbsPath); err != nil {
			return errwrap.Wrap(err, link.AbsPath)
		}
	}

	n2 = &Node{FileInfo: info, Path: link.Path, AbsPath: link.AbsPath}

	if info.IsDir() {
		s2 := &Visitor{n2, v.Visited, v.Options}
		if !v.Options.NotRecursive {
			if err = s2.visit(n2); err != nil {
				return errwrap.Wrap(err, "Visit %q", link.AbsPath)
			}
		}
	}

	n.Link = n2
	n.FileInfo = &SymlinkInfo{n.FileInfo, n2.FileInfo}
	return
}
