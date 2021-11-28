package main

import (
	"testing"
)

func TestBase64EncodeNoOver(t *testing.T) {
	var input string = "ABCDEF"
	var answer string = "QUJDREVG"
	result := base64encode(input)
	if result != answer {
		t.Fail()
	}
}

func TestBase64EncodeOver1Char(t *testing.T) {
	var input string = "ABCDEFG"
	var answer string = "QUJDREVGRw=="
	result := base64encode(input)
	if result != answer {
		t.Fail()
	}
}

func TestBase64EncodeOver2Char(t *testing.T) {
	var input string = "ABCDEFGH"
	var answer string = "QUJDREVGR0g="
	result := base64encode(input)
	if result != answer {
		t.Fail()
	}
}

func TestBase64EncodeReturnCode(t *testing.T) {
	var input string = "ABCDEFGH\r\naaa"
	var answer string = "QUJDREVGR0gNCmFhYQ=="
	result := base64encode(input)
	if result != answer {
		t.Fail()
	}
}

func TestGetChar(t *testing.T) {
	var position uint = 3
	var result string = "D"

	if result != string(getChar(position)) {
		t.Fail()
	}
}
func TestGetCharPosition(t *testing.T) {
	charByte := []byte("D")
	var result uint = 3

	if result != getPosition(charByte[0]) {
		t.Fail()
	}
}
