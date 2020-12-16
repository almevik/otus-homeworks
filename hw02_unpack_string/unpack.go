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

const escapeRune rune = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var strBuilder strings.Builder

	runes := []rune(str)

	// Если пустая строка
	if len(runes) == 0 {
		return "", nil
	}

	var prevRune rune
	prevType := emptyType

	for _, curRune := range runes {
		currentType := defineRuneType(prevRune, curRune, prevType)
		prevType = currentType

		switch {
		case currentType == symbolType:
			strBuilder.WriteRune(curRune)

		case currentType == escapeSymbolType:

		case prevType == escapeSymbolType:
			if currentType == escapeSymbolType || currentType == symbolType || currentType == countType {
				strBuilder.WriteRune(curRune)
				prevType = emptyType
			}

		case currentType == countType:
			if curRune == '0' {
				subStr := getZeroCountStr(strBuilder.String())

				strBuilder.Reset()
				strBuilder.WriteString(subStr)

				break
			}

			count, err := strconv.Atoi(string(curRune))

			if err != nil {
				panic("Can't convert countRune to int")
			}

			// count - 1 чтобы учитывать уже имеющийся символ
			strBuilder.WriteString(strings.Repeat(string(prevRune), count-1))

		case currentType == errorType:
			return "", ErrInvalidString
		}

		prevRune = curRune
	}

	return strBuilder.String(), nil
}

// getZeroCountStr возвращает подстроку без последнего символа.
func getZeroCountStr(strCur string) string {
	var subStrBuilder strings.Builder

	runesSubStr := []rune(strCur)[:len(strCur)-1]

	for _, runeSubStr := range runesSubStr {
		subStrBuilder.WriteRune(runeSubStr)
	}

	return subStrBuilder.String()
}

// defineRuneType определяет тип руны.
func defineRuneType(prevRune rune, currentRune rune, prevType runeType) runeType {
	isNumber := unicode.IsDigit(currentRune)
	isPrevNumber := unicode.IsDigit(prevRune)

	switch {
	case prevType == emptyType && isNumber:
		return errorType
	case prevType != symbolType && isPrevNumber && isNumber:
		return errorType
	case (prevRune != escapeRune && isNumber) || (prevType == symbolType && isNumber):
		return countType
	case prevRune == escapeRune && currentRune == escapeRune && prevType == escapeType:
		return symbolType
	case prevType != escapeType && currentRune == escapeRune:
		return escapeType
	default:
		return symbolType
	}
}
