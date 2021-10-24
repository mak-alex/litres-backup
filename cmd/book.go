package cmd

import (
	"github.com/mak-alex/litres-backup/pkg/conf"
	"github.com/mak-alex/litres-backup/tools"

	"github.com/mak-alex/litres-backup/pkg/litres"
	"github.com/spf13/cobra"
)

var bookCmd = &cobra.Command{
	Use:    "book [OPTION]...",
	Short:  "Back up all your books, search for necessary ones, and much more",
	PreRun: loadConfig,
	Run: func(cmd *cobra.Command, args []string) {
		params := &litres.Litres{
			Login:              conf.GlobalConfig.Login,
			Password:           conf.GlobalConfig.Password,
			Library:            conf.GlobalConfig.Library,
			Format:             conf.FilterBook.Format,
			Progress:           conf.FilterBook.Progress,
			Verbose:            conf.GlobalConfig.Verbose,
			Debug:              conf.GlobalConfig.Debug,
			NormalizedName:     conf.FilterBook.NormalizedName,
			Available4Download: conf.FilterBook.ShowAvailable4Download,
		}
		l := litres.New(params)
		_ = tools.MakeDirectory(l.Library)

		if conf.FilterBook.BookID != 0 {
			l.Print(l.LoadBooksByBookId(conf.FilterBook.BookID, true))
			return
		}
		l.Print(l.LoadPurchasedBooks(&conf.FilterBook.BookTitle, conf.FilterBook.Offset, conf.FilterBook.MaxCount))
	},
}

func init() {
	bookCmd.Flags().StringVarP(&conf.FilterBook.BookTitle, "title", "t", "", "Search book by title, ex: 'Девушка, которая играла с огнем'")
	bookCmd.PersistentFlags().IntVarP(&conf.FilterBook.BookID, "id", "i", 0, "Download or print book by № from available books for download")
	bookCmd.PersistentFlags().IntVarP(&conf.FilterBook.Offset, "offset", "o", 0, "default offset 0")
	bookCmd.PersistentFlags().IntVarP(&conf.FilterBook.MaxCount, "max_count", "m", 100, "default max count 100")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.Progress, "progress", "b", false, "Show progress bar")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.ShowDescription, "description", "s", false, "Show description by book id or name")
	rootCmd.PersistentFlags().StringVarP(&conf.FilterBook.Format, "format", "f", "fb2.zip", "Downloading format. 'list' for available")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.NormalizedName, "normalized_name", "n", false, "Normalize book name")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.ShowAvailable4Download, "available", "a", false, "Display a list of available books for download")
}
