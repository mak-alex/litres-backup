## litres-backup - Assistant for backing up your books from a book resource litres.ru

### Opportunities:

- View a list of available books for download
- Download the entire list of available books for download
- Search and download from the list of available books to download
- Check for the existence of the downloaded book
- Generating a configuration file (~/.litres-backup/config/config.yaml)
- Logging (~/.litres-backup/logs/litres-backup.log)

More information on the [wiki](https://github.com/mak-alex/litres-backup/wiki)

### Install
```bash
$ git clone https://github.com/mak-alex/litres-backup.git $GOPATH/src/github.com/mak-alex/litres-backup
$ cd $GOPATH/src/github.com/mak-alex/litres-backup && go install
$ $GOPATH/bin/litres-backup --help
```

## How to use
### litres-backup --help
```
litres-backup - Assistant for backing up your books from a book resource litres.ru

Usage:
  litres-backup [command]

Available Commands:
  book        Back up all your books, search for necessary ones, and much more
  config      Initial config generation
  help        Help about any command

Flags:
      --config string   filepath to config.yaml (default "~/.litres-backup/config/config.yaml")
  -d, --debug           print lots of debugging information
  -f, --format string   Downloading format. 'list' for available (default "fb2.zip")
  -h, --help            help for litres-backup
  -v, --verbose         be verbose (this is the default)
      --version         version for litres-backup

Use "litres-backup [command] --help" for more information about a command.
```

### litres-backup config --help
```
Initial config generation

Usage:
  litres-backup config [flags]

Flags:
  -h, --help                 help for config
  -l, --library string       The directory where the books will be saved (default "/tmp")
      --logFile string       log file name (default "~/.litres-backup/logs/litres-backup.log")
      --logFileAge int       log file save max days (default 7)
      --logFileBackups int   number of log backup (default 3)
      --logFileCompress      compress log file (default true)
      --logFileMaxSize int   max log file size (default 1024)
  -m, --mode                 production mode, short output logs (default true)
  -p, --password string      password
  -u, --user string          username

Global Flags:
      --config string   filepath to config.yaml (default "~/.litres-backup/config/config.yaml")
  -d, --debug           print lots of debugging information
  -f, --format string   Downloading format. 'list' for available (default "fb2.zip")
  -v, --verbose         be verbose (this is the default)

```

### litres-backup book --help
```
Back up all your books, search for necessary ones, and much more

Usage:
  litres-backup book [OPTION]... [flags]

Flags:
  -a, --available         Display a list of available books for download
  -s, --description       Show description by book id or name
  -h, --help              help for book
  -i, --id string         Download or print book by № from available books for download
  -n, --normalized_name   Normalize book name
  -b, --progress          Show progress bar
  -t, --title string      Search book by title, ex: 'Девушка, которая играла с огнем'

Global Flags:
      --config string   filepath to config.yaml (default "~/.litres-backup/config/config.yaml")
  -d, --debug           print lots of debugging information
  -f, --format string   Downloading format. 'list' for available (default "fb2.zip")
  -v, --verbose         be verbose (this is the default)

```
