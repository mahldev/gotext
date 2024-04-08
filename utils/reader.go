package utils

import (
	"math/rand"
	"os"
	"strings"
)

func ReadWordFile() []string {
	content, err := os.ReadFile("./words.txt")
	if err != nil {
		panic(err.Error())
	}
	lines := strings.Split(string(content), "\n")

	return lines
}

func ReadWordFileN(n uint64) []string {
	content, err := os.ReadFile("./words.txt")
	if err != nil {
		panic(err.Error())
	}
	lines := strings.Split(string(content), "\n")

	if int(n) > len(lines) {
		n = uint64(len(lines))
	}

	rand.Shuffle(len(lines), func(i, j int) {
		lines[i], lines[j] = lines[j], lines[i]
	})

	result := make([]string, n)
	copy(result, lines[:n])

	return result
}
