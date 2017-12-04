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
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wvanbergen/kazoo-go"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete topics",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("topic name is require")
			os.Exit(-1)
		}

		cfg := kazoo.NewConfig()
		cfg.Logger = log.New(ioutil.Discard, "[kafka-cli] ", log.LstdFlags)
		if verbose {
			cfg.Logger = log.New(os.Stderr, "[kafka-cli] ", log.LstdFlags)
		}

		servers := strings.Split(zookeepers, ",")
		kz, err := kazoo.NewKazoo(servers, cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		defer kz.Close()

		for _, topic := range args {
			err = kz.DeleteTopic(topic)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}
	},
}

func init() {
	topicsCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
