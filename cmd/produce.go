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
	"bufio"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

// produceCmd represents the produce command
var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "produce message",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("topic is required")
			os.Exit(-1)
		}

		topic := args[0]

		exitOnError := func(err error) {
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}
		producer, err := sarama.NewAsyncProducerFromClient(c)
		exitOnError(err)

		ch := make(chan []byte)
		defer producer.Close()
		go func() {
			if err := readMessages(os.Stdin, ch); err != nil {
				exitOnError(err)
			}
		}()
		for msg := range ch {
			select {
			case producer.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(msg)}:
			case ack := <-producer.Successes():
				fmt.Println(ack)
			case err := <-producer.Errors():
				exitOnError(err)
			}
		}

	},
}

func readMessages(f *os.File, ch chan []byte) error {
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		ch <- line[0 : len(line)-1]
	}
	return nil
}

func init() {
	RootCmd.AddCommand(produceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// produceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// produceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
