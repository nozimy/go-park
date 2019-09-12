package main

import (
	"bytes"
	"io"
	"testing"
)

var testOkInput = `12
12
21
45
2
65
6
2
32`

var testOkResult = `65
45
32
21
12
6
2
`

var testOkInput2 = `Napkin
Apple
January
BOOK
January
Hauptbahnhof
Book
Go`

var testOkResult2 = `Napkin
January
Hauptbahnhof
Go
BOOK
Apple
`

func TestOk(t *testing.T) {
	in := bytes.NewBufferString(testOkInput)
	out := bytes.NewBuffer(nil)
	args := map[string]string{
		"filename": "dataNums.txt",
		"-o":       "output.txt",
		"-f":       "-f",
		"-u":       "-u",
		"-k":       "0",
		"-r":       "-r",
		"-n":       "-n",
	}
	err := mySort(in, out, args)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != testOkResult {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestOk2(t *testing.T) {
	in := bytes.NewBufferString(testOkInput2)
	out := bytes.NewBuffer(nil)
	args := map[string]string{
		"filename": "dataNums.txt",
		"-o":       "output.txt",
		"-f":       "-f",
		"-u":       "-u",
		"-k":       "0",
		"-r":       "-r",
	}
	err := mySort(in, out, args)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != testOkResult2 {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestFail(t *testing.T) {
	in := bytes.NewBufferString(testOkInput2)
	out := bytes.NewBuffer(nil)
	args := map[string]string{
		"-k": "a2",
	}
	err := mySort(in, out, args)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}

func TestFail2(t *testing.T) {
	in := bytes.NewBufferString(testOkInput2)
	out := bytes.NewBuffer(nil)
	args := map[string]string{
		"-k": "2",
	}
	err := mySort(in, out, args)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}

func TestFail3(t *testing.T) {
	var in io.Reader
	var out io.Writer
	args := map[string]string{}
	err := mySort(in, out, args)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}
