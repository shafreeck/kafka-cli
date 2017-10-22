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

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

type Topic struct {
	Name         string
	Partitions   []int32   //partition ids
	Replicas     [][]int32 //replica id of partition
	LatestOffset []int64   //latest offset of partition
}

func (t Topic) String() string {
	topic := fmt.Sprintf("%-15s\t", t.Name)
	for _, p := range t.Partitions {
		topic += fmt.Sprintf("%d%d:%d ", p, t.Replicas[p], t.LatestOffset[p])
	}
	return topic
}

func displayTopics() {
	names, err := c.Topics()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("%-15s\tpartition[replicaid...]:offset ...\n", "topic")
	fmt.Printf("--------------------------------------------------\n")
	for i := range names {
		t := Topic{Name: names[i]}
		ps, err := c.Partitions(names[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		if len(ps) < 1 {
			continue
		}

		replicas := make([][]int32, ps[len(ps)-1]+1)
		offsets := make([]int64, ps[len(ps)-1]+1)

		for _, p := range ps {
			rs, err := c.Replicas(names[i], p)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			replicas[p] = rs

			ofs, err := c.GetOffset(names[i], p, sarama.OffsetNewest)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			offsets[p] = ofs
		}
		t.Partitions = ps
		t.Replicas = replicas
		t.LatestOffset = offsets
		fmt.Println(t)
	}

}

// topicCmd represents the topic command
var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "list all topics",
	Run: func(cmd *cobra.Command, args []string) {
		displayTopics()
	},
}

func init() {
	RootCmd.AddCommand(topicsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// topicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// topicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
