package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

type wordStruct struct {
	word  string
	count int
}

func Top10(str string) []string {
	if len(str) == 0 {
		return nil
	}

	words := make(map[string]int)
	var sb strings.Builder

	for _, r := range []rune(str) {
		if r == ' ' || r == '\n' || r == '\t' {
			if sb.Len() > 0 {
				words[sb.String()] = words[sb.String()] + 1
				sb.Reset()
			}
			continue
		}
		sb.WriteRune(r)
	}

	var wordsSlice []wordStruct

	for word, count := range words {
		wordsSlice = append(wordsSlice, wordStruct{word, count})
	}

	sort.Slice(wordsSlice, func(i, j int) bool {
		return wordsSlice[i].count > wordsSlice[j].count
	})

	var top10 []string

	for _, value := range wordsSlice[:10] {
		top10 = append(top10, value.word)
	}

	return top10
}
