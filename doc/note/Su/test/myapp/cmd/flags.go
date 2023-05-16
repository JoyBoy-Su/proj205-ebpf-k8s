/*
Copyright © 2023 jiadisu@fudan.edu.cn

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
	"fmt"

	"github.com/spf13/cobra"
)

// global variable for config and flags
var Name string
var Save bool
var Password string
var Json bool
var Yaml bool

// flagsCmd represents the flags command
var flagsCmd = &cobra.Command{
	Use:   "flags",
	Short: "Test flags",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flags called")
		fmt.Println("cmd.Name = " + Name)
		if Save {
			fmt.Println("cmd.Save is true")
		}
		if Json {
			fmt.Println("cmd.format is JSON")
		} else {
			fmt.Println("cmd.format is YAML")
		}
	},
}

func init() {
	rootCmd.AddCommand(flagsCmd)
	// name flag
	flagsCmd.Flags().StringVar(&Name, "name", "default name", "my name")
	flagsCmd.MarkFlagRequired("name")
	// password flag
	flagsCmd.Flags().StringVar(&Password, "password", "", "my password")
	flagsCmd.MarkFlagsRequiredTogether("name", "password")
	// save flag
	flagsCmd.Flags().BoolVar(&Save, "save", false, "for save")
	// json and yaml flag
	flagsCmd.Flags().BoolVar(&Json, "json", false, "format json")
	flagsCmd.Flags().BoolVar(&Yaml, "yaml", false, "format yaml")
	flagsCmd.MarkFlagsMutuallyExclusive("json", "yaml")
}
