package magazine

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
)

func CreatePDF(mag *Magazine) error {
	log.Print("Generating PDF")
	files := []string{}
	sort.Sort(ByNumber(mag.Pages))
	for _, page := range mag.Pages {
		files = append(files, page.TmpFile)
	}
	// Output file
	files = append(files, mag.formatMagazineTitle())
	cmd := exec.Command("convert", files...)
	c, _ := cmd.CombinedOutput()
	fmt.Println(string(c))
	return nil
}

type ByNumber []*Page

func (a ByNumber) Len() int           { return len(a) }
func (a ByNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNumber) Less(i, j int) bool { return a[i].PageNumber < a[j].PageNumber }
