package cmd

import (
	"github.com/mak-alex/litres-backup/pkg/conf"
	"github.com/mak-alex/litres-backup/pkg/consts"
	"github.com/mak-alex/litres-backup/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Initial config generation",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("path")

		if configPath == "" {
			configPath = filepath.Join(tools.DefaultConfigPath(), "config", consts.DefaultConfigFile)
		}

		if err := viper.Unmarshal(&conf.GlobalConfig); err != nil {
			log.Printf("Marshalling config to global struct variable: %+v\n", zap.Error(err))
		}

		if err := conf.SaveConfig(configPath); err != nil {
			log.Printf("Saving config failed: %+v\n", zap.Error(err))
		}
		log.Printf("config file is saved success and path is %s", configPath)
	},
}

func init() {
	viper.SetEnvPrefix("litres-backup")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	configCmd.Flags().BoolVarP(&conf.GlobalConfig.Mode, "mode", "m", true, "production mode, short output logs")
	_ = viper.BindPFlag("user", configCmd.Flags().Lookup("user"))

	configCmd.Flags().StringVarP(&conf.GlobalConfig.Login, "user", "u", "", "username")
	_ = viper.BindPFlag("user", configCmd.Flags().Lookup("user"))

	configCmd.Flags().StringVarP(&conf.GlobalConfig.Password, "password", "p", "", "password")
	_ = viper.BindPFlag("password", configCmd.Flags().Lookup("password"))

	configCmd.Flags().StringVarP(&conf.GlobalConfig.Library, "library", "l", filepath.Join(os.TempDir(), "litres"), "The directory where the books will be saved")
	_ = viper.BindPFlag("library", configCmd.Flags().Lookup("library"))

	// Log
	configCmd.Flags().StringVar(&conf.GlobalConfig.Log.FileLocation, "logFile", filepath.Join(tools.DefaultConfigPath(), "logs", consts.DefaultLogFileName), "log file name")
	_ = viper.BindPFlag("Log.FileName", configCmd.Flags().Lookup("logFile"))

	configCmd.Flags().IntVar(&conf.GlobalConfig.Log.MaxSize, "logFileMaxSize", 1024, "max log file size")
	_ = viper.BindPFlag("Log.MaxSize", configCmd.Flags().Lookup("logFileMaxSize"))

	configCmd.Flags().IntVar(&conf.GlobalConfig.Log.MaxBackups, "logFileBackups", 3, "number of log backup")
	_ = viper.BindPFlag("Log.MaxBackups", configCmd.Flags().Lookup("logFileBackups"))

	configCmd.Flags().IntVar(&conf.GlobalConfig.Log.MaxAge, "logFileAge", 7, "log file save max days")
	_ = viper.BindPFlag("Log.MaxAge", configCmd.Flags().Lookup("logFileAge"))

	configCmd.Flags().BoolVar(&conf.GlobalConfig.Log.Compress, "logFileCompress", true, "compress log file")
	_ = viper.BindPFlag("Log.Compress", configCmd.Flags().Lookup("logFileCompress"))
}
