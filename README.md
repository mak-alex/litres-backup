## BackLitr - Assistant for backing up your books from a book resource litres.ru

### Opportunities:

- View a list of available books for download
- Download the entire list of available books for download
- Search and download from the list of available books to download
- Check for the existence of the downloaded book
- Generating a configuration file (~/.backlitr/config/config.yaml)
- Logging (~/.backlitr/logs/backlitr.log)

More information on the [wiki](https://github.com/mak-alex/backlitr/wiki)

### BackLitr --help
```
BackLitr - Assistant for backing up your books from a book resource litres.ru

Usage:
  BackLitr [command]

Available Commands:
  book        Back up all your books, search for necessary ones, and much more
  config      Initial config generation
  help        Help about any command

Flags:
      --config string   filepath to config.yaml (default "/home/mak/.backlitr/config/config.yaml")
  -d, --debug           print lots of debugging information
  -f, --format string   Downloading format. 'list' for available (default "fb2.zip")
  -h, --help            help for BackLitr
  -v, --verbose         be verbose (this is the default)
      --version         version for BackLitr

Use "BackLitr [command] --help" for more information about a command.
```

### BackLitr config --help
```
Initial config generation

Usage:
  BackLitr config [flags]

Flags:
  -h, --help                 help for config
  -l, --library string       The directory where the books will be saved (default "/tmp")
      --logFile string       log file name (default "/home/mak/.backlitr/logs/backlitr.log")
      --logFileAge int       log file save max days (default 7)
      --logFileBackups int   number of log backup (default 3)
      --logFileCompress      compress log file (default true)
      --logFileMaxSize int   max log file size (default 1024)
      --mode string          server run mode (default "debug")
  -p, --password string      password
  -u, --user string          username

Global Flags:
      --config string   filepath to config.yaml (default "/home/mak/.backlitr/config/config.yaml")
  -d, --debug           print lots of debugging information
  -f, --format string   Downloading format. 'list' for available (default "fb2.zip")
  -v, --verbose         be verbose (this is the default)
```

### BackLitr book --help
```
Back up all your books, search for necessary ones, and much more

Usage:
  BackLitr book [OPTION]... [flags]

Flags:
  -i, --book_id int                   Download or print book by № from available books for download (default -1)
  -h, --help                          help for book
  -n, --normalized_name               Normalize book name
  -b, --progress                      Show progress bar
  -t, --search_by_title string        Search book by title, ex: 'Девушка, которая играла с огнем'
  -a, --show_available_for_download   Display a list of available books for download

Global Flags:
      --config string   filepath to config.yaml (default "/home/mak/.backlitr/config/config.yaml")
  -d, --debug           print lots of debugging information
  -f, --format string   Downloading format. 'list' for available (default "fb2.zip")
  -v, --verbose         be verbose (this is the default)
```
