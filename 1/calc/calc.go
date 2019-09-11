package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const numbers = "1234567890."
const operators = "+-*/"
const plus = '+'
const minus = '-'
const multiply = '*'
const divide = '/'
const parenOpen = '('
const parenClose = ')'
const endOfLine = '\r'
const malformedExpressionErr = "malformed expression"
const invalidBracketsErr = "invalid brackets"

type StackItem struct {
	operator rune
	num      float64
}

type Stack struct {
	items []StackItem
}

func (s *Stack) New() *Stack {
	s.items = []StackItem{}
	return s
}

func (s *Stack) Push(t StackItem) {
	s.items = append(s.items, t)
}

func (s *Stack) Pop() *StackItem {
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	return &item
}

func (s *Stack) IsEmpty() bool {
	return s.items == nil || len(s.items) == 0
}

func (s *Stack) Peek() *StackItem {
	item := s.items[len(s.items)-1]
	return &item
}

func main() {
	in := bytes.NewBufferString(strings.Join(os.Args[1:], ""))
	err := calc(in, os.Stdout)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func calc(input io.Reader, output io.Writer) error {
	in := bufio.NewScanner(input)
	var line string

	for in.Scan() {
		line = in.Text() + "\r"
	}

	if !validateBrackets(line) {
		return fmt.Errorf(invalidBracketsErr)
	}

	var stack Stack
	stack.New()

	var runeNum []rune
	var num float64
	var lastOperator rune
	var lastCh string

	for _, ch := range line {
		strCh := string(ch)

		if charIsOperator(strCh) && charIsOperator(lastCh) {
			return fmt.Errorf(malformedExpressionErr)
		}

		if charIsNumber(strCh) {
			runeNum = append(runeNum, ch)
		} else if charIsOperator(strCh) {
			num, _ = strconv.ParseFloat(string(runeNum), 64)
			runeNum = nil

			if getPrecedence(ch) >= getPrecedence(lastOperator) {
				if ch == minus {
					runeNum = append(runeNum, ch)
					ch = plus
				}

				stack.Push(StackItem{
					operator: ch,
					num:      num,
				})
			} else {
				for !stack.IsEmpty() {
					item := stack.Pop()
					num, _ = performOperator(item.num, num, item.operator)
				}

				stack.Push(StackItem{
					operator: ch,
					num:      num,
				})
			}

			lastOperator = ch
		} else if ch == parenOpen {
			stack.Push(StackItem{
				operator: ch,
				num:      0,
			})

			lastOperator = ch
		} else if ch == parenClose {
			num, _ = strconv.ParseFloat(string(runeNum), 64)

			for stack.Peek().operator != parenOpen {
				item := stack.Pop()
				num, _ = performOperator(item.num, num, item.operator)
			}

			stack.Pop()
			lastOperator = ch
			runeNum = []rune(fmt.Sprintf("%f", num))
		} else if ch == endOfLine {
			if charIsOperator(lastCh) {
				return fmt.Errorf(malformedExpressionErr)
			}

			num, _ = strconv.ParseFloat(string(runeNum), 64)

			for !stack.IsEmpty() {
				item := stack.Pop()
				num, _ = performOperator(item.num, num, item.operator)
			}
		} else {
			return fmt.Errorf("invalid character %s", strCh)
		}
		lastCh = strCh
	}

	_, err := fmt.Fprint(output, num)

	return err
}

func charIsNumber(ch string) bool {
	return strings.Index(numbers, ch) != -1
}

func charIsOperator(ch string) bool {
	return strings.Index(operators, ch) != -1 && ch != ""
}

func performOperator(a float64, b float64, operator rune) (float64, error) {
	switch operator {
	case plus:
		return a + b, nil
	case minus:
		return a - b, nil
	case multiply:
		return a * b, nil
	case divide:
		return a / b, nil
	default:
		return 0, errors.New("perform operator error")
	}
}

func validateBrackets(str string) bool {
	count := 0

	for _, code := range str {
		ch := string(code)
		if ch == "(" {
			count++
		} else if ch == ")" {
			if count == 0 {
				return false
			} else {
				count--
			}
		}
	}

	if count == 0 {
		return true
	}

	return false
}

func getPrecedence(operator rune) int {
	switch operator {
	case plus:
		fallthrough
	case minus:
		return 1
	case multiply:
		fallthrough
	case divide:
		return 2
	}

	return -1
}
