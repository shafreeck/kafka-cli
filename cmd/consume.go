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
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

type consumerOptions struct {
	Partition int32
	Offset    int64
}

var consumerOpt consumerOptions

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "consume from kafka",
	Run: func(cmd *cobra.Command, args []string) {
		topics := args
		if len(topics) == 0 {
			fmt.Println("topic is required\n")
			displayTopics()
			os.Exit(-1)
		}

		consumer, err := sarama.NewConsumerFromClient(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		messages := make(chan *sarama.ConsumerMessage, cfg.ChannelBufferSize)
		consume := func(ctx context.Context, pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					messages <- msg
				case <-ctx.Done():
					pc.Close()
					return
				}
			}
		}
		exitOnError := func(err error) {
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		for _, topic := range topics {
			// consume certain partition
			if consumerOpt.Partition >= 0 {
				pc, err := consumer.ConsumePartition(topic, consumerOpt.Partition, consumerOpt.Offset)
				exitOnError(err)

				go consume(ctx, pc)
				continue
			}
			// consume all partitions
			partitions, err := consumer.Partitions(topic)
			exitOnError(err)

			for _, p := range partitions {
				pc, err := consumer.ConsumePartition(topic, p, consumerOpt.Offset)
				exitOnError(err)
				go consume(ctx, pc)
			}
		}

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		for {
			select {
			case msg := <-messages:
				fmt.Printf("%v(%v:%v) %s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			case <-signals:
				cancel()
				return
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(consumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	consumeCmd.Flags().DurationVar(&cfg.Consumer.MaxWaitTime, "maxwaittime", 250*time.Millisecond, "the maximum amount of time the broker will wait for Consumer.Fetch.Min bytes bytes to become available")
	consumeCmd.Flags().DurationVar(&cfg.Consumer.MaxProcessingTime, "maxprocessingtime", 100*time.Millisecond, "the maximum amount of time the consumer expects a message takes to process for the user.")
	consumeCmd.Flags().DurationVar(&cfg.Consumer.Retry.Backoff, "retry.backoff", 2*time.Second, "how long to wait after a failing to read from a partition before trying again")
	consumeCmd.Flags().Int32Var(&cfg.Consumer.Fetch.Min, "fetch.min", 1, "the minimum number of message bytes to fetch in a request")
	consumeCmd.Flags().Int32Var(&cfg.Consumer.Fetch.Default, "fetch.default", 32768, "the default number of message bytes to fetch from the broker in each request")
	consumeCmd.Flags().Int32Var(&cfg.Consumer.Fetch.Max, "fetch.max", 0, "the maximum number of message bytes to fetch from the broker in a single request, 0 means no limit")
	consumeCmd.Flags().DurationVar(&cfg.Consumer.Offsets.CommitInterval, "offsets.commitinterval", 1*time.Second, "how frequently to commit updated offsets")
	consumeCmd.Flags().Int64Var(&cfg.Consumer.Offsets.Initial, "offsets.initial", sarama.OffsetNewest, "the initial offset to use if no offset was previously committed")
	consumeCmd.Flags().Int64Var(&consumerOpt.Offset, "offset", sarama.OffsetNewest, "offset to consume")
	consumeCmd.Flags().Int32Var(&consumerOpt.Partition, "partition", -1, "partition to consume")
}
