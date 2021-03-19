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
			list := l.GetBooks(nil, &conf.FilterBook.SearchByTitle)
			l.ShowAvailable4Download(list.Fb2Book)
			os.Exit(0)
		} else {
			_, _ = l.DownloadBooks(nil, &conf.FilterBook.SearchByTitle, &conf.FilterBook.BookID)
		}

	},
}

func init() {
	bookCmd.Flags().StringVarP(&conf.FilterBook.SearchByTitle, "search_by_title", "t", "", "Search book by title, ex: 'Девушка, которая играла с огнем'")
	bookCmd.PersistentFlags().IntVarP(&conf.FilterBook.BookID, "book_id", "i", -1, "Download or print book by № from available books for download")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.Progress, "progress", "b", false, "Show progress bar")
	rootCmd.PersistentFlags().StringVarP(&conf.FilterBook.Format, "format", "f", "fb2.zip", "Downloading format. 'list' for available")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.NormalizedName, "normalized_name", "n", false, "Normalize book name")
	bookCmd.PersistentFlags().BoolVarP(&conf.FilterBook.ShowAvailable4Download, "show_available_for_download", "a", false, "Display a list of available books for download")
}
