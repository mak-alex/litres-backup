package main

import (
	"flag"

	"github.com/mak-alex/backlitr/litres"
)

func main() {
	login := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	bookPath := flag.String("l", "/tmp", "The directory where the books will be saved")
	format := flag.String("f", "list", "Downloading format. 'list' for available")
	search := flag.String("s", "", "Search book by title, ex: 'Девушка, которая играла с огнем'")
	progress := flag.Bool("b", false, "Show progress bar")
	normalizedName := flag.Bool("n", false, "Normalize book name")
	verbose := flag.Bool("v", false, "be verbose (this is the default)")
	debug := flag.Bool("d", false, "print lots of debugging information")
	flag.Parse()

	params := &litres.Litres{
		Login:          *login,
		Password:       *password,
		BookPath:       *bookPath,
		Format:         *format,
		Progress:       *progress,
		Verbose:        *verbose,
		Debug:          *debug,
		NormalizedName: *normalizedName,
	}

	l := litres.New(params)
	l.DownloadBooks("", *search)
}
