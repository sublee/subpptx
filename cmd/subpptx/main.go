package main

import (
	"os"

	"github.com/sublee/subpptx/pkg/subpptx"
)

func main() {
	path := os.Args[1]
	outputPath := os.Args[2]
	subpptx.InjectSubtitles(path, outputPath)
}
