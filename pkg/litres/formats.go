package litres

import (
	"fmt"
	"os"
	"strings"
)

var (
	formats = []string{
		"fb2.zip",
		"html",
		"html.zip",
		"txt",
		"txt.zip",
		"rtf.zip",
		"a4.pdf",
		"a6.pdf",
		"mobi.prc",
		"epub",
		"ios.epub",
		"fb3",
	}
)

func (l *Litres) existFormat() bool {
	for _, format := range formats {
		if strings.Contains(l.Format, format) {
			return true
		}
	}
	return false
}

func (l *Litres) printFormats() {
	fmt.Println("Available formats:")
	for _, format := range formats {
		fmt.Println("\t -", format)
	}
	os.Exit(0)
}
