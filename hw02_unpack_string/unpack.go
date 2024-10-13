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
	escaped := false
	allowSecondDigit := false

	for i, r := range input {
		if escaped {
			// Если предыдущий символ был '\', то просто добавляем текущий символ
			result.WriteRune(r)
			prevRune = r
			prevRuneSize = len(string(r))
			escaped = false
			allowSecondDigit = true
			continue
		}

		if r == '\\' {
			// Если текущий символ '\', то устанавливаем флаг экранирования
			escaped = true
			continue
		}

		if unicode.IsDigit(r) {
			if i == 0 || (unicode.IsDigit(prevRune) && !allowSecondDigit) {
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
		allowSecondDigit = false
	}

	if escaped {
		// Если строка заканчивается на '\', это ошибка
		return "", ErrInvalidString
	}

	return result.String(), nil
}
