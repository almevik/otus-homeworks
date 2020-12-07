package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"log"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	// Place your code here
	//unicode.IsDigit
	//strings.Builder
	//strings.Repeat
	//strconv.Atoi
	var res string

	for _, char := range []rune(str) {
		res = res + string(char)

		// Если встретили цифру
		if unicode.IsDigit(char) {

		}

		if !unicode.IsDigit(char) {

		}
	}

	return res, nil
}

func itoa(i int) (s string) {
	minus := ``
	if i == 0 {
		return "0"
	}
	if i < 0 {
		i *= -1
		minus = `-`
	}

	for i != 0 {
		s = string('0'+i%10) + s
		i = i / 10
	}

	return minus + s
}

func main() {
	type pair struct {
		input    string
		expected string
	}

	test := []pair{
		{"a4bc2d5e", "aaaabccddddde"},
		{"abcd", "abcd"},
		{"3abc", ""},
		{"45", ""},
		{"aaa10b", ""},
		{"aaa0b", "aab"},
		{"", ""},
		{"d\n5abc", "d\n\n\n\n\nabc"},
	}

	for _, t := range test {
		expected, err := Unpack(t.input)

		if err != nil {
			log.Fatalf(err.Error())
		}

		if expected == t.expected {
			fmt.Printf("%s - %s\n", expected, t.expected, "OK")
		} else {
			fmt.Printf("%s - %s\n", expected, t.expected, "FAIL")
		}
	}
}
