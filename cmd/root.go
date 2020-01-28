/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wpe_merge/wpe_merge/account"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile               string
	infile, outfile       *os.File
	url                   string
	maxConcurrentRequests int64

	streamer *account.WPStreamer
	ctx      context.Context
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wpe_merge <input_file> <output_file>",
	Short: "Merges input account info with data on the server",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 2 {
			return errors.New("required input_file and output_file")
		}

		infile, err = os.Open(args[0])
		if err != nil {
			return errors.Wrapf(err, "could not open file `%s`", args[0])
		}

		outfile, err = os.Create(args[1])
		if err != nil {
			infile.Close()
			return errors.Wrapf(err, "could not create file `%s`", args[1])
		}

		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		client := account.NewWPClient(url)
		streamer = account.NewWPStreamer(client, account.WithMaxConcurrentRequests(maxConcurrentRequests))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		go func() {
			<-sigCh
			cancel()
		}()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := logrus.WithContext(ctx)

		if err := streamer.Stream(ctx, infile, outfile); err != nil {
			infile.Close()
			outfile.Close()
			os.Remove(outfile.Name())
			log.WithError(err).Error("Could not stream data")
			return
		}

		infile.Close()
		outfile.Close()
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wpe_merge.yaml)")
	rootCmd.PersistentFlags().StringVar(&url, "url", "http://interview.wpengine.io/", "URL to connect to the WPE server")
	rootCmd.PersistentFlags().Int64Var(&maxConcurrentRequests, "max-concurrent-requests", 10, "max concurrent requests to make to the WPE server")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("max-concurrent-requests", rootCmd.PersistentFlags().Lookup("max-concurrent-requests"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".wpe_merge" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".wpe_merge")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
