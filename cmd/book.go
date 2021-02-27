package cmd

import (
	"os"

	"github.com/mak-alex/backlitr/litres"
	"github.com/spf13/cobra"
)

var bookCmd = &cobra.Command{
	Use:   "book [OPTION]...",
	Short: "Back up all your books, search for necessary ones, and much more",
	Run: func(cmd *cobra.Command, args []string) {
		params := &litres.Litres{
			Login:              config.Login,
			Password:           config.Password,
			BookPath:           config.BookPath,
			Format:             config.Format,
			Progress:           config.Progress,
			Verbose:            config.Verbose,
			Debug:              config.Debug,
			NormalizedName:     config.NormalizedName,
			Available4Download: config.ShowAvailable4Download,
		}
		l := litres.New(params)

		if l.Available4Download {
			list := l.GetBooks(nil, &config.SearchByTitle)
			l.ShowAvailable4Download(list.Fb2Book)
			os.Exit(0)
		} else {
			_, _ = l.DownloadBooks(nil, &config.SearchByTitle, &config.BookID)
		}

	},
}

func init() {
	bookCmd.Flags().StringVarP(&config.SearchByTitle, "search_by_title", "t", "", "Search book by title, ex: 'Девушка, которая играла с огнем'")
	bookCmd.PersistentFlags().IntVarP(&config.BookID, "book_id", "i", -1, "Download or print book by № from available books for download")
	bookCmd.PersistentFlags().BoolVarP(&config.Progress, "progress", "b", false, "Show progress bar")
	bookCmd.PersistentFlags().BoolVarP(&config.NormalizedName, "normalized_name", "n", false, "Normalize book name")
	bookCmd.PersistentFlags().BoolVarP(&config.ShowAvailable4Download, "show_available_for_download", "a", false, "Display a list of available books for download")
}
