package cmd

import (
	"fmt"
	"github.com/mak-alex/backlitr/pkg/consts"
	"github.com/mak-alex/backlitr/pkg/logger"
	"github.com/mak-alex/backlitr/tools"
	"github.com/spf13/viper"
	"os"
	"path/filepath"

	"github.com/mak-alex/backlitr/pkg/conf"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "BackLitr",
		Short:   "BackLitr - Assistant for backing up your books from a book resource litres.ru",
		Version: "0.0.1",
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Executing root command")
	}
}

func init() {
	viper.SetConfigType("toml")
	rootCmd.AddCommand(
		configCmd,
		bookCmd,
	)

	rootCmd.PersistentFlags().BoolVarP(&conf.GlobalConfig.Verbose, "verbose", "v", false, "be verbose (this is the default)")
	rootCmd.PersistentFlags().BoolVarP(&conf.GlobalConfig.Debug, "debug", "d", false, "print lots of debugging information")
	rootCmd.PersistentFlags().StringVar(&conf.GlobalConfig.ConfPath, "config", filepath.Join(tools.DefaultConfigPath(), "config", consts.DefaultConfigFile), "filepath to config.yaml")
}

// Load the configuration from file
func loadConfig(cmd *cobra.Command, args []string) {
	err := conf.LoadConf(conf.GlobalConfig.ConfPath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Println("Please configure the server first before starting it: backlitr config --help")
		os.Exit(-1)
	}
	logger.Work = logger.NewLogger(conf.GlobalConfig.Log, logger.InstanceZapLogger)
}
