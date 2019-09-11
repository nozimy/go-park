package main

import (
	"bytes"
	"testing"
)

func TestOk(t *testing.T) {
	in := bytes.NewBufferString("2+2*2/2")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != "4" {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestOk2(t *testing.T) {
	in := bytes.NewBufferString("2*(1+3)")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != "8" {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestOk3(t *testing.T) {
	in := bytes.NewBufferString("100/(100*(100+(100/100)-1)/100)")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != "1" {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestOk4(t *testing.T) {
	in := bytes.NewBufferString("-2+2")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != "0" {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestOk5(t *testing.T) {
	in := bytes.NewBufferString("22+3*4+15/(1+3)*2+0.5")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != "36.375" {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestOk6(t *testing.T) {
	in := bytes.NewBufferString("2+2-1")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != "3" {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestFail1(t *testing.T) {
	in := bytes.NewBufferString("2+*2")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}

func TestFail2(t *testing.T) {
	in := bytes.NewBufferString("2+2+")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}

func TestFail3(t *testing.T) {
	in := bytes.NewBufferString("2^2%2")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}

func TestFail4(t *testing.T) {
	in := bytes.NewBufferString("10/(2+3))")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}

func TestFail5(t *testing.T) {
	in := bytes.NewBufferString("(10/(2+3)")
	out := bytes.NewBuffer(nil)
	err := calc(in, out)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}
