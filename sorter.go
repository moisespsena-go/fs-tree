package fstree

import "sort"

type Sorter func(children Children) (names []string, err error)

func NameSorter(children Children) (names []string, err error) {
	for name := range children {
		names = append(names, name)
	}
	sort.Strings(names)
	return
}
