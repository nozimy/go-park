package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	args := make(map[string]string)
	var prevArg string

	for _, arg := range os.Args[1:] {
		if prevArg == "-o" || prevArg == "-k" {
			args[prevArg] = arg
			prevArg = arg
			continue
		}

		if strings.HasSuffix(arg, ".txt") {
			args["filename"] = arg
			prevArg = arg
			continue
		}

		args[arg] = arg
		prevArg = arg
	}

	file, _ := os.Open(args["filename"])
	output := os.Stdout

	if args["-o"] != "" {
		output, _ = os.Create(args["-o"])
	}

	err := mySort(file, output, args)
	closeErr := file.Close()

	if err != nil && closeErr != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mySort(input io.Reader, output io.Writer, args map[string]string) error {
	in := bufio.NewScanner(input)
	var table [][]string

	for in.Scan() {
		line := in.Text()

		if len(line) > 0 {
			var row []string

			for _, word := range strings.Split(line, " ") {
				row = append(row, word)
			}

			table = append(table, row)
		}
	}

	desc := false
	caseInsensitive := false
	valuesType := "string"
	unique := false
	col := 0
	var e error

	if args["-k"] != "" {
		col, e = strconv.Atoi(args["-k"])

		if e != nil {
			return fmt.Errorf(e.Error())
		}
	}

	if args["-r"] != "" {
		desc = true
	}

	if args["-f"] != "" {
		caseInsensitive = true
	}

	if args["-n"] != "" {
		valuesType = "number"
	}

	if args["-u"] != "" {
		unique = true
	}

	if table != nil && len(table[0]) < col {
		return fmt.Errorf("column index error")
	}

	sort.Slice(table, func(i, j int) bool {
		a := table[i][col]
		b := table[j][col]

		if caseInsensitive && valuesType != "number" {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

		var a1 int
		var b1 int

		if valuesType == "number" {
			a1, _ = strconv.Atoi(a)
			b1, _ = strconv.Atoi(b)

			if desc {
				return a1 > b1
			}

			return a1 < b1
		}

		if desc {
			return a > b
		}

		return a < b
	})

	if unique {
		uniques := make(map[string]bool)
		var ln int

		for _, row := range table {
			key := row[col]

			if caseInsensitive {
				key = strings.ToLower(key)
			}

			if uniques[key] {
				continue
			}

			uniques[key] = true
			table[ln] = row
			ln++
		}

		table = table[:ln]
	}

	for _, row := range table {
		_, e := fmt.Fprintln(output, strings.Join(row, " "))

		if e != nil {
			return e
		}
	}

	return nil
}
