package main

import "github.com/moisespsena-go/fs-tree"

func main() {
	t, err := fstree.New("sample/data")
	if err != nil {
		panic(err)
	}
	println(t.Visit())
	println(t.Root)
}
