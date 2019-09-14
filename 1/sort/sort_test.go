package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

var testOkInputString = `12
12
21
45
2
65
6
2
32
`

var testOkInput = [][]string{
	{"12"},
	{"12"},
	{"21"},
	{"45"},
	{"2"},
	{"65"},
	{"6"},
	{"2"},
	{"32"},
}

var testOkResult = [][]string{
	{"65"},
	{"45"},
	{"32"},
	{"21"},
	{"12"},
	{"6"},
	{"2"},
}

var testOkInput2 = [][]string{
	{"Napkin"},
	{"Apple"},
	{"January"},
	{"BOOK"},
	{"January"},
	{"Hauptbahnhof"},
	{"Book"},
	{"Go"},
}

var testOkResult2 = [][]string{
	{"Napkin"},
	{"January"},
	{"Hauptbahnhof"},
	{"Go"},
	{"BOOK"},
	{"Apple"},
}

func TestOkReadInput(t *testing.T) {
	in := bytes.NewBufferString(testOkInputString)
	table, err := getTableFromReader(in)
	require.Equal(t, nil, err)
	require.Equal(t, testOkInput, table)
}

func TestOkWriteResult(t *testing.T) {
	out := bytes.NewBuffer(nil)
	err := writeResult(out, testOkInput)
	require.Equal(t, nil, err)
	require.Equal(t, testOkInputString, out.String())
}

func TestOk(t *testing.T) {
	config := sortConfig{
		desc:            true,
		caseInsensitive: true,
		unique:          true,
		valuesAreNumber: true,
		column:          0,
		outputFilename:  "output.txt",
		inputFilename:   "dataNums.tx",
	}
	table, err := mySort(testOkInput, config)
	require.Equal(t, nil, err)
	require.Equal(t, testOkResult, table)
}

func TestOk2(t *testing.T) {
	config := sortConfig{
		desc:            true,
		caseInsensitive: true,
		unique:          true,
		valuesAreNumber: false,
		column:          0,
		outputFilename:  "output.txt",
		inputFilename:   "data.tx",
	}
	table, err := mySort(testOkInput2, config)
	require.Equal(t, nil, err)
	require.Equal(t, testOkResult2, table)
}

func TestFail(t *testing.T) {
	_, err := mySort(nil, sortConfig{})
	require.EqualError(t, err, columnIndexError)
}

func TestFail2(t *testing.T) {
	config := sortConfig{
		column: 2,
	}
	_, err := mySort(testOkInput, config)
	require.EqualError(t, err, columnIndexError)
}
