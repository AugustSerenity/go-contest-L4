package service

import (
	"bufio"
	"io"
)

func SplitInput(r io.Reader) [][]string {
	scanner := bufio.NewScanner(r)
	chunk := make([]string, 0, 1000)
	var chunks [][]string

	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) == 1000 {
			chunks = append(chunks, chunk)
			chunk = make([]string, 0, 1000)
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}
