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

	"fudan.edu.cn/swz/bpf/bpf"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear the cache of the package",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// 循环处理所有的package
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			panic(err)
		}
		var packages []string
		// 判断是否删除全部
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			panic(err)
		}
		if all {
			packages = bpf.PackageList()
		} else {
			packages = args
		}
		for _, package_name := range packages {
			exist, err := bpf.PackageExist(package_name)
			if err != nil {
				panic(err)
			}
			if !exist {
				fmt.Printf("package %s not exist\n", package_name)
			} else {
				// 删除对应的package
				bpf.PackageDelete(package_name, force)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)

	clearCmd.Flags().BoolP("force", "f", false, "forced clear (delete the associated instance)")
	clearCmd.Flags().BoolP("all", "a", false, "clear all package")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
