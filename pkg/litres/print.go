package litres

import (
	"fmt"
	"github.com/mak-alex/litres-backup/pkg/conf"
	"github.com/mak-alex/litres-backup/pkg/model"
	"github.com/mak-alex/litres-backup/pkg/table"
	"os"
	"text/tabwriter"
)

func (l *Litres) Print(books *model.CatalitFb2Books) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	header := []string{
		fmt.Sprintf("%10s", "ID"),
		"Author",
		"Title",
	}
	colors := []table.Color{
		table.Bright,
		table.Bright,
		table.Bright,
	}
	table.PrintRow(writer, table.PaintRow(colors, header))

	pPrint := func(b model.Fb2Book) {
		table.PrintRow(writer, table.PaintRow(colors, []string{
			fmt.Sprintf("%10s", b.GetID()), fmt.Sprintf("%30s", b.GetAuthor()), b.GetTitle(),
		}))

		_ = writer.Flush()
	}

	if conf.FilterBook.BookTitle != "" || conf.FilterBook.BookID != "" {
		b := books.FindBook(&conf.FilterBook.BookID, &conf.FilterBook.BookTitle)
		if b != nil {
			pPrint(*b)
			if conf.FilterBook.ShowDescription {
				fmt.Println(fmt.Sprintf("\t\t%s", b.GetDescription()))
			}
		}
	} else {
		for i, b := range books.Fb2Book {
			if i > 0 && conf.FilterBook.ShowDescription {
				fmt.Printf("\n\n\n")
			}
			pPrint(b)
			if conf.FilterBook.ShowDescription {
				fmt.Println(fmt.Sprintf("\t\t%s", b.GetDescription()))
			}
		}
	}

}
