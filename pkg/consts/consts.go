package consts

const (
	//Debug has verbose message
	Debug = "debug"

	//Info is default log level
	Info = "info"

	//Warn is for logging messages about possible issues
	Warn = "warn"

	//Error is for logging errors
	Error = "error"

	//Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	Fatal = "fatal"

	// DefaultConfigFile name of config file (toml format)
	DefaultConfigFile = "config.yaml"

	// DefaultLogFileName
	DefaultLogFileName = "litres-backup.log"

	TB = 1000000000000
	GB = 1000000000
	MB = 1000000
	KB = 1000

	// urls
	BaseUrl      = "http://robot.litres.ru/"
	AuthorizeUrl = BaseUrl + "pages/catalit_authorise/"
	//GenresUrl       = baseUrl + "pages/catalit_genres/"
	//AuthorsUrl      = baseUrl + "pages/catalit_persons/"
	CatalogUrl = BaseUrl + "pages/catalit_browser/"
	//TrialsUrl       = baseUrl + "static/trials/"
	//PurchaseUrl     = baseUrl + "pages/purchase_book/"
	DownloadBookUrl = BaseUrl + "pages/catalit_download_book/"
)
