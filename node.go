package fstree

import (
	"os"
	"path"

	"github.com/go-errors/errors"
)

type ValidFunc func(n *Node) bool

type Children map[string]*Node

func AllValid(*Node) bool { return true }

type Link struct {
	AbsPath string
	Path    string
}

type Node struct {
	os.FileInfo
	Link     *Node
	Path     string
	AbsPath  string
	children Children
	sorted   []string
	Data     interface{}
}

func (n *Node) Child(name string) (child *Node, ok bool) {
	child, ok = n.children[name]
	return
}

func (n *Node) MustChild(name string) (child *Node) {
	return n.children[name]
}

func (n *Node) Children() Children {
	return n.children
}

func (n *Node) Count() uint {
	return uint(len(n.children))
}

func (n *Node) IsEmpty() bool {
	return len(n.children) == 0
}

func (n *Node) Sorted() []string {
	return n.sorted
}

func (n *Node) walkDir(pth string, cb func(n *Node) error) (err error) {
	if n.children == nil {
		return nil
	}
	for name, child := range n.children {
		if child.IsDir() {
			if err = cb(child); err != nil {
				return
			}
			if err = child.walkDir(path.Join(pth, name), cb); err != nil {
				return
			}
		}
	}
	return
}

func (n *Node) WalkDir(cb func(n *Node) error) (err error) {
	return n.walkDir(".", cb)
}

func (n *Node) walk(pth string, depth uint, cb func(pth string, level, index uint, dir, n *Node) error) (err error) {
	if n.children == nil {
		return nil
	}
	var npth string
	if len(n.sorted) > 0 {
		var child *Node
		for i, name := range n.sorted {
			npth = path.Join(pth, name)
			child = n.children[name]
			if err = cb(npth, depth+1, uint(i), n, child); err != nil {
				return
			}
			if child.IsDir() {
				if err = child.walk(npth, depth+1, cb); err != nil {
					return
				}
			}
		}
	} else {
		for name, child := range n.children {
			npth = path.Join(pth, name)
			if err = cb(npth, depth+1, 0, n, child); err != nil {
				return
			}
			if child.IsDir() {
				if err = child.walk(npth, depth+1, cb); err != nil {
					return
				}
			}
		}
	}
	return
}

func (n *Node) Walk(cb func(pth string, depth, index uint, dir, n *Node) error) error {
	return n.walk(".", 0, cb)
}

func (n *Node) SortR(sorter Sorter) (err error) {
	if n.children == nil {
		return nil
	}

	var sorted []string

	if sorter == nil {
		sorter = NameSorter
	}

	if sorted, err = sorter(n.children); err != nil {
		return
	}

	if len(sorted) != len(n.children) {
		return errors.New("Invalid sorted length")
	}

	n.sorted = sorted

	for _, child := range n.children {
		if child.IsDir() {
			if err = child.SortR(sorter); err != nil {
				return
			}
		}
	}
	return
}

func (n *Node) Sort(sorter Sorter) (err error) {
	if n.children == nil {
		return nil
	}
	var sorted []string

	if sorter == nil {
		sorter = NameSorter
	}

	if sorted, err = sorter(n.children); err != nil {
		return
	}

	if len(sorted) != len(n.children) {
		return errors.New("Invalid sorted length")
	}

	n.sorted = sorted
	return
}

func (n *Node) Each(cb func(n *Node) error) (err error) {
	if n.children == nil {
		return nil
	}
	if len(n.sorted) > 0 {
		for _, name := range n.sorted {
			if err = cb(n.children[name]); err != nil {
				return
			}
		}
	} else {
		for _, child := range n.children {
			if err = cb(child); err != nil {
				return
			}
		}
	}
	return
}

func (n *Node) Merge(root *Node) {
	for name, on := range root.children {
		c, ok := n.children[name]
		if !ok || !on.IsDir() {
			n.children[name] = on
		} else {
			for ocname, oc := range on.children {
				if !oc.IsDir() {
					c.children[ocname] = oc
				} else {
					c.Merge(oc)
				}
			}
		}
	}
}
