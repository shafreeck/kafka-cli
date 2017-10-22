// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg *sarama.Config = sarama.NewConfig()
var brokers string
var c sarama.Client

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kafka-cli",
	Short: "kafka-cli utility",
	Long: `kafka-cli is a console util tool to access kafka cluster
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		c, err = sarama.NewClient(strings.Split(brokers, ","), cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if c != nil {
		c.Close()
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	cfg = sarama.NewConfig()

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kafka-cli.yaml)")
	RootCmd.PersistentFlags().DurationVar(&cfg.Net.DialTimeout, "net.dialtimeout", 30*time.Second, "timeout of dialing to brokers")
	RootCmd.PersistentFlags().DurationVar(&cfg.Net.ReadTimeout, "net.readtimeout", 30*time.Second, "timeout of reading messages")
	RootCmd.PersistentFlags().DurationVar(&cfg.Net.WriteTimeout, "net.writetimeout", 30*time.Second, "timeout of writing messages")
	RootCmd.PersistentFlags().IntVar(&cfg.Net.MaxOpenRequests, "net.maxopenrequests", 5, "how many outstanding requests a connection is allowed to have before sending on it blocks")

	RootCmd.PersistentFlags().DurationVar(&cfg.Net.KeepAlive, "net.keepalive", 0, "keepalive period, 0 means disabled")
	RootCmd.PersistentFlags().IntVar(&cfg.ChannelBufferSize, "buffersize", 256, "internal channel buffer size")
	RootCmd.PersistentFlags().StringVar(&cfg.ClientID, "clientid", "kafka-cli", "a user-provided string sent with every request to the brokers for logging debugging, and auditing purposes")
	RootCmd.PersistentFlags().DurationVar(&cfg.Metadata.RefreshFrequency, "metadata.refresh", 10*time.Minute, "metadata refresh frequency")
	RootCmd.PersistentFlags().IntVar(&cfg.Metadata.Retry.Max, "metadata.retry.max", 3, "total number to request metadata when the cluster has a leader election")
	RootCmd.PersistentFlags().DurationVar(&cfg.Metadata.Retry.Backoff, "metadata.retry.backoff", 250*time.Millisecond, "backoff between retrying")

	RootCmd.PersistentFlags().StringVar(&brokers, "brokers", "127.0.0.1:9092", "broker list, delimited by comma")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".kafka-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")      // adding home directory as first search path
	viper.AutomaticEnv()              // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
