package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	numbers                = "1234567890."
	operators              = "+-*/"
	plus                   = '+'
	minus                  = '-'
	multiply               = '*'
	divide                 = '/'
	parenOpen              = '('
	parenClose             = ')'
	endOfLine              = '\r'
	malformedExpressionErr = "malformed expression"
	invalidBracketsErr     = "invalid brackets"
)

type StackItem struct {
	operator rune
	num      float64
}

type Stack struct {
	items []StackItem
}

func NewStack() Stack {
	return Stack{}
}

func (s *Stack) PushItem(t StackItem) {
	s.items = append(s.items, t)
}

func (s *Stack) Push(num float64, operator rune) {
	s.PushItem(StackItem{
		operator: operator,
		num:      num,
	})
}

func (s *Stack) Pop() *StackItem {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return &item
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Peek() *StackItem {
	item := s.items[len(s.items)-1]
	return &item
}

func main() {
	in := bytes.NewBufferString(strings.Join(os.Args[1:], ""))

	var err error
	line, err := readExpression(in)
	if err != nil {
		log.Fatalf("failed to readExpression() %s", err)
	}
	num, err := calc(line)
	if err != nil {
		log.Fatalf("failed to calc() %s", err)
	}
	err = writeResult(os.Stdout, num)
	if err != nil {
		log.Fatalf("failed to writeResult() %s", err)
	}
}

func readExpression(input io.Reader) (string, error) {
	in := bufio.NewScanner(input)
	var line string

	for in.Scan() {
		line = in.Text() + "\r"
	}

	if in.Err() != nil {
		return line, in.Err()
	}

	return line, nil
}

func writeResult(output io.Writer, num float64) error {
	_, err := fmt.Fprint(output, num)

	return err
}

func calc(line string) (float64, error) {
	var num float64

	if !validateBrackets(line) {
		return num, fmt.Errorf(invalidBracketsErr)
	}

	stack := NewStack()
	var runeNum []rune
	var lastOperator rune
	var lastChar string

	for _, char := range line {
		strChar := string(char)

		if charIsOperator(strChar) && charIsOperator(lastChar) {
			return num, fmt.Errorf(malformedExpressionErr)
		}

		var err error

		switch {
		case charIsNumber(strChar):
			runeNum = append(runeNum, char)
		case charIsOperator(strChar):
			if runeNum != nil {
				num, err = strconv.ParseFloat(string(runeNum), 64)
				if err != nil {
					return num, err
				}
				runeNum = nil
			}

			if getPrecedence(char) >= getPrecedence(lastOperator) {
				if char == minus {
					runeNum = append(runeNum, char)
					char = plus
				}

				stack.Push(num, char)
			} else {
				num, err = pullStack(&stack, num)
				if err != nil {
					return num, err
				}

				stack.Push(num, char)
			}

			lastOperator = char
		case char == parenOpen:
			stack.Push(0, char)

			lastOperator = char
		case char == parenClose:
			num, err = strconv.ParseFloat(string(runeNum), 64)
			if err != nil {
				return num, err
			}
			num, err = pullStack(&stack, num)
			if err != nil {
				return num, err
			}

			stack.Pop()
			lastOperator = char
			runeNum = []rune(fmt.Sprintf("%f", num))
		case char == endOfLine:
			if charIsOperator(lastChar) {
				return num, fmt.Errorf(malformedExpressionErr)
			}

			num, err = strconv.ParseFloat(string(runeNum), 64)
			if err != nil {
				return num, err
			}
			num, err = pullStack(&stack, num)
			if err != nil {
				return num, err
			}
		default:
			return num, fmt.Errorf("invalid character %s", strChar)
		}

		lastChar = strChar
	}

	return num, nil
}

func charIsNumber(ch string) bool {
	return strings.Index(numbers, ch) != -1 && ch != ""
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
	case plus, minus:
		return 1
	case multiply, divide:
		return 2
	}

	return -1
}

func pullStack(stack *Stack, num float64) (float64, error) {
	var err error
	for !stack.IsEmpty() && stack.Peek().operator != parenOpen {
		item := stack.Pop()
		num, err = performOperator(item.num, num, item.operator)
		if err != nil {
			return num, err
		}
	}

	return num, nil
}
