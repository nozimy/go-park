package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type ByCol [][]string

func (a ByCol) Len() int           { return len(a) }
func (a ByCol) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCol) Less(i, j int) bool { return a[i][0] < a[j][0] }

func main() {
	var table [][]string
	filename := os.Args[1]
	fmt.Println("filename " + filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) > 0 {
			var row []string
			for _, word := range strings.Split(line, " ") {
				row = append(row, word)
			}
			table = append(table, row)
		}
	}
	sort.Sort(ByCol(table))
	for _, row := range table {
		fmt.Println(strings.Join(row, " "))
	}
}
