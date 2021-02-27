## BackLitr - Assistant for backing up your books from a book resource litres.ru

More information on the [wiki](https://github.com/mak-alex/backlitr/wiki)

### BackLitr --help
```
BackLitr - Assistant for backing up your books from a book resource litres.ru

Usage:
  BackLitr [command]

Available Commands:
  book        Back up all your books, search for necessary ones, and much more
  help        Help about any command

Flags:
      --config string     config file (default is $HOME/.cobra.yaml)
  -d, --debug             print lots of debugging information
  -h, --help              help for BackLitr
  -l, --library string    The directory where the books will be saved (default "/tmp")
  -p, --password string   password
  -u, --user string       username
  -v, --verbose           be verbose (this is the default)
      --version           version for BackLitr

Use "BackLitr [command] --help" for more information about a command.
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
  -f, --format string     Downloading format. 'list' for available (default "fb2.zip")
  -t, --search_by_title string        Search book by title, ex: 'Девушка, которая играла с огнем'
  -a, --show_available_for_download   Display a list of available books for download

Global Flags:
      --config string     config file (default is $HOME/.cobra.yaml)
  -d, --debug             print lots of debugging information
  -l, --library string    The directory where the books will be saved (default "/tmp")
  -p, --password string   password
  -u, --user string       username
  -v, --verbose           be verbose (this is the default)
```
