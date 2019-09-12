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

const paramOutputFile = "-o"
const paramSortColumn = "-k"
const paramFilename = "filename"
const paramDescSort = "-r"
const paramCaseInsensitive = "-f"
const paramSortNumberValueType = "-n"
const paramUniqueFilter = "-u"
const columnsSeparator = " "
const emptyString = ""
const stringTypeName = "string"
const intTypeName = "int"

type sortConfig struct {
	desc            bool
	caseInsensitive bool
	unique          bool
	valuesType      string
	column          int
}

func main() {
	args := getParsedArgs()

	file, _ := os.Open(args[paramFilename])
	output := os.Stdout

	if args[paramOutputFile] != emptyString {
		output, _ = os.Create(args[paramOutputFile])
	}

	err := mySort(file, output, args)
	closeErr := file.Close()

	if err != nil && closeErr != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mySort(input io.Reader, output io.Writer, args map[string]string) error {
	if input == nil || output == nil {
		return fmt.Errorf("io argument error")
	}

	config, e := getSortConfig(args)

	if e != nil {
		return fmt.Errorf(e.Error())
	}

	table := getTableFromReader(input)

	if table != nil && len(table[0]) < config.column {
		return fmt.Errorf("column index error")
	}

	sort.Slice(table, func(i, j int) bool {
		a := table[i][config.column]
		b := table[j][config.column]

		if config.caseInsensitive && config.valuesType != intTypeName {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

		if config.valuesType == intTypeName {
			aNum, _ := strconv.Atoi(a)
			bNum, _ := strconv.Atoi(b)

			return numberComparator(aNum, bNum, config.desc)
		}

		return stringComparator(a, b, config.desc)
	})

	if config.unique {
		table = getUnique(table, config.column, config.caseInsensitive)
	}

	for _, row := range table {
		_, e := fmt.Fprintln(output, strings.Join(row, columnsSeparator))

		if e != nil {
			return e
		}
	}

	return nil
}

func getParsedArgs() map[string]string {
	args := make(map[string]string)
	var prevArg string
	var key, val string

	for _, arg := range os.Args[1:] {
		if arg == paramOutputFile || arg == paramSortColumn {
			key = arg
			val = emptyString
		}

		if prevArg == paramOutputFile || prevArg == paramSortColumn {
			key = prevArg
			val = arg
		}

		if strings.HasSuffix(arg, ".txt") {
			key = paramFilename
			val = arg
		}

		args[key] = val
		prevArg = arg
	}

	return args
}

func getTableFromReader(input io.Reader) [][]string {
	in := bufio.NewScanner(input)
	var table [][]string

	for in.Scan() {
		line := in.Text()

		if len(line) > 0 {
			var row []string

			for _, word := range strings.Split(line, columnsSeparator) {
				row = append(row, word)
			}

			table = append(table, row)
		}
	}

	return table
}

func numberComparator(a, b int, desc bool) bool {
	if desc {
		return a > b
	}

	return a < b
}

func stringComparator(a, b string, desc bool) bool {
	if desc {
		return a > b
	}

	return a < b
}

func getUnique(table [][]string, col int, caseInsensitive bool) [][]string {
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

	return table
}

func getSortConfig(args map[string]string) (sortConfig, error) {
	var config = sortConfig{false, false, false, stringTypeName, 0}

	var e error

	if args[paramSortColumn] != emptyString {
		config.column, e = strconv.Atoi(args[paramSortColumn])

		if e != nil {
			return config, e
		}
	}

	if args[paramDescSort] != emptyString {
		config.desc = true
	}

	if args[paramCaseInsensitive] != emptyString {
		config.caseInsensitive = true
	}

	if args[paramSortNumberValueType] != emptyString {
		config.valuesType = intTypeName
	}

	if args[paramUniqueFilter] != emptyString {
		config.unique = true
	}

	return config, nil
}
