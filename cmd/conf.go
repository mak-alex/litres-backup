package cmd

import (
	"github.com/mak-alex/backlitr/pkg/conf"
	"github.com/mak-alex/backlitr/pkg/consts"
	"github.com/mak-alex/backlitr/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"path/filepath"
	"strings"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Initial config generation",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("path")
		if err := conf.FillRuntimePaths(); err != nil {
			log.Fatal("Filling config", zap.Error(err))
		}

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
	viper.SetEnvPrefix("backlitr")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// run mode
	configCmd.Flags().String("mode", "debug", "server run mode")
	_ = viper.BindPFlag("Mode", configCmd.Flags().Lookup("mode"))

	configCmd.Flags().StringVarP(&conf.GlobalConfig.Login, "user", "u", "", "username")
	_ = viper.BindPFlag("user", configCmd.Flags().Lookup("user"))

	configCmd.Flags().StringVarP(&conf.GlobalConfig.Password, "password", "p", "", "password")
	_ = viper.BindPFlag("password", configCmd.Flags().Lookup("password"))

	configCmd.Flags().StringVarP(&conf.GlobalConfig.Library, "library", "l", "/tmp", "The directory where the books will be saved")
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
