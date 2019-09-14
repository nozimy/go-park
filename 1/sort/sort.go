package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	columnsSeparator = " "
	columnIndexError = "column index error"
)

type sortConfig struct {
	desc            bool
	caseInsensitive bool
	unique          bool
	valuesAreNumber bool
	column          int
	outputFilename  string
	inputFilename   string
}

func NewSortConfig() sortConfig {
	return sortConfig{false, false, false, false, 0, "", ""}
}

func main() {
	config := getParsedArgs()

	file, err := os.Open(config.inputFilename)
	if err != nil {
		log.Fatalf("failed to os.Open() %s", err)
	}
	output := os.Stdout
	if len(config.outputFilename) != 0 {
		output, err = os.Create(config.outputFilename)
		if err != nil {
			log.Fatalf("failed to os.Create() %s", err)
		}
	}

	table, err := getTableFromReader(file)
	if err != nil {
		log.Fatalf("failed to getTableFromReader() %s", err)
	}
	err = file.Close()
	if err != nil {
		log.Fatalf("failed to close file %s", err)
	}

	table, err = mySort(table, config)
	if err != nil {
		log.Fatalf("failed to mySort() %s", err)
	}

	err = writeResult(output, table)
	if err != nil {
		log.Fatalf("failed to writeResult() %s", err)
	}
}

func mySort(table [][]string, config sortConfig) ([][]string, error) {
	if table == nil || len(table[0]) < config.column {
		return table, fmt.Errorf(columnIndexError)
	}

	sort.Slice(table, func(i, j int) bool {
		a := table[i][config.column]
		b := table[j][config.column]

		if config.caseInsensitive && !config.valuesAreNumber {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}

		if config.valuesAreNumber {
			aNum, _ := strconv.Atoi(a)
			bNum, _ := strconv.Atoi(b)

			return numberComparator(aNum, bNum, config.desc)
		}

		return stringComparator(a, b, config.desc)
	})

	if config.unique {
		table = getUnique(table, config.column, config.caseInsensitive)
	}

	return table, nil
}

func getParsedArgs() sortConfig {
	config := NewSortConfig()
	flag.BoolVar(&config.desc, "r", false, "-r - сортировка по убыванию")
	flag.BoolVar(&config.caseInsensitive, "f", false, "-f - игнорировать регистр букв")
	flag.BoolVar(&config.unique, "u", false, "-u - выводить только первое среди нескольких равных")
	flag.BoolVar(&config.valuesAreNumber, "n", false, "-n - сортировка чисел")
	flag.IntVar(&config.column, "k", 0, "-k <номер столбца> - сортировать по столбцу (разделитель столбцов по умолчанию можно оставить пробел)")
	flag.StringVar(&config.outputFilename, "o", "", "-o <файл> - выводить в файл, без этой опции выводить в stdout")
	flag.Parse()
	config.inputFilename = flag.Arg(0)

	return config
}

func getTableFromReader(input io.Reader) ([][]string, error) {
	var table [][]string

	if input == nil {
		return table, fmt.Errorf("io argument error")
	}

	in := bufio.NewScanner(input)

	for in.Scan() {
		line := in.Text()

		if len(line) > 0 {
			table = append(table, strings.Split(line, columnsSeparator))
		}
	}

	if in.Err() != nil {
		return table, in.Err()
	}

	return table, nil
}

func writeResult(output io.Writer, table [][]string) error {
	if output == nil {
		return fmt.Errorf("io argument error")
	}

	for _, row := range table {
		_, err := fmt.Fprintln(output, strings.Join(row, columnsSeparator))

		if err != nil {
			return err
		}
	}

	return nil
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
