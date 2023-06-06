/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"fudan.edu.cn/swz/bpf/bpf"
	"fudan.edu.cn/swz/bpf/kube"
	"github.com/spf13/cobra"
)

type LogOptions struct {
	namespace  string
	name       string
	flowOption bool
}

var logOptions LogOptions

func (o *LogOptions) Validate(cmd *cobra.Command, args []string) error {

	return nil
}

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "log the status of a running bpf",
	PreRunE: func(c *cobra.Command, args []string) error {
		return logOptions.Validate(c, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("log called")
		fmt.Println(logOptions)
		inst_name := args[0]
		if logOptions.flowOption {
			kube.FllowLog("default", inst_name)
		} else {
			kube.GetPodLog("default", inst_name)
		}
	},
}

func init() {
	logCmd.Flags().StringVar(&logOptions.name, "name", "", " Ebpf name to log")
	logCmd.Flags().StringVar(&logOptions.name, "namespace", bpf.BPF_NAMESPACE, "Namespace")
	logCmd.Flags().BoolVar(&logOptions.flowOption, "flow", false, "Watch")
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
