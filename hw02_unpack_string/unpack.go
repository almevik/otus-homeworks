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

type runeIterator struct {
	curRune  rune
	prevRune rune
	curType  runeType
	prevType runeType
}

func (rI runeIterator) isErrorType(isNumber bool, isPrevNumber bool) bool {
	return (rI.prevType == emptyType && isNumber) ||
		(rI.prevRune == escapeRune && rI.curRune == lineBreakerRune) ||
		(rI.prevRune == escapeRune && rI.curRune == tabRune) ||
		(rI.curRune == spaceRune) ||
		(rI.prevType != symbolType && isPrevNumber && isNumber)
}

func (rI runeIterator) isCountType(isNumber bool) bool {
	return (rI.prevRune != escapeRune && isNumber) || (rI.prevType == symbolType && isNumber)
}

func (rI runeIterator) isEscapeType() bool {
	return rI.prevType != escapeType && rI.curRune == escapeRune
}

// defineRuneType определяет тип руны.
func (rI *runeIterator) defineRuneType() {
	isNumber := unicode.IsDigit(rI.curRune)
	isPrevNumber := unicode.IsDigit(rI.prevRune)

	if rI.isErrorType(isNumber, isPrevNumber) {
		rI.curType = errorType
		return
	}

	if rI.isCountType(isNumber) {
		rI.curType = countType
		return
	}

	if rI.isEscapeType() {
		rI.curType = escapeType
		return
	}

	rI.curType = symbolType
}

type strStruct struct {
	runeIterator
	strBuilder strings.Builder
}

var ErrInvalidString = errors.New("invalid string")

// Записывает символ в слайс строки.
func (sS *strStruct) write() {
	switch {
	case sS.curType == symbolType:
		sS.strBuilder.WriteRune(sS.curRune)

	case sS.prevType == escapeSymbolType:
		if sS.curType == escapeSymbolType || sS.curType == symbolType || sS.curType == countType {
			sS.strBuilder.WriteRune(sS.curRune)
			sS.curType = emptyType
		}

	case sS.curType == countType:
		if sS.curRune == '0' {
			subStr := getZeroCountStr(sS.strBuilder.String())

			sS.strBuilder.Reset()
			sS.strBuilder.WriteString(subStr)
		} else {
			count, _ := strconv.Atoi(string(sS.curRune))

			// count - 1 чтобы учитывать уже имеющийся символ
			sS.strBuilder.WriteString(strings.Repeat(string(sS.prevRune), count-1))
		}
	}

	sS.prevRune = sS.curRune
	sS.prevType = sS.curType
}

func Unpack(str string) (string, error) {
	var sS strStruct

	runes := []rune(str)

	// Если пустая строка
	if len(runes) == 0 {
		return "", nil
	}

	sS.prevType = emptyType

	for _, curRune := range runes {
		sS.curRune = curRune
		sS.defineRuneType()

		if sS.curType == errorType {
			return "", ErrInvalidString
		}

		sS.write()
	}

	return sS.strBuilder.String(), nil
}

// getZeroCountStr возвращает подстроку без последнего символа.
func getZeroCountStr(strCur string) string {
	if len(strCur) == 0 {
		return ""
	}

	return strCur[:len(strCur)-1]
}
