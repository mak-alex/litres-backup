package conf

import (
	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/tools"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type BookFilter struct {
	BookID                 int
	SearchByTitle          string
	Format                 string
	NormalizedName         bool
	ShowAvailable4Download bool
	Progress               bool
}

type Conf struct {
	Login    string `json:"login" yaml:"login" toml:"login,omitempty"`
	Password string `json:"password" yaml:"password" toml:"password,omitempty"`
	Library  string `json:"library" yaml:"library" toml:"library,omitempty"`

	Verbose bool    `json:"verbose" yaml:"verbose" toml:"verbose,omitempty"`
	Debug   bool    `json:"debug" yaml:"debug" toml:"debug,omitempty"`
	Log     LogConf `json:"log" yaml:"log" toml:"log,omitempty"` // Logs

	Mode     bool   `json:"mode,omitempty" yaml:"mode,omitempty" toml:"mode,omitempty"`
	ConfPath string `json:"config_path" yaml:"config_path" toml:"config_path,omitempty"`
}

// LogConf represents parameters of log
type LogConf struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool

	EnableConsole     bool
	EnableFile        bool
	ConsoleJSONFormat bool
	ConsoleLevel      string

	FileJSONFormat bool
	FileLevel      string
	FileLocation   string
}

// SaveConf save global parameters to configFile
func SaveConfig(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0775)
		if err != nil {
			return errors.Wrapf(err, "creating dir %s", dir)
		}
	}

	cf, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = cf.Close()
	}()

	consoleLogLevel := consts.Warn
	if !GlobalConfig.Mode {
		consoleLogLevel = consts.Debug
	}
	cwd := tools.DefaultConfigPath()
	GlobalConfig.Log = LogConf{
		MaxSize:           2,
		MaxAge:            2,
		MaxBackups:        2,
		Compress:          true,
		EnableConsole:     true,
		ConsoleLevel:      consoleLogLevel,
		ConsoleJSONFormat: false,
		EnableFile:        true,
		FileLevel:         consoleLogLevel,
		FileJSONFormat:    false,
		FileLocation:      filepath.Join(cwd, "logs", "litres-backup.log"),
	}

	err = toml.NewEncoder(cf).Encode(GlobalConfig)
	if err != nil {
		return err
	}
	return nil
}

// LoadConf from configFile
// the function has side effect updating global var Conf
func LoadConf(path string) error {
	t := LoadConfToVar(path, &GlobalConfig)
	return t
}

// LoadConfToVar - ..
func LoadConfToVar(path string, v *Conf) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.Errorf("Unable to load config file %s", path)
	}

	viper.SetConfigFile(path)
	if err = viper.ReadInConfig(); err != nil {
		return errors.Wrapf(err, "reading config")
	}

	if err := viper.Unmarshal(&v); err != nil {
		return errors.Wrapf(err, "marshalling config to global struct variable")
	}

	if !strings.EqualFold(v.ConfPath, path) {
		v.ConfPath = path
	}

	return nil
}
