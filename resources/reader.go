package resources

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func ReadWordFileN(n *uint64, lang string) []string {
	fileToRead := fmt.Sprintf("./resources/langs/%s.txt", lang)
	content, err := os.ReadFile(fileToRead)
	if err != nil {
		panic(err.Error())
	}

	lines := strings.Split(string(content), "\n")

	if n == nil || int(*n) > len(lines) {
		if n == nil {
			return lines
		}
		*n = uint64(len(lines))
	}

	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})

	result := make([]string, *n)
	copy(result, lines[:*n])

	return result
}
