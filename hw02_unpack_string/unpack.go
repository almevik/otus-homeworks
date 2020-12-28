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
	strings.Builder
}

var ErrInvalidString = errors.New("invalid string")
var ErrUnpackString = errors.New("cannot unpack string")

// Записывает символ в слайс строки.
func (sS *strStruct) write() error {
	switch {
	case sS.curType == symbolType:
		_, err := sS.WriteRune(sS.curRune)

		if err != nil {
			return err
		}

	case sS.prevType == escapeSymbolType:
		if sS.curType == escapeSymbolType || sS.curType == symbolType || sS.curType == countType {
			_, err := sS.WriteRune(sS.curRune)

			if err != nil {
				return err
			}
			sS.curType = emptyType
		}

	case sS.curType == countType:
		if sS.curRune == '0' {
			subStr := getZeroCountStr(sS.String())

			sS.Reset()
			_, err := sS.WriteString(subStr)

			if err != nil {
				return err
			}
		} else {
			count, _ := strconv.Atoi(string(sS.curRune))

			// count - 1 чтобы учитывать уже имеющийся символ
			_, err := sS.WriteString(strings.Repeat(string(sS.prevRune), count-1))

			if err != nil {
				return err
			}
		}
	}

	sS.prevRune = sS.curRune
	sS.prevType = sS.curType

	return nil
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

		if sS.write() != nil {
			return "", ErrUnpackString
		}
	}

	return sS.String(), nil
}

// getZeroCountStr возвращает подстроку без последнего символа.
func getZeroCountStr(strCur string) string {
	if len(strCur) == 0 {
		return ""
	}

	return strCur[:len(strCur)-1]
}
