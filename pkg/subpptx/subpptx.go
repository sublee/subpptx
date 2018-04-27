package subpptx

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/mholt/archiver"
)

// Finds and removes elements
func removeElements(elem *etree.Element, path string) {
	for _, e := range elem.FindElements(path) {
		e.Parent().RemoveChild(e)
	}
}

// Injects a subtitle for one slide file.
func injectSubtitle(slidePath string, done chan bool) {
	defer func() {
		done <- true
	}()

	fmt.Println(slidePath)

	relsPath := filepath.Join(filepath.Dir(slidePath), "_rels", filepath.Base(slidePath)+".rels")

	rels := etree.NewDocument()
	rels.ReadFromFile(relsPath)

	notesRel := rels.FindElement("//Relationship[@Type='http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide']")
	if notesRel == nil {
		// This slide does not have notes
		return
	}

	notesPath := notesRel.SelectAttrValue("Target", "")
	notesPath = filepath.Join(filepath.Dir(slidePath), notesPath)

	notes := etree.NewDocument()
	notes.ReadFromFile(notesPath)

	// Find the subtitle body in the notes
	e := notes.FindElement("//p:ph[@type='body']")
	e = e.Parent().Parent().Parent()
	txBody := e.FindElement("./p:txBody")

	// Eliminate harmful tags
	removeElements(txBody, ".//a:pPr")
	removeElements(txBody, ".//a:rPr")

	// Create a new footer
	tmpDoc := etree.NewDocument()
	tmpDoc.ReadFromString(`
		<p:sp>
			<p:nvSpPr>
				<p:cNvPr id="999999" name="Subtitle" />
				<p:cNvSpPr><a:spLocks noGrp="1" /></p:cNvSpPr>
				<p:nvPr><p:ph type="ftr" sz="quarter" idx="11" /></p:nvPr>
			</p:nvSpPr>
			<p:spPr />
		</p:sp>
	`)
	sp := tmpDoc.Root()
	sp.AddChild(txBody)

	// Inject the footer
	slide := etree.NewDocument()
	slide.ReadFromFile(slidePath)
	spTree := slide.FindElement("//p:spTree")
	spTree.AddChild(sp)

	slide.WriteToFile(slidePath)
}

func InjectSubtitles(pptxPath string, outputPath string) {
	// Create temp dir
	pptxBase := filepath.Base(pptxPath)
	tempDir, err := ioutil.TempDir("", pptxBase)
	if err != nil {
		panic("failed to create temp dir")
	}

	// Unzip PPTX into temp dir
	archiver.Zip.Open(pptxPath, tempDir)

	// Collect files in ppt/slides
	slidesDir := filepath.Join(tempDir, "ppt", "slides")
	slidesFiles, err := ioutil.ReadDir(slidesDir)
	if err != nil {
		panic("PPTX directory structure corruptted")
	}

	// Inject subtitles over slide*.xml
	count := 0
	done := make(chan bool)

	for _, file := range slidesFiles {
		if strings.HasPrefix(file.Name(), "slide") && strings.HasSuffix(file.Name(), ".xml") {
			slidePath := filepath.Join(slidesDir, file.Name())
			go injectSubtitle(slidePath, done)
			count += 1
		}
	}

	// Join the goroutines
	for i := 0; i < count; i++ {
		<-done
	}

	// Save as a new PPTX
	rootFiles, err := ioutil.ReadDir(tempDir)
	if err != nil {
		panic("failed to read temp dir")
	}
	filenames := make([]string, len(rootFiles))
	for i, file := range rootFiles {
		filenames[i] = filepath.Join(tempDir, file.Name())
	}

	archiver.Zip.Make(outputPath, filenames)
}