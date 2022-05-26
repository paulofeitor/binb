package main

import (
	"fmt"
	"os"

	"github.com/paulofeitor/binb/internal/pkg/config"
	"github.com/paulofeitor/binb/internal/pkg/database"
	"github.com/paulofeitor/binb/internal/pkg/server"
	"github.com/paulofeitor/binb/internal/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configFile string
	flags      *pflag.FlagSet
	vp         *viper.Viper
	cmd        = &cobra.Command{
		Use:   "binb",
		Short: "Cober and Viper together at last",
		Long:  `Demonstrate how to get cobra flags to bind to viper properly`,
		Run: func(_ *cobra.Command, _ []string) {
			run()
		},
	}
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	vp = viper.New()
	var err error
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "Path to config")
	cmd.PersistentFlags().String("log-type", "JSON", "Log format")
	cmd.PersistentFlags().String("log-level", "DEBUG", "Log level")
	cmd.PersistentFlags().String("database-user", "binb", "Log level")
	cmd.PersistentFlags().String("database-pass", "password123", "Log level")
	cmd.PersistentFlags().String("database-host", "localhost:3306", "Log level")
	cmd.PersistentFlags().String("database-name", "binb", "Log level")

	flags = cmd.PersistentFlags()

	err = vp.BindPFlag("log.type", flags.Lookup("log-type"))
	exitOnErr(err)
	err = vp.BindPFlag("log.level", flags.Lookup("log-level"))
	exitOnErr(err)
	err = vp.BindPFlag("database.user", flags.Lookup("database-user"))
	exitOnErr(err)
	err = vp.BindPFlag("database.pass", flags.Lookup("database-pass"))
	exitOnErr(err)
	err = vp.BindPFlag("database.host", flags.Lookup("database-host"))
	exitOnErr(err)
	err = vp.BindPFlag("database.name", flags.Lookup("database-name"))
	exitOnErr(err)
}

func initConfig() {
	if configFile == "" {
		return
	}

	vp.SetConfigFile(configFile)
	if err := vp.ReadInConfig(); err != nil {
		fmt.Println("error reading config:", err)
		os.Exit(1)
	}
}

func getConfig() (config.Configuration, error) {
	var c config.Configuration
	err := vp.Unmarshal(&c)
	return c, err
}

func run() {
	conf, err := getConfig()
	exitOnErr(err)

	db, err := database.New(conf)
	exitOnErr(err)

	w, err := worker.New(conf, db)
	exitOnErr(err)

	err = w.Start()
	exitOnErr(err)

	s := server.New(conf, db)
	s.Start()
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
