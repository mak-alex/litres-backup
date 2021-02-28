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
)

const (
	// config mode status production
	ProductionMode = "production"

	// DefaultConfigFile name of config file (toml format)
	DefaultConfigFile = "config.yaml"

	// DefaultWorkdirName name of working directory
	DefaultWorkdirName = "config"

	// DefaultPidFilename is default filename of pid file
	DefaultPidFilename = "backlitr.pid"

	// DefaultLockFilename is default filename of lock file
	DefaultLockFilename = "backlitr.lock"

	// DefaultLogFileName
	DefaultLogFileName = "backlitr.log"

	// server file dir
	DefaultSystemDataDirName = "system-data"

	// frontend static file dir
	DefaultFrontendStaticDirName = "public"

	// user file upload file dir
	DefaultUserDataDirName = "user-data"

	// temp file dir
	DefaultTempDirName = "framework-temp"
)
