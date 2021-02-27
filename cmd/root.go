package cmd

import (
	"fmt"

	"github.com/mak-alex/backlitr/conf"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config = &conf.Conf{}
var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:     "BackLitr",
		Short:   "BackLitr - Assistant for backing up your books from a book resource litres.ru",
		Version: "0.0.1",
	}
)

// Execute executes the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Executing root command")
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.Login, "user", "u", "", "username")
	rootCmd.PersistentFlags().StringVarP(&config.Password, "password", "p", "", "password")
	rootCmd.PersistentFlags().StringVarP(&config.BookPath, "library", "l", "/tmp", "The directory where the books will be saved")
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "be verbose (this is the default)")
	rootCmd.PersistentFlags().BoolVarP(&config.Debug, "debug", "d", false, "print lots of debugging information")

	_ = viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	_ = viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("password"))
	_ = viper.BindPFlag("library", rootCmd.PersistentFlags().Lookup("library"))

	rootCmd.AddCommand(bookCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".backlitr")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
