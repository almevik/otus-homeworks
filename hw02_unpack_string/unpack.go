package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type runeType int

const (
	symbolType runeType = iota
	countType
	escapeSymbolType
	escapeType
	errorType
	emptyType
)

const (
	escapeRune      rune = '\\'
	lineBreakerRune rune = 'n'
	tabRune         rune = 't'
	spaceRune       rune = ' '
)

type strStruct struct {
	curRune    rune
	prevRune   rune
	curType    runeType
	prevType   runeType
	strBuilder strings.Builder
}

var ErrInvalidString = errors.New("invalid string")

// Записывает символ в слайс строки.
func (strStruct *strStruct) write() {
	switch {
	case strStruct.curType == symbolType:
		strStruct.strBuilder.WriteRune(strStruct.curRune)

	case strStruct.prevType == escapeSymbolType:
		if strStruct.curType == escapeSymbolType || strStruct.curType == symbolType || strStruct.curType == countType {
			strStruct.strBuilder.WriteRune(strStruct.curRune)
			strStruct.curType = emptyType
		}

	case strStruct.curType == countType:
		if strStruct.curRune == '0' {
			subStr := getZeroCountStr(strStruct.strBuilder.String())

			strStruct.strBuilder.Reset()
			strStruct.strBuilder.WriteString(subStr)
		} else {
			count, _ := strconv.Atoi(string(strStruct.curRune))

			// count - 1 чтобы учитывать уже имеющийся символ
			strStruct.strBuilder.WriteString(strings.Repeat(string(strStruct.prevRune), count-1))
		}
	}

	strStruct.prevRune = strStruct.curRune
	strStruct.prevType = strStruct.curType
}

func Unpack(str string) (string, error) {
	var strStruct strStruct

	runes := []rune(str)

	// Если пустая строка
	if len(runes) == 0 {
		return "", nil
	}

	strStruct.prevType = emptyType

	for _, curRune := range runes {
		strStruct.curRune = curRune
		strStruct.curType = defineRuneType(strStruct)

		if strStruct.curType == errorType {
			return "", ErrInvalidString
		}

		strStruct.write()
	}

	return strStruct.strBuilder.String(), nil
}

// getZeroCountStr возвращает подстроку без последнего символа.
func getZeroCountStr(strCur string) string {
	if len(strCur) == 0 {
		return ""
	}

	return strCur[:len(strCur)-1]
}

// defineRuneType определяет тип руны.
func defineRuneType(strStruct strStruct) runeType {
	isNumber := unicode.IsDigit(strStruct.curRune)
	isPrevNumber := unicode.IsDigit(strStruct.prevRune)

	if isErrorType(strStruct, isNumber, isPrevNumber) {
		return errorType
	}

	if isCountType(strStruct, isNumber) {
		return countType
	}

	if isEscapeType(strStruct) {
		return escapeType
	}

	return symbolType
}

func isErrorType(strStruct strStruct, isNumber bool, isPrevNumber bool) bool {
	return (strStruct.prevType == emptyType && isNumber) ||
		(strStruct.prevRune == escapeRune && strStruct.curRune == lineBreakerRune) ||
		(strStruct.prevRune == escapeRune && strStruct.curRune == tabRune) ||
		(strStruct.curRune == spaceRune) ||
		(strStruct.prevType != symbolType && isPrevNumber && isNumber)
}

func isCountType(strStruct strStruct, isNumber bool) bool {
	return (strStruct.prevRune != escapeRune && isNumber) || (strStruct.prevType == symbolType && isNumber)
}

func isEscapeType(strStruct strStruct) bool {
	return strStruct.prevType != escapeType && strStruct.curRune == escapeRune
}
