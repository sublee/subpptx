package main

import (
	"github.com/spf13/cobra"

	"github.com/sublee/subpptx/pkg/subpptx"
)

var cmd = &cobra.Command{
	Use:   "subpptx [pptx] [output]",
	Short: "Inject subtitles in MS PowerPoint files",

	Args: cobra.ExactArgs(2),
	DisableFlagsInUseLine: true,

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		outputPath := args[1]
		subpptx.InjectSubtitles(path, outputPath)
	},
}

func main() {
	cmd.Execute()
}
