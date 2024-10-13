package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	var prevRune rune
	var prevRuneSize int

	for i, r := range input {
		if unicode.IsDigit(r) {
			if i == 0 || unicode.IsDigit(prevRune) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(r))
			if count == 0 {
				// Удаляем последний добавленный символ
				resultStr := result.String()
				result.Reset()
				result.WriteString(resultStr[:len(resultStr)-prevRuneSize])
			} else {
				result.WriteString(strings.Repeat(string(prevRune), count-1))
			}
		} else {
			result.WriteRune(r)
			prevRuneSize = len(string(r))
		}
		prevRune = r
	}

	return result.String(), nil
}
