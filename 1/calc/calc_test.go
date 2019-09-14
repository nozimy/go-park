package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOkReadExpression(t *testing.T) {
	in := bytes.NewBufferString("2+2*2/2")
	line, err := readExpression(in)
	require.Equal(t, nil, err)
	require.Equal(t, "2+2*2/2\r", line)
}

func TestOkWriteResult(t *testing.T) {
	out := bytes.NewBuffer(nil)
	err := writeResult(out, 5)
	require.Equal(t, nil, err)
	require.Equal(t, "5", out.String())
}

func TestOk(t *testing.T) {
	num, err := calc("2+2*2/2\r")
	require.Equal(t, nil, err)
	require.Equal(t, 4.0, num)
}

func TestOk2(t *testing.T) {
	num, err := calc("2*(1+3)\r")
	require.Equal(t, nil, err)
	require.Equal(t, 8.0, num)
}

func TestOk3(t *testing.T) {
	num, err := calc("100/(100*(100+(100/100)-1)/100)\r")
	require.Equal(t, nil, err)
	require.Equal(t, 1.0, num)
}

func TestOk4(t *testing.T) {
	num, err := calc("-2+2\r")
	require.Equal(t, nil, err)
	require.Equal(t, 0.0, num)
}

func TestOk5(t *testing.T) {
	num, err := calc("22+3*4+15/(1+3)*2+0.5\r")
	require.Equal(t, nil, err)
	require.Equal(t, 36.375, num)
}

func TestOk6(t *testing.T) {
	num, err := calc("2+2-1\r")
	require.Equal(t, nil, err)
	require.Equal(t, 3.0, num)
}

func TestFail1(t *testing.T) {
	_, err := calc("2+*2\r")
	require.EqualError(t, err, malformedExpressionErr)
}

func TestFail2(t *testing.T) {
	_, err := calc("2+2+\r")
	require.EqualError(t, err, malformedExpressionErr)
}

func TestFail3(t *testing.T) {
	_, err := calc("2^2%2\r")
	require.EqualError(t, err, "invalid character ^")
}

func TestFail4(t *testing.T) {
	_, err := calc("10/(2+3))\r")
	require.EqualError(t, err, invalidBracketsErr)
}

func TestFail5(t *testing.T) {
	_, err := calc("(10/(2+3)\r")
	require.EqualError(t, err, invalidBracketsErr)
}
