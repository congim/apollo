// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shaocongcong/apollo/agent"
	_ "github.com/shaocongcong/apollo/plugins/collector/all"
	_ "github.com/shaocongcong/apollo/plugins/data/all"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: startAgent,
}

func init() {
	RootCmd.AddCommand(agentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startAgent(cmd *cobra.Command, args []string) {

	// new agent
	agentSer, err := agent.CreateAgent()
	if err != nil {
		log.Fatal(err)
	}

	// init conf
	agent.LoadConfig("agent.toml")

	// start agent
	go agentSer.Run()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	<-chSig
	// stop agent
	agentSer.Stop()
}