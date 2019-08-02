// Package cmd contains all the commands used by the people project
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "people",
	Short: "Command to manage address book",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file(default is $HOME/.peoplerc)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.SetDefault("logfile", "/tmp/people.log")
}

func initConfig() {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		viper.AddConfigPath(user.HomeDir)
		viper.SetConfigName(".people")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	file, err := os.OpenFile(viper.GetString("logfile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
