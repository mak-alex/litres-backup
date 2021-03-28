package cmd

import (
	"github.com/mak-alex/litres-backup/pkg/conf"
	"os"

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

		if l.Available4Download {
			var search string
			if conf.FilterBook.BookTitle != "" {
				search = conf.FilterBook.BookTitle
			} else if conf.FilterBook.BookID != "" {
				search = conf.FilterBook.BookID
			}
			l.Print(l.GetBooks(nil, &search))
			os.Exit(0)
		}
		_, _ = l.DownloadBooks(nil, &conf.FilterBook.BookTitle, &conf.FilterBook.BookID)
	},
}

func init() {
	bookCmd.Flags().StringVarP(&conf.FilterBook.BookTitle, "title", "t", "", "Search book by title, ex: 'Девушка, которая играла с огнем'")
	bookCmd.PersistentFlags().StringVarP(&conf.FilterBook.BookID, "id", "i", "", "Download or print book by № from available books for download")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.Progress, "progress", "b", false, "Show progress bar")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.ShowDescription, "description", "s", false, "Show description by book id or name")
	rootCmd.PersistentFlags().StringVarP(&conf.FilterBook.Format, "format", "f", "fb2.zip", "Downloading format. 'list' for available")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.NormalizedName, "normalized_name", "n", false, "Normalize book name")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.ShowAvailable4Download, "available", "a", false, "Display a list of available books for download")
}
