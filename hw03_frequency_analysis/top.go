package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	wordCount := make(map[string]int)

	// Разделяем текст на слова
	words := strings.Fields(text)

	// Подсчитываем частоту каждого слова
	for _, word := range words {
		if len(word) > 0 && word != "-" {
			word = cleanAndLowercaseWord(word)
			wordCount[word]++
		}
	}

	// Создаем слайс для сортировки
	type wordFrequency struct {
		word  string
		count int
	}

	frequencies := make([]wordFrequency, 0, len(wordCount))
	for word, count := range wordCount {
		frequencies = append(frequencies, wordFrequency{word, count})
	}

	// Сортируем по частоте и лексикографически
	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}
		return frequencies[i].count > frequencies[j].count
	})

	// Формируем результат
	result := make([]string, 0, 10)
	for i := 0; i < len(frequencies) && i < 10; i++ {
		result = append(result, frequencies[i].word)
	}

	return result
}

// Дополнительная функция для очистки слов от знаков препинания и приведения к нижнему регистру.
func cleanAndLowercaseWord(input string) string {
	// Регулярное выражение для удаления знаков препинания по краям слова
	re := regexp.MustCompile(`^[\p{P}]+|[\p{P}]+$`)

	// Приводим строку к нижнему регистру и удаляем знаки препинания по краям
	cleanedWord := re.ReplaceAllString(input, "")
	return strings.ToLower(cleanedWord)
}
