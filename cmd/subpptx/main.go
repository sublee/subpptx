package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/sublee/subpptx/pkg/subpptx"
)

func runCommand(path string, outputPath string) {
	l := log.New(os.Stderr, "", 0)
	t := time.Now()
	monitor := make(chan int)

	go subpptx.InjectSubtitles(path, outputPath, monitor)

	l.Printf(`processing "%s"...`, path)

	numSlides := <-monitor
	step := numSlides / 10

	l.Printf("%d slides found (%s elapsed)", numSlides, time.Since(t))

	for i := 0; i < numSlides; i++ {
		<-monitor

		if step == 0 {
			// Avoid zero division error
			continue
		}
		if (i+1)%step == 0 {
			percent := (float64(i+1) / float64(numSlides)) * 100
			l.Printf("%3.0f%% of slides modified (%s elapsed)", percent, time.Since(t))
		}
	}

	<-monitor
	l.Printf(`saved at "%s" (%s elapsed)`, outputPath, time.Since(t))

	<-monitor
	l.Printf("done (%s elapsed)", time.Since(t))
}

var cmd = &cobra.Command{
	Use:   "subpptx <pptx> <output>",
	Short: "Inject subtitles in MS PowerPoint files",

	Args: cobra.ExactArgs(2),
	DisableFlagsInUseLine: true,

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		outputPath := args[1]
		runCommand(path, outputPath)
	},
}

func main() {
	cmd.Execute()
}
